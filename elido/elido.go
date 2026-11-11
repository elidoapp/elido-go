// Package elido is the Go SDK for the Elido URL-shortener API. It
// mirrors the @elido/sdk (TypeScript) shape so call sites read about
// the same in either language.
//
// Example:
//
//	c, err := elido.New(elido.Options{
//	    BaseURL: "https://api.elido.app",
//	    Token:   "elido_…",
//	})
//	if err != nil { /* … */ }
//	defer c.Close()
//
//	link, err := c.Links.Create(ctx, 42, elido.CreateLinkInput{
//	    DomainID:        1,
//	    DestinationURL:  "https://example.com/landing",
//	    Tags:            []string{"spring-2026"},
//	})
//
// The SDK is hand-written instead of generated. The Elido OpenAPI
// spec is the source of truth for the wire shape; this package keeps
// the Go ergonomics (typed inputs, context.Context, errors as values).
package elido

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Options configures the client. Only BaseURL is required.
type Options struct {
	// Base URL of the API (e.g. "https://api.elido.app",
	// "http://localhost:8085"). Trailing slash is stripped.
	BaseURL string

	// Per-service overrides. Default to BaseURL when empty.
	AnalyticsBaseURL string
	BillingBaseURL   string
	WebhookBaseURL   string
	DomainBaseURL    string
	QRBaseURL        string
	MetadataBaseURL  string

	// Auth (first non-empty wins, in this order).
	Token     string // Bearer API key — production scripting
	Cookie    string // raw Cookie header — server-side proxies
	DevUserID string // X-Dev-User-ID — local-dev only

	// HTTPClient lets the caller bring their own transport (custom
	// timeouts, RoundTripper instrumentation, etc). When nil, the SDK
	// builds a 10s-timeout client.
	HTTPClient *http.Client
}

// Client is the SDK entry point. Resource namespaces are exposed as
// fields and share the same auth state.
type Client struct {
	opts Options
	http *http.Client

	// Per-service base URLs, normalised once at construction.
	base, analytics, billing, webhook, domain, qr, metadata string

	// Resource namespaces.
	Workspaces   *Workspaces
	Links        *Links
	Folders      *Folders
	APIKeys      *APIKeys
	MachineUsers *MachineUsers
	Invitations  *Invitations
	Domains      *Domains
	Analytics    *Analytics
	Billing      *Billing
	Pixels       *Pixels
	Compliance   *Compliance
	Galleries    *Galleries
	Microsites   *Microsites
	SmsCampaigns *SmsCampaigns
	SmsOptOuts   *SmsOptOuts
}

// New builds a Client. Returns an error only when BaseURL is empty.
func New(opts Options) (*Client, error) {
	if strings.TrimSpace(opts.BaseURL) == "" {
		return nil, errors.New("elido: BaseURL is required")
	}
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	c := &Client{
		opts:      opts,
		http:      httpClient,
		base:      strings.TrimRight(opts.BaseURL, "/"),
		analytics: orFallback(opts.AnalyticsBaseURL, opts.BaseURL),
		billing:   orFallback(opts.BillingBaseURL, opts.BaseURL),
		webhook:   orFallback(opts.WebhookBaseURL, opts.BaseURL),
		domain:    orFallback(opts.DomainBaseURL, opts.BaseURL),
		qr:        orFallback(opts.QRBaseURL, opts.BaseURL),
		metadata:  orFallback(opts.MetadataBaseURL, opts.BaseURL),
	}
	c.Workspaces = &Workspaces{c: c}
	c.Links = &Links{c: c}
	c.Folders = &Folders{c: c}
	c.APIKeys = &APIKeys{c: c}
	c.MachineUsers = &MachineUsers{c: c}
	c.Invitations = &Invitations{c: c}
	c.Domains = &Domains{c: c}
	c.Analytics = &Analytics{c: c}
	c.Billing = &Billing{c: c}
	c.Pixels = &Pixels{c: c}
	c.Compliance = &Compliance{c: c}
	c.Galleries = &Galleries{c: c}
	c.Microsites = &Microsites{c: c}
	c.SmsCampaigns = &SmsCampaigns{c: c}
	c.SmsOptOuts = &SmsOptOuts{c: c}
	return c, nil
}

// Close is a noop today (the underlying http.Client has nothing to
// close in the std-lib transport) but kept on the API so callers can
// wire it into defer cleanup without churning when we add a pooled
// transport.
func (c *Client) Close() error { return nil }

func orFallback(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return strings.TrimRight(fallback, "/")
	}
	return strings.TrimRight(s, "/")
}

// APIError is returned for any non-2xx response. The wire body is
// unmarshalled into Body when it parses as JSON; otherwise it stays
// as a raw string.
type APIError struct {
	Status     int
	Message    string
	Body       any
	BodyString string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("elido api %d: %s", e.Status, e.Message)
}

// IsAPIError reports whether err is an APIError. Errors from the
// transport layer (DNS, TCP, JSON parse on success path) are returned
// as plain errors and don't satisfy this check.
func IsAPIError(err error) bool {
	var ae *APIError
	return errors.As(err, &ae)
}

func (c *Client) request(
	ctx context.Context,
	base, method, path string,
	body any,
	out any,
	query url.Values,
) error {
	full := base + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("elido: marshal body: %w", err)
		}
		reader = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, full, reader)
	if err != nil {
		return fmt.Errorf("elido: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	switch {
	case c.opts.Token != "":
		req.Header.Set("Authorization", "Bearer "+c.opts.Token)
	case c.opts.Cookie != "":
		req.Header.Set("Cookie", c.opts.Cookie)
	case c.opts.DevUserID != "":
		req.Header.Set("X-Dev-User-ID", c.opts.DevUserID)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("elido: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		ae := &APIError{Status: resp.StatusCode, BodyString: string(raw)}
		var parsed map[string]any
		if json.Unmarshal(raw, &parsed) == nil {
			ae.Body = parsed
			if msg, ok := parsed["error"].(string); ok && msg != "" {
				ae.Message = msg
			}
		}
		if ae.Message == "" {
			ae.Message = fmt.Sprintf("%s %s → %d", method, path, resp.StatusCode)
		}
		return ae
	}

	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("elido: decode response: %w", err)
	}
	return nil
}

// joinTags is a tiny helper for the links list endpoint that takes
// a comma-separated tag filter rather than repeated query keys.
func joinTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ",")
}

// itoa wraps strconv so call sites stay short.
func itoa(i int64) string { return strconv.FormatInt(i, 10) }

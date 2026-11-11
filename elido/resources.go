package elido

import (
	"context"
	"net/url"
)

// Workspaces — /v1/workspaces

type Workspaces struct{ c *Client }

type Workspace struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PlanTier  string `json:"plan_tier"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WorkspacesListResponse struct {
	Items []struct {
		Workspace
		Role string `json:"role"`
	} `json:"items"`
}

func (w *Workspaces) List(ctx context.Context) (*WorkspacesListResponse, error) {
	out := &WorkspacesListResponse{}
	return out, w.c.request(ctx, w.c.base, "GET", "/v1/workspaces", nil, out, nil)
}

func (w *Workspaces) Get(ctx context.Context, workspaceID int64) (*Workspace, error) {
	out := &Workspace{}
	return out, w.c.request(ctx, w.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID), nil, out, nil)
}

type CreateWorkspaceInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (w *Workspaces) Create(ctx context.Context, in CreateWorkspaceInput) (*Workspace, error) {
	out := &Workspace{}
	return out, w.c.request(ctx, w.c.base, "POST", "/v1/workspaces", in, out, nil)
}

type UpdateWorkspaceInput struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

func (w *Workspaces) Update(ctx context.Context, workspaceID int64, patch UpdateWorkspaceInput) (*Workspace, error) {
	out := &Workspace{}
	return out, w.c.request(ctx, w.c.base, "PATCH", "/v1/workspaces/"+itoa(workspaceID), patch, out, nil)
}

// Links — /v1/workspaces/{id}/links

type Links struct{ c *Client }

type Link struct {
	ID             int64    `json:"id"`
	WorkspaceID    int64    `json:"workspace_id"`
	DomainID       int64    `json:"domain_id"`
	Slug           string   `json:"slug"`
	DestinationURL string   `json:"destination_url"`
	Title          *string  `json:"title"`
	Tags           []string `json:"tags"`
	Status         string   `json:"status"`
	ExpiresAt      *string  `json:"expires_at"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	RedirectStatus int      `json:"redirect_status,omitempty"`
	MaxClicks      *int     `json:"max_clicks,omitempty"`
	FolderID       *int64   `json:"folder_id,omitempty"`
	RotationMode   string   `json:"rotation_mode,omitempty"`
	ShowSplash     bool     `json:"show_splash,omitempty"`
	// end-to-end encrypted link. When IsEncrypted is true,
	// DestinationURL is the sentinel "encrypted://" and the real target
	// is recovered client-side via AES-GCM with the key in the URL
	// fragment (`#k=…`).
	IsEncrypted bool `json:"is_encrypted,omitempty"`
}

type ListLinksOptions struct {
	Limit  int
	Offset int
	Tags   []string
	Q      string
	Folder string // numeric ID as string, or "none" for unfiled
}

type PaginatedLinks struct {
	Items  []Link `json:"items"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (l *Links) List(ctx context.Context, workspaceID int64, opts ListLinksOptions) (*PaginatedLinks, error) {
	q := url.Values{}
	if opts.Limit > 0 {
		q.Set("limit", itoa(int64(opts.Limit)))
	}
	if opts.Offset > 0 {
		q.Set("offset", itoa(int64(opts.Offset)))
	}
	if t := joinTags(opts.Tags); t != "" {
		q.Set("tags", t)
	}
	if opts.Q != "" {
		q.Set("q", opts.Q)
	}
	if opts.Folder != "" {
		q.Set("folder", opts.Folder)
	}
	out := &PaginatedLinks{}
	return out, l.c.request(ctx, l.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/links", nil, out, q)
}

type CreateLinkInput struct {
	DomainID       int64    `json:"domain_id"`
	Slug           string   `json:"slug,omitempty"`
	DestinationURL string   `json:"destination_url,omitempty"`
	Title          *string  `json:"title,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	RedirectStatus *int     `json:"redirect_status,omitempty"`
	ShowSplash     *bool    `json:"show_splash,omitempty"`
	// set IsEncrypted with EncryptedTargetB64/EncIvB64
	// for an end-to-end encrypted link. DestinationURL should be empty
	// in that case (the server forces it to the sentinel). The
	// encryption itself happens in the browser via Web Crypto API —
	// this SDK only carries the wire shape.
	IsEncrypted          bool   `json:"is_encrypted,omitempty"`
	EncryptedTargetB64   string `json:"encrypted_target,omitempty"`
	EncIvB64             string `json:"enc_iv,omitempty"`
	EncAlg               string `json:"enc_alg,omitempty"`
	EncPassphraseSaltB64 string `json:"enc_passphrase_salt,omitempty"`
	EncKdf               string `json:"enc_kdf,omitempty"`
	EncKdfIters          int32  `json:"enc_kdf_iters,omitempty"`
}

// RotateLinkEncryptionInput carries a fresh ciphertext + IV for an
// existing E2E encrypted link.
type RotateLinkEncryptionInput struct {
	EncryptedTargetB64   string `json:"encrypted_target"`
	EncIvB64             string `json:"enc_iv"`
	EncAlg               string `json:"enc_alg"`
	EncPassphraseSaltB64 string `json:"enc_passphrase_salt,omitempty"`
	EncKdf               string `json:"enc_kdf,omitempty"`
	EncKdfIters          int32  `json:"enc_kdf_iters,omitempty"`
}

// RotateEncryption posts a fresh crypto envelope for an existing E2E
// encrypted link. The link MUST already be encrypted (the endpoint
// rejects the request otherwise). The new key never leaves the
// caller's process — the workspace member pastes the original URL into
// the dashboard, the browser re-encrypts it with a freshly-generated
// key, and only the ciphertext reaches this method.
func (l *Links) RotateEncryption(ctx context.Context, workspaceID, linkID int64, in RotateLinkEncryptionInput) error {
	return l.c.request(ctx, l.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/links/"+itoa(linkID)+"/encryption/rotate",
		in, nil, nil)
}

func (l *Links) Create(ctx context.Context, workspaceID int64, in CreateLinkInput) (*Link, error) {
	out := &Link{}
	return out, l.c.request(ctx, l.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/links", in, out, nil)
}

// UpdateLinkInput is a free-form patch — only the fields you want to
// change need a non-nil value. Maps onto the PATCH endpoint
// which uses COALESCE under the hood.
type UpdateLinkInput struct {
	Slug           *string  `json:"slug,omitempty"`
	DestinationURL *string  `json:"destination_url,omitempty"`
	Title          *string  `json:"title,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Status         *string  `json:"status,omitempty"`
	ExpiresAt      *string  `json:"expires_at,omitempty"`
	Password       *string  `json:"password,omitempty"`
	RedirectStatus *int     `json:"redirect_status,omitempty"`
	MaxClicks      *int     `json:"max_clicks,omitempty"`
	FolderID       *int64   `json:"folder_id,omitempty"`
	RotationMode   *string  `json:"rotation_mode,omitempty"`
	ShowSplash     *bool    `json:"show_splash,omitempty"`
}

func (l *Links) Update(ctx context.Context, workspaceID, linkID int64, patch UpdateLinkInput) (*Link, error) {
	out := &Link{}
	return out, l.c.request(ctx, l.c.base, "PATCH", "/v1/workspaces/"+itoa(workspaceID)+"/links/"+itoa(linkID), patch, out, nil)
}

func (l *Links) Delete(ctx context.Context, workspaceID, linkID int64) error {
	return l.c.request(ctx, l.c.base, "DELETE", "/v1/workspaces/"+itoa(workspaceID)+"/links/"+itoa(linkID), nil, nil, nil)
}

// Folders — /v1/workspaces/{id}/folders

type Folders struct{ c *Client }

type Folder struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	LinkCount   int64  `json:"link_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type FoldersListResponse struct {
	Items []Folder `json:"items"`
}

func (f *Folders) List(ctx context.Context, workspaceID int64) (*FoldersListResponse, error) {
	out := &FoldersListResponse{}
	return out, f.c.request(ctx, f.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/folders", nil, out, nil)
}

type CreateFolderInput struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

func (f *Folders) Create(ctx context.Context, workspaceID int64, in CreateFolderInput) (*Folder, error) {
	out := &Folder{}
	return out, f.c.request(ctx, f.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/folders", in, out, nil)
}

// APIKeys — /v1/workspaces/{id}/api-keys

type APIKeys struct{ c *Client }

type APIKey struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Prefix     string  `json:"prefix"`
	Role       string  `json:"role"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
	ExpiresAt  *string `json:"expires_at"`
	RevokedAt  *string `json:"revoked_at"`
}

type CreateAPIKeyInput struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type CreateAPIKeyResponse struct {
	APIKey
	Token string `json:"token"` // shown once
}

func (k *APIKeys) Create(ctx context.Context, workspaceID int64, in CreateAPIKeyInput) (*CreateAPIKeyResponse, error) {
	out := &CreateAPIKeyResponse{}
	return out, k.c.request(ctx, k.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/api-keys", in, out, nil)
}

// MachineUsers — /v1/workspaces/{id}/machine-users.
//
// A machine user is a non-human identity owned by a workspace. Tokens minted
// under it authenticate the same way as a PAT (Authorization: Bearer
// elido_...) but the request context carries a MachineUser instead of a
// User; RBAC reads the role off the machine_users row.

type MachineUsers struct{ c *Client }

type MachineUser struct {
	ID               int64   `json:"id"`
	WorkspaceID      int64   `json:"workspace_id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Role             string  `json:"role"`
	CreatedBy        *int64  `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
	LastUsedAt       *string `json:"last_used_at"`
	DisabledAt       *string `json:"disabled_at"`
	ActiveTokenCount int64   `json:"active_token_count"`
}

type MachineUsersListResponse struct {
	Items []MachineUser `json:"items"`
}

type CreateMachineUserInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Role        string `json:"role,omitempty"`
}

type UpdateMachineUserInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Role        *string `json:"role,omitempty"`
	Disabled    *bool   `json:"disabled,omitempty"`
}

type MachineUserToken struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Prefix         string  `json:"prefix"`
	Role           string  `json:"role"`
	LastUsedAt     *string `json:"last_used_at,omitempty"`
	LastUsedIP     *string `json:"last_used_ip,omitempty"`
	CreatedAt      string  `json:"created_at"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
	RevokedAt      *string `json:"revoked_at,omitempty"`
	TotalRequests  int64   `json:"total_requests,omitempty"`
	FailedRequests int64   `json:"failed_requests,omitempty"`
}

type MachineUserTokensListResponse struct {
	Items []MachineUserToken `json:"items"`
}

type MintMachineUserTokenInput struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// MintMachineUserTokenResponse carries the freshly minted plaintext token.
// Plaintext is returned exactly once — callers MUST persist it.
type MintMachineUserTokenResponse struct {
	Token     MachineUserToken `json:"token"`
	Plaintext string           `json:"plaintext"`
}

func (m *MachineUsers) List(ctx context.Context, workspaceID int64) (*MachineUsersListResponse, error) {
	out := &MachineUsersListResponse{}
	return out, m.c.request(ctx, m.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users", nil, out, nil)
}

func (m *MachineUsers) Create(ctx context.Context, workspaceID int64, in CreateMachineUserInput) (*MachineUser, error) {
	out := &MachineUser{}
	return out, m.c.request(ctx, m.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users", in, out, nil)
}

func (m *MachineUsers) Get(ctx context.Context, workspaceID, machineUserID int64) (*MachineUser, error) {
	out := &MachineUser{}
	return out, m.c.request(ctx, m.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID), nil, out, nil)
}

func (m *MachineUsers) Update(ctx context.Context, workspaceID, machineUserID int64, patch UpdateMachineUserInput) (*MachineUser, error) {
	out := &MachineUser{}
	return out, m.c.request(ctx, m.c.base, "PATCH",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID), patch, out, nil)
}

func (m *MachineUsers) Delete(ctx context.Context, workspaceID, machineUserID int64) error {
	return m.c.request(ctx, m.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID), nil, nil, nil)
}

func (m *MachineUsers) ListTokens(ctx context.Context, workspaceID, machineUserID int64) (*MachineUserTokensListResponse, error) {
	out := &MachineUserTokensListResponse{}
	return out, m.c.request(ctx, m.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID)+"/tokens", nil, out, nil)
}

func (m *MachineUsers) MintToken(ctx context.Context, workspaceID, machineUserID int64, in MintMachineUserTokenInput) (*MintMachineUserTokenResponse, error) {
	out := &MintMachineUserTokenResponse{}
	return out, m.c.request(ctx, m.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID)+"/tokens", in, out, nil)
}

func (m *MachineUsers) RevokeToken(ctx context.Context, workspaceID, machineUserID, tokenID int64) error {
	return m.c.request(ctx, m.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/machine-users/"+itoa(machineUserID)+"/tokens/"+itoa(tokenID), nil, nil, nil)
}

// Invitations — /v1/workspaces/{id}/invitations

type Invitations struct{ c *Client }

type Invitation struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Token       string `json:"token,omitempty"` // returned once on create
	CreatedAt   string `json:"created_at"`
}

type CreateInvitationInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (i *Invitations) Create(ctx context.Context, workspaceID int64, in CreateInvitationInput) (*Invitation, error) {
	out := &Invitation{}
	return out, i.c.request(ctx, i.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/invitations", in, out, nil)
}

// Domains — /v1/workspaces/{id}/domains
//
// Custom-domain verification is handled by a background reconcile loop.
// Verification is asynchronous: Create returns a row with
// status="pending" and the DNS instructions the user has to publish;
// callers then poll Get() until status flips to "verified" (and
// eventually "active" once Caddy issues a cert).

type Domains struct{ c *Client }

type Domain struct {
	ID          int64   `json:"id"`
	WorkspaceID int64   `json:"workspace_id"`
	Hostname    string  `json:"hostname"`
	VerifyToken *string `json:"verify_token"`
	Status      string  `json:"status"`
	IsWildcard  bool    `json:"is_wildcard"`
	CreatedAt   string  `json:"created_at"`
	VerifiedAt  *string `json:"verified_at"`
	FailedAt    *string `json:"failed_at"`
	FailReason  *string `json:"fail_reason"`
}

type DomainsListResponse struct {
	Items []Domain `json:"items"`
}

func (d *Domains) List(ctx context.Context, workspaceID int64) (*DomainsListResponse, error) {
	out := &DomainsListResponse{}
	return out, d.c.request(ctx, d.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/domains", nil, out, nil)
}

// CreateDomainInput is the request body for POST
// /v1/workspaces/{id}/domains. Hostname must be a lowercase FQDN —
// wildcards use the literal "*." prefix (e.g. "*.example.com").
type CreateDomainInput struct {
	Hostname   string `json:"hostname"`
	IsWildcard bool   `json:"is_wildcard,omitempty"`
}

// CreateDomainResponse is the body of a successful Create. Instructions
// is a small map shaped `{"txt_record": "...", "cname_record": "..."}`
// the user has to publish before the server will verify the row.
type CreateDomainResponse struct {
	Domain       Domain            `json:"domain"`
	Instructions map[string]string `json:"instructions"`
}

// Create registers a new pending custom domain. Verification runs in
// the background — poll Get() until status flips to "verified".
func (d *Domains) Create(ctx context.Context, workspaceID int64, in CreateDomainInput) (*CreateDomainResponse, error) {
	out := &CreateDomainResponse{}
	path := "/v1/workspaces/" + itoa(workspaceID) + "/domains"
	return out, d.c.request(ctx, d.c.base, "POST", path, in, out, nil)
}

// Get fetches a single domain by id. Returns *APIError with Status 404
// when the id doesn't exist in this workspace.
func (d *Domains) Get(ctx context.Context, workspaceID, domainID int64) (*Domain, error) {
	out := &Domain{}
	path := "/v1/workspaces/" + itoa(workspaceID) + "/domains/" + itoa(domainID)
	return out, d.c.request(ctx, d.c.base, "GET", path, nil, out, nil)
}

// Delete hard-deletes a custom domain. Idempotent on the wire — 404 is
// returned as an APIError, callers that want "delete or no-op" should
// check IsAPIError + Status == 404.
func (d *Domains) Delete(ctx context.Context, workspaceID, domainID int64) error {
	path := "/v1/workspaces/" + itoa(workspaceID) + "/domains/" + itoa(domainID)
	return d.c.request(ctx, d.c.base, "DELETE", path, nil, nil, nil)
}

// VerifyDNS pushes the row to the front of the verification
// queue, so the next pass picks it up immediately rather than waiting
// up to VerifyInterval. The returned Domain is the current row —
// status is usually still "pending" at this point; callers poll Get()
// to observe the transition to "verified".
func (d *Domains) VerifyDNS(ctx context.Context, workspaceID, domainID int64) (*Domain, error) {
	out := &Domain{}
	path := "/v1/workspaces/" + itoa(workspaceID) + "/domains/" + itoa(domainID) + "/verify"
	return out, d.c.request(ctx, d.c.base, "POST", path, nil, out, nil)
}

// AppLinksInput is the PATCH body for the app-links endpoint. Either
// field may be nil — the column is left unchanged when nil.
type AppLinksInput struct {
	AppleAppSiteAssociation any `json:"apple_app_site_association,omitempty"`
	AndroidAssetLinks       any `json:"android_asset_links,omitempty"`
}

// UpdateAppLinks pushes Universal Links / App Links manifests for a
// custom domain. The body is stored verbatim and served at the
// well-known path.
func (d *Domains) UpdateAppLinks(ctx context.Context, workspaceID, domainID int64, in AppLinksInput) (*Domain, error) {
	out := &Domain{}
	path := "/v1/workspaces/" + itoa(workspaceID) + "/domains/" + itoa(domainID) + "/app-links"
	return out, d.c.request(ctx, d.c.base, "PATCH", path, in, out, nil)
}

// Analytics — /v1/analytics/...

type Analytics struct{ c *Client }

type TimeseriesPoint struct {
	TS    string `json:"ts"`
	Count int64  `json:"count"`
}

type TimeseriesResponse struct {
	Points []TimeseriesPoint `json:"points"`
}

type TimeseriesOptions struct {
	From     string // YYYY-MM-DD
	To       string // YYYY-MM-DD
	Interval string // "hour" | "day"
	LinkID   int64
}

func (a *Analytics) Timeseries(ctx context.Context, workspaceID int64, opts TimeseriesOptions) (*TimeseriesResponse, error) {
	q := url.Values{"from": {opts.From}, "to": {opts.To}}
	if opts.Interval != "" {
		q.Set("interval", opts.Interval)
	}
	if opts.LinkID > 0 {
		q.Set("link_id", itoa(opts.LinkID))
	}
	out := &TimeseriesResponse{}
	return out, a.c.request(ctx, a.c.analytics, "GET",
		"/v1/analytics/workspaces/"+itoa(workspaceID)+"/timeseries", nil, out, q)
}

type TopLink struct {
	LinkID int64 `json:"link_id"`
	Count  int64 `json:"count"`
}

type TopLinksResponse struct {
	Items []TopLink `json:"items"`
}

type TopLinksOptions struct {
	From  string
	To    string
	Limit int
}

func (a *Analytics) TopLinks(ctx context.Context, workspaceID int64, opts TopLinksOptions) (*TopLinksResponse, error) {
	q := url.Values{"from": {opts.From}, "to": {opts.To}}
	if opts.Limit > 0 {
		q.Set("limit", itoa(int64(opts.Limit)))
	}
	out := &TopLinksResponse{}
	return out, a.c.request(ctx, a.c.analytics, "GET",
		"/v1/analytics/workspaces/"+itoa(workspaceID)+"/links/top", nil, out, q)
}

// LTVCohort is one row of the lifetime-value report — engagement
// (visitors / clicks) from ClickHouse joined with attributed revenue
// (Stripe / Shopify conversions) from Postgres, scoped to a single
// currency. Cohorts with no conversions in the requested currency
// report RevenueCents = 0.
type LTVCohort struct {
	CohortWeek             string  `json:"cohort_week"`
	LinksAcquired          uint64  `json:"links_acquired"`
	Visitors               uint64  `json:"visitors"`
	ClicksTotal            uint64  `json:"clicks_total"`
	RevenueCents           int64   `json:"revenue_cents"`
	Currency               string  `json:"currency"`
	ClicksPerLink          float64 `json:"clicks_per_link"`
	RevenuePerLinkCents    float64 `json:"revenue_per_link_cents"`
	RevenuePerVisitorCents float64 `json:"revenue_per_visitor_cents"`
}

type LTVResponse struct {
	Currency string      `json:"currency"`
	Cohorts  []LTVCohort `json:"cohorts"`
}

type LTVOptions struct {
	From     string // YYYY-MM-DD
	To       string // YYYY-MM-DD
	LinkID   int64
	Currency string // 3-letter ISO; default "USD" on the server
}

// LTV returns the lifetime-value report — one row per acquisition cohort
// (ISO week of a link's first click).
func (a *Analytics) LTV(ctx context.Context, workspaceID int64, opts LTVOptions) (*LTVResponse, error) {
	q := url.Values{"from": {opts.From}, "to": {opts.To}}
	if opts.LinkID > 0 {
		q.Set("link_id", itoa(opts.LinkID))
	}
	if opts.Currency != "" {
		q.Set("currency", opts.Currency)
	}
	out := &LTVResponse{}
	return out, a.c.request(ctx, a.c.analytics, "GET",
		"/v1/analytics/workspaces/"+itoa(workspaceID)+"/ltv", nil, out, q)
}

// Billing — /v1/billing/...  (billing service)

type Billing struct{ c *Client }

type Plan struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	Tier            string         `json:"tier"`
	PriceCents      int64          `json:"price_cents"`
	Currency        string         `json:"currency"`
	BillingInterval string         `json:"billing_interval"`
	Features        map[string]any `json:"features"`
	IsActive        bool           `json:"is_active"`
}

type PlansResponse struct {
	Items []Plan `json:"items"`
}

func (b *Billing) Plans(ctx context.Context) (*PlansResponse, error) {
	out := &PlansResponse{}
	return out, b.c.request(ctx, b.c.billing, "GET", "/v1/billing/plans", nil, out, nil)
}

type Subscription struct {
	ID                 int64   `json:"id"`
	WorkspaceID        int64   `json:"workspace_id"`
	PlanID             int64   `json:"plan_id"`
	Status             string  `json:"status"`
	CurrentPeriodStart *string `json:"current_period_start"`
	CurrentPeriodEnd   *string `json:"current_period_end"`
}

func (b *Billing) Subscription(ctx context.Context, workspaceID int64) (*Subscription, error) {
	out := &Subscription{}
	return out, b.c.request(ctx, b.c.billing, "GET",
		"/v1/billing/workspaces/"+itoa(workspaceID)+"/subscription", nil, out, nil)
}

type CheckoutResponse struct {
	SubscriptionID  int64  `json:"subscription_id"`
	Status          string `json:"status"`
	RedirectURL     string `json:"redirect_url,omitempty"`
	LiqPayData      string `json:"liqpay_data,omitempty"`
	LiqPaySignature string `json:"liqpay_signature,omitempty"`
}

func (b *Billing) Checkout(ctx context.Context, workspaceID, planID int64) (*CheckoutResponse, error) {
	body := map[string]int64{"plan_id": planID}
	out := &CheckoutResponse{}
	return out, b.c.request(ctx, b.c.billing, "POST",
		"/v1/billing/workspaces/"+itoa(workspaceID)+"/checkout", body, out, nil)
}

// Compliance — /v1/workspaces/{id}/compliance and /compliance/evidence

type Compliance struct{ c *Client }

type WorkspaceCompliance struct {
	WorkspaceID               int64  `json:"workspace_id"`
	Mode                      string `json:"mode"`
	AuditRetentionDays        *int32 `json:"audit_retention_days,omitempty"`
	IPAllowlistEnabled        bool   `json:"ip_allowlist_enabled"`
	DefaultAuditRetentionDays int32  `json:"default_audit_retention_days"`
}

type UpdateComplianceInput struct {
	Mode                *string `json:"mode,omitempty"`
	AuditRetentionDays  *int32  `json:"audit_retention_days,omitempty"`
	ResetAuditRetention bool    `json:"reset_audit_retention,omitempty"`
}

// EvidenceSnapshot mirrors the JSON returned by the SOC 2 / HIPAA
// evidence-collection endpoint. The Go shape stays denormalised
// (one struct per logical area) so callers can pull `snap.Auditing`
// without crawling a generic map.
type EvidenceSnapshot struct {
	WorkspaceID    int64                `json:"workspace_id"`
	CollectedAt    string               `json:"collected_at"`
	Posture        EvidencePosture      `json:"posture"`
	Access         EvidenceAccess       `json:"access"`
	Auditing       EvidenceAuditing     `json:"auditing"`
	Integrations   EvidenceIntegrations `json:"integrations"`
	ControlsMatrix string               `json:"controls_matrix_url"`
}

type EvidencePosture struct {
	Mode                      string `json:"mode"`
	AuditRetentionDays        *int32 `json:"audit_retention_days,omitempty"`
	DefaultAuditRetentionDays int32  `json:"default_audit_retention_days"`
	IPAllowlistEnabled        bool   `json:"ip_allowlist_enabled"`
	IPRulesCount              int64  `json:"ip_rules_count"`
}

type EvidenceAccess struct {
	MembersTotal       int64            `json:"members_total"`
	MembersByRole      map[string]int64 `json:"members_by_role"`
	CustomRolesCount   int64            `json:"custom_roles_count"`
	PendingInvitations int64            `json:"pending_invitations"`
	APIKeysActive      int64            `json:"api_keys_active"`
	APIKeysLastCreated *string          `json:"api_keys_last_created,omitempty"`
}

type EvidenceAuditing struct {
	EventsLast30Days int64   `json:"events_last_30_days"`
	EventsLast90Days int64   `json:"events_last_90_days"`
	OldestRetained   *string `json:"oldest_retained,omitempty"`
	NewestEvent      *string `json:"newest_event,omitempty"`
}

type EvidenceIntegrations struct {
	SSOBound       bool  `json:"sso_bound"`
	SCIMBound      bool  `json:"scim_bound"`
	SCIMEnabled    bool  `json:"scim_enabled"`
	BAACountSigned int64 `json:"baa_signed_count"`
}

func (cmp *Compliance) Get(ctx context.Context, workspaceID int64) (*WorkspaceCompliance, error) {
	out := &WorkspaceCompliance{}
	return out, cmp.c.request(ctx, cmp.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/compliance", nil, out, nil)
}

func (cmp *Compliance) Update(ctx context.Context, workspaceID int64, in UpdateComplianceInput) (*WorkspaceCompliance, error) {
	out := &WorkspaceCompliance{}
	return out, cmp.c.request(ctx, cmp.c.base, "PUT",
		"/v1/workspaces/"+itoa(workspaceID)+"/compliance", in, out, nil)
}

// Evidence returns the SOC 2 / HIPAA evidence snapshot. Pair with
// EvidenceBundleURL when an auditor wants the downloadable ZIP.
func (cmp *Compliance) Evidence(ctx context.Context, workspaceID int64) (*EvidenceSnapshot, error) {
	out := &EvidenceSnapshot{}
	return out, cmp.c.request(ctx, cmp.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/compliance/evidence", nil, out, nil)
}

// EvidenceBundleURL returns the URL of the ZIP-streamed evidence
// bundle. Callers typically GET it themselves with their own
// Authorization header so they can stream the body to disk without
// buffering a 100k-event audit log in memory.
func (cmp *Compliance) EvidenceBundleURL(workspaceID int64) string {
	return cmp.c.base + "/v1/workspaces/" + itoa(workspaceID) + "/compliance/evidence.zip"
}

// Pixels — /v1/workspaces/{id}/pixels

type Pixels struct{ c *Client }

type WorkspacePixels struct {
	WorkspaceID int64   `json:"workspace_id"`
	MetaPixelID *string `json:"meta_pixel_id"`
	GoogleID    *string `json:"google_id"`
	LinkedInID  *string `json:"linkedin_id"`
	TikTokID    *string `json:"tiktok_id"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type UpsertPixelsInput struct {
	MetaPixelID *string `json:"meta_pixel_id,omitempty"`
	GoogleID    *string `json:"google_id,omitempty"`
	LinkedInID  *string `json:"linkedin_id,omitempty"`
	TikTokID    *string `json:"tiktok_id,omitempty"`
}

func (p *Pixels) Get(ctx context.Context, workspaceID int64) (*WorkspacePixels, error) {
	out := &WorkspacePixels{}
	return out, p.c.request(ctx, p.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/pixels", nil, out, nil)
}

func (p *Pixels) Upsert(ctx context.Context, workspaceID int64, in UpsertPixelsInput) (*WorkspacePixels, error) {
	out := &WorkspacePixels{}
	return out, p.c.request(ctx, p.c.base, "PUT",
		"/v1/workspaces/"+itoa(workspaceID)+"/pixels", in, out, nil)
}

func (p *Pixels) Delete(ctx context.Context, workspaceID int64) error {
	return p.c.request(ctx, p.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/pixels", nil, nil, nil)
}

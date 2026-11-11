package elido

import (
	"context"
	"net/url"
	"strconv"
)

// Microsites — /v1/workspaces/{id}/microsites.
//
// A microsite is a structured PRODUCT landing page — hero, specs,
// manuals, warranty, FAQ, related products. Public renderer lives at
// /p/<slug> on the bio renderer; QR codes (GS1 Digital Link or otherwise)
// can point at the public URL like any other short link.

type Microsites struct{ c *Client }

type Microsite struct {
	ID             int64   `json:"id"`
	WorkspaceID    int64   `json:"workspace_id"`
	Slug           string  `json:"slug"`
	ProductName    string  `json:"product_name"`
	SKU            *string `json:"sku,omitempty"`
	GTIN           *string `json:"gtin,omitempty"`
	Manufacturer   *string `json:"manufacturer,omitempty"`
	ProductVersion *string `json:"product_version,omitempty"`
	Tagline        string  `json:"tagline"`
	Description    string  `json:"description"`
	HeroImageURL   *string `json:"hero_image_url,omitempty"`
	HeroVideoURL   *string `json:"hero_video_url,omitempty"`
	PrimaryColor   *string `json:"primary_color,omitempty"`
	SecondaryColor *string `json:"secondary_color,omitempty"`
	Theme          string  `json:"theme"`
	Visibility     string  `json:"visibility"`
	CustomDomainID *int64  `json:"custom_domain_id,omitempty"`
	SEOTitle       *string `json:"seo_title,omitempty"`
	SEODescription *string `json:"seo_description,omitempty"`
	SEOOGImageURL  *string `json:"seo_og_image_url,omitempty"`
	Locale         string  `json:"locale"`
	CreatedBy      *int64  `json:"created_by,omitempty"`
	PublishedAt    *string `json:"published_at,omitempty"`
	ArchivedAt     *string `json:"archived_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type MicrositeSection struct {
	ID          int64          `json:"id"`
	MicrositeID int64          `json:"microsite_id"`
	Kind        string         `json:"kind"`
	Title       *string        `json:"title,omitempty"`
	Data        map[string]any `json:"data"`
	SortOrder   int32          `json:"sort_order"`
	IsVisible   bool           `json:"is_visible"`
}

// MicrositeDetail = microsite + embedded sections (GET-by-id shape).
type MicrositeDetail struct {
	Microsite
	Sections []MicrositeSection `json:"sections"`
}

type MicrositeFormSubmission struct {
	ID             int64          `json:"id"`
	MicrositeID    int64          `json:"microsite_id"`
	SectionID      *int64         `json:"section_id,omitempty"`
	FormKind       string         `json:"form_kind"`
	Payload        map[string]any `json:"payload"`
	SubmitterEmail *string        `json:"submitter_email,omitempty"`
	UA             *string        `json:"ua,omitempty"`
	CreatedAt      string         `json:"created_at"`
}

type ListMicrositesOptions struct {
	Limit      int
	Offset     int
	Q          string
	Visibility string
	Locale     string
	Archived   bool
}

type MicrositesListResponse struct {
	Items  []Microsite `json:"items"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

func (m *Microsites) List(ctx context.Context, workspaceID int64, opts ListMicrositesOptions) (*MicrositesListResponse, error) {
	q := url.Values{}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Q != "" {
		q.Set("q", opts.Q)
	}
	if opts.Visibility != "" {
		q.Set("visibility", opts.Visibility)
	}
	if opts.Locale != "" {
		q.Set("locale", opts.Locale)
	}
	if opts.Archived {
		q.Set("archived", "true")
	}
	out := &MicrositesListResponse{}
	return out, m.c.request(ctx, m.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/microsites", nil, out, q)
}

type CreateMicrositeInput struct {
	Slug           string  `json:"slug"`
	ProductName    string  `json:"product_name"`
	SKU            *string `json:"sku,omitempty"`
	GTIN           *string `json:"gtin,omitempty"`
	Manufacturer   *string `json:"manufacturer,omitempty"`
	ProductVersion *string `json:"product_version,omitempty"`
	Tagline        string  `json:"tagline,omitempty"`
	Description    string  `json:"description,omitempty"`
	HeroImageURL   *string `json:"hero_image_url,omitempty"`
	HeroVideoURL   *string `json:"hero_video_url,omitempty"`
	PrimaryColor   *string `json:"primary_color,omitempty"`
	SecondaryColor *string `json:"secondary_color,omitempty"`
	Theme          string  `json:"theme,omitempty"`
	Visibility     string  `json:"visibility,omitempty"`
	Locale         string  `json:"locale,omitempty"`
}

func (m *Microsites) Create(ctx context.Context, workspaceID int64, in CreateMicrositeInput) (*Microsite, error) {
	out := &Microsite{}
	return out, m.c.request(ctx, m.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/microsites", in, out, nil)
}

func (m *Microsites) Get(ctx context.Context, workspaceID, micrositeID int64) (*MicrositeDetail, error) {
	out := &MicrositeDetail{}
	return out, m.c.request(ctx, m.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID), nil, out, nil)
}

type UpdateMicrositeInput struct {
	ProductName    *string `json:"product_name,omitempty"`
	SKU            *string `json:"sku,omitempty"`
	GTIN           *string `json:"gtin,omitempty"`
	Manufacturer   *string `json:"manufacturer,omitempty"`
	ProductVersion *string `json:"product_version,omitempty"`
	Tagline        *string `json:"tagline,omitempty"`
	Description    *string `json:"description,omitempty"`
	HeroImageURL   *string `json:"hero_image_url,omitempty"`
	HeroVideoURL   *string `json:"hero_video_url,omitempty"`
	PrimaryColor   *string `json:"primary_color,omitempty"`
	SecondaryColor *string `json:"secondary_color,omitempty"`
	Theme          *string `json:"theme,omitempty"`
	Visibility     *string `json:"visibility,omitempty"`
	CustomDomainID *int64  `json:"custom_domain_id,omitempty"`
	SEOTitle       *string `json:"seo_title,omitempty"`
	SEODescription *string `json:"seo_description,omitempty"`
	SEOOGImageURL  *string `json:"seo_og_image_url,omitempty"`
	Locale         *string `json:"locale,omitempty"`
	Published      *bool   `json:"published,omitempty"`
	Archived       *bool   `json:"archived,omitempty"`
}

func (m *Microsites) Update(ctx context.Context, workspaceID, micrositeID int64, in UpdateMicrositeInput) (*Microsite, error) {
	out := &Microsite{}
	return out, m.c.request(ctx, m.c.base, "PATCH", "/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID), in, out, nil)
}

func (m *Microsites) Delete(ctx context.Context, workspaceID, micrositeID int64) error {
	return m.c.request(ctx, m.c.base, "DELETE", "/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID), nil, nil, nil)
}

type CreateMicrositeSectionInput struct {
	Kind      string         `json:"kind"`
	Title     *string        `json:"title,omitempty"`
	Data      map[string]any `json:"data"`
	SortOrder int32          `json:"sort_order,omitempty"`
	IsVisible *bool          `json:"is_visible,omitempty"`
}

func (m *Microsites) CreateSection(ctx context.Context, workspaceID, micrositeID int64, in CreateMicrositeSectionInput) (*MicrositeSection, error) {
	out := &MicrositeSection{}
	return out, m.c.request(ctx, m.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID)+"/sections", in, out, nil)
}

type UpdateMicrositeSectionInput struct {
	Title     *string        `json:"title,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	SortOrder *int32         `json:"sort_order,omitempty"`
	IsVisible *bool          `json:"is_visible,omitempty"`
}

func (m *Microsites) UpdateSection(ctx context.Context, workspaceID, micrositeID, sectionID int64, in UpdateMicrositeSectionInput) (*MicrositeSection, error) {
	out := &MicrositeSection{}
	return out, m.c.request(ctx, m.c.base, "PATCH",
		"/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID)+"/sections/"+itoa(sectionID), in, out, nil)
}

func (m *Microsites) DeleteSection(ctx context.Context, workspaceID, micrositeID, sectionID int64) error {
	return m.c.request(ctx, m.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID)+"/sections/"+itoa(sectionID), nil, nil, nil)
}

type ReorderMicrositeSectionEntry struct {
	SectionID int64 `json:"section_id"`
	SortOrder int32 `json:"sort_order"`
}

type reorderSectionsBody struct {
	Entries []ReorderMicrositeSectionEntry `json:"entries"`
}

type ReorderMicrositeSectionsResponse struct {
	OK    bool `json:"ok"`
	Count int  `json:"count"`
}

func (m *Microsites) ReorderSections(ctx context.Context, workspaceID, micrositeID int64, entries []ReorderMicrositeSectionEntry) (*ReorderMicrositeSectionsResponse, error) {
	body := reorderSectionsBody{Entries: entries}
	out := &ReorderMicrositeSectionsResponse{}
	return out, m.c.request(ctx, m.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID)+"/sections/reorder", body, out, nil)
}

type ListMicrositeSubmissionsOptions struct {
	Limit  int
	Offset int
}

type MicrositeSubmissionsListResponse struct {
	Items  []MicrositeFormSubmission `json:"items"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
	Total  int                       `json:"total"`
}

func (m *Microsites) ListSubmissions(ctx context.Context, workspaceID, micrositeID int64, opts ListMicrositeSubmissionsOptions) (*MicrositeSubmissionsListResponse, error) {
	q := url.Values{}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}
	out := &MicrositeSubmissionsListResponse{}
	return out, m.c.request(ctx, m.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/microsites/"+itoa(micrositeID)+"/submissions", nil, out, q)
}

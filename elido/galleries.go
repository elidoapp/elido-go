package elido

import (
	"context"
	"net/url"
	"strconv"
)

// Galleries — /v1/workspaces/{id}/galleries.
//
// A gallery is a workspace-owned, public-facing showcase that aggregates
// short links under a single curated URL: press kits, campaign hubs,
// "all my links" pages. Many per workspace, distinct from bio pages
// (one per human).

type Galleries struct{ c *Client }

type Gallery struct {
	ID             int64   `json:"id"`
	WorkspaceID    int64   `json:"workspace_id"`
	Slug           string  `json:"slug"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Theme          string  `json:"theme"`
	Visibility     string  `json:"visibility"`
	AccentColor    *string `json:"accent_color,omitempty"`
	CoverImageURL  *string `json:"cover_image_url,omitempty"`
	CustomDomainID *int64  `json:"custom_domain_id,omitempty"`
	SEOOGImageURL  *string `json:"seo_og_image_url,omitempty"`
	SEOTitle       *string `json:"seo_title,omitempty"`
	SEODescription *string `json:"seo_description,omitempty"`
	CreatedBy      *int64  `json:"created_by,omitempty"`
	PublishedAt    *string `json:"published_at,omitempty"`
	ArchivedAt     *string `json:"archived_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type GallerySection struct {
	ID          int64  `json:"id"`
	GalleryID   int64  `json:"gallery_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sort_order"`
}

type GalleryItem struct {
	ID          int64   `json:"id"`
	GalleryID   int64   `json:"gallery_id"`
	LinkID      int64   `json:"link_id"`
	SectionID   *int64  `json:"section_id,omitempty"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	SortOrder   int32   `json:"sort_order"`
}

// GalleryDetail is the GET-by-id shape: gallery + embedded sections + items.
type GalleryDetail struct {
	Gallery
	Sections []GallerySection `json:"sections"`
	Items    []GalleryItem    `json:"items"`
}

type ListGalleriesOptions struct {
	Limit      int
	Offset     int
	Q          string
	Visibility string // "" | public | unlisted | private
	Archived   bool
}

type GalleriesListResponse struct {
	Items  []Gallery `json:"items"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

func (g *Galleries) List(ctx context.Context, workspaceID int64, opts ListGalleriesOptions) (*GalleriesListResponse, error) {
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
	if opts.Archived {
		q.Set("archived", "true")
	}
	out := &GalleriesListResponse{}
	return out, g.c.request(ctx, g.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/galleries", nil, out, q)
}

type CreateGalleryInput struct {
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	Description   string  `json:"description,omitempty"`
	Theme         string  `json:"theme,omitempty"`
	Visibility    string  `json:"visibility,omitempty"`
	AccentColor   *string `json:"accent_color,omitempty"`
	CoverImageURL *string `json:"cover_image_url,omitempty"`
}

func (g *Galleries) Create(ctx context.Context, workspaceID int64, in CreateGalleryInput) (*Gallery, error) {
	out := &Gallery{}
	return out, g.c.request(ctx, g.c.base, "POST", "/v1/workspaces/"+itoa(workspaceID)+"/galleries", in, out, nil)
}

func (g *Galleries) Get(ctx context.Context, workspaceID, galleryID int64) (*GalleryDetail, error) {
	out := &GalleryDetail{}
	return out, g.c.request(ctx, g.c.base, "GET", "/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID), nil, out, nil)
}

type UpdateGalleryInput struct {
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	CoverImageURL  *string `json:"cover_image_url,omitempty"`
	Theme          *string `json:"theme,omitempty"`
	AccentColor    *string `json:"accent_color,omitempty"`
	Visibility     *string `json:"visibility,omitempty"`
	CustomDomainID *int64  `json:"custom_domain_id,omitempty"`
	SEOOGImageURL  *string `json:"seo_og_image_url,omitempty"`
	SEOTitle       *string `json:"seo_title,omitempty"`
	SEODescription *string `json:"seo_description,omitempty"`
	Published      *bool   `json:"published,omitempty"`
	Archived       *bool   `json:"archived,omitempty"`
}

func (g *Galleries) Update(ctx context.Context, workspaceID, galleryID int64, in UpdateGalleryInput) (*Gallery, error) {
	out := &Gallery{}
	return out, g.c.request(ctx, g.c.base, "PATCH", "/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID), in, out, nil)
}

func (g *Galleries) Delete(ctx context.Context, workspaceID, galleryID int64) error {
	return g.c.request(ctx, g.c.base, "DELETE", "/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID), nil, nil, nil)
}

type CreateGallerySectionInput struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
}

func (g *Galleries) CreateSection(ctx context.Context, workspaceID, galleryID int64, in CreateGallerySectionInput) (*GallerySection, error) {
	out := &GallerySection{}
	return out, g.c.request(ctx, g.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/sections", in, out, nil)
}

type UpdateGallerySectionInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int32  `json:"sort_order,omitempty"`
}

func (g *Galleries) UpdateSection(ctx context.Context, workspaceID, galleryID, sectionID int64, in UpdateGallerySectionInput) (*GallerySection, error) {
	out := &GallerySection{}
	return out, g.c.request(ctx, g.c.base, "PATCH",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/sections/"+itoa(sectionID), in, out, nil)
}

func (g *Galleries) DeleteSection(ctx context.Context, workspaceID, galleryID, sectionID int64) error {
	return g.c.request(ctx, g.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/sections/"+itoa(sectionID), nil, nil, nil)
}

type CreateGalleryItemInput struct {
	LinkID      int64   `json:"link_id"`
	SectionID   *int64  `json:"section_id,omitempty"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	SortOrder   int32   `json:"sort_order,omitempty"`
}

func (g *Galleries) CreateItem(ctx context.Context, workspaceID, galleryID int64, in CreateGalleryItemInput) (*GalleryItem, error) {
	out := &GalleryItem{}
	return out, g.c.request(ctx, g.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/items", in, out, nil)
}

type UpdateGalleryItemInput struct {
	SectionID   *int64  `json:"section_id,omitempty"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	SortOrder   *int32  `json:"sort_order,omitempty"`
}

func (g *Galleries) UpdateItem(ctx context.Context, workspaceID, galleryID, itemID int64, in UpdateGalleryItemInput) (*GalleryItem, error) {
	out := &GalleryItem{}
	return out, g.c.request(ctx, g.c.base, "PATCH",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/items/"+itoa(itemID), in, out, nil)
}

func (g *Galleries) DeleteItem(ctx context.Context, workspaceID, galleryID, itemID int64) error {
	return g.c.request(ctx, g.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/items/"+itoa(itemID), nil, nil, nil)
}

// ReorderGalleryItemsEntry pairs an item id with its new (section,
// sort_order). Section may be nil to drop the item back into the
// ungrouped block.
type ReorderGalleryItemsEntry struct {
	ItemID    int64  `json:"item_id"`
	SectionID *int64 `json:"section_id"`
	SortOrder int32  `json:"sort_order"`
}

type reorderItemsBody struct {
	Entries []ReorderGalleryItemsEntry `json:"entries"`
}

type ReorderGalleryItemsResponse struct {
	OK    bool `json:"ok"`
	Count int  `json:"count"`
}

func (g *Galleries) ReorderItems(ctx context.Context, workspaceID, galleryID int64, entries []ReorderGalleryItemsEntry) (*ReorderGalleryItemsResponse, error) {
	body := reorderItemsBody{Entries: entries}
	out := &ReorderGalleryItemsResponse{}
	return out, g.c.request(ctx, g.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/galleries/"+itoa(galleryID)+"/items/reorder", body, out, nil)
}

package elido

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Galleries.Create — POSTs JSON body and returns the typed gallery.
func TestGalleriesCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v1/workspaces/42/galleries"; got != want {
			t.Fatalf("path=%q want %q", got, want)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"slug":"press-kit"`) {
			t.Fatalf("body missing slug: %s", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           7,
			"workspace_id": 42,
			"slug":         "press-kit",
			"title":        "Press",
			"description":  "",
			"theme":        "classic",
			"visibility":   "unlisted",
			"created_at":   "2026-05-13T00:00:00Z",
			"updated_at":   "2026-05-13T00:00:00Z",
		})
	})
	g, err := c.Galleries.Create(context.Background(), 42, CreateGalleryInput{
		Slug:  "press-kit",
		Title: "Press",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if g.ID != 7 || g.Slug != "press-kit" {
		t.Fatalf("unexpected gallery: %+v", g)
	}
}

// ReorderItems — POSTs entries inside an envelope object.
func TestGalleriesReorderItems(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v1/workspaces/42/galleries/7/items/reorder"; got != want {
			t.Fatalf("path=%q want %q", got, want)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"entries":`) {
			t.Fatalf("body missing entries: %s", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "count": 2})
	})
	resp, err := c.Galleries.ReorderItems(context.Background(), 42, 7, []ReorderGalleryItemsEntry{
		{ItemID: 1, SortOrder: 0},
		{ItemID: 2, SortOrder: 1},
	})
	if err != nil {
		t.Fatalf("ReorderItems: %v", err)
	}
	if !resp.OK || resp.Count != 2 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

// List — query-string parameters survive the round trip.
func TestGalleriesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v1/workspaces/42/galleries"; got != want {
			t.Fatalf("path=%q want %q", got, want)
		}
		q := r.URL.Query()
		if q.Get("visibility") != "public" {
			t.Fatalf("visibility=%q", q.Get("visibility"))
		}
		if q.Get("archived") != "true" {
			t.Fatalf("archived=%q", q.Get("archived"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items":  []map[string]any{},
			"limit":  50,
			"offset": 0,
		})
	})
	resp, err := c.Galleries.List(context.Background(), 42, ListGalleriesOptions{
		Visibility: "public",
		Archived:   true,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if resp.Limit != 50 {
		t.Fatalf("limit=%d", resp.Limit)
	}
}

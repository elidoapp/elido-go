package elido

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer builds a one-handler httptest server. The caller's
// handler runs after the test inspects the request, so cases stay
// terse: each test asserts on the request and then returns a body.
func newTestServer(t *testing.T, h http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c, err := New(Options{BaseURL: srv.URL, Token: "elido_test"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c, srv
}

func TestNewRequiresBaseURL(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when BaseURL is empty")
	}
}

func TestWorkspacesListSendsBearer(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v1/workspaces"; got != want {
			t.Fatalf("path=%q want %q", got, want)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer elido_test" {
			t.Fatalf("auth=%q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{"id": 1, "name": "Acme", "role": "owner"}},
		})
	})
	resp, err := c.Workspaces.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Name != "Acme" {
		t.Fatalf("unexpected items: %+v", resp.Items)
	}
}

func TestLinksListSerialisesQuery(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("limit"); got != "20" {
			t.Fatalf("limit=%q", got)
		}
		if got := q.Get("tags"); got != "a,b" {
			t.Fatalf("tags=%q", got)
		}
		if got := q.Get("folder"); got != "42" {
			t.Fatalf("folder=%q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []any{}, "limit": 20, "offset": 0,
		})
	})
	_, err := c.Links.List(context.Background(), 7, ListLinksOptions{
		Limit:  20,
		Tags:   []string{"a", "b"},
		Folder: "42",
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestErrorEnvelopeUnwrapsMessage(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "not found"})
	})
	_, err := c.Workspaces.Get(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsAPIError(err) {
		t.Fatalf("expected APIError, got %T", err)
	}
	var ae *APIError
	_ = err
	// errors.As is what IsAPIError uses internally; pull it out for the
	// assertion. Avoiding the import dance here keeps the test concise.
	for cur := err; cur != nil; cur = nil {
		if x, ok := cur.(*APIError); ok {
			ae = x
		}
	}
	if ae == nil || ae.Status != 404 || ae.Message != "not found" {
		t.Fatalf("unexpected error: %+v", err)
	}
}

func TestDeleteReturnsNilOn204(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.Links.Delete(context.Background(), 7, 99); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestDevUserIDHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Dev-User-ID"); got != "1" {
			t.Fatalf("dev-user-id=%q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{}})
	}))
	defer srv.Close()

	c, err := New(Options{BaseURL: srv.URL, DevUserID: "1"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := c.Workspaces.List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
}

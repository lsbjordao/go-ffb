package ffb

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testClient(t *testing.T, h http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c, err := NewClient(WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestFetchFamilies(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/families" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["Fabaceae","Rubiaceae"]`))
	}))

	got, err := c.FetchFamilies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Fabaceae", "Rubiaceae"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestPathsAreEscaped(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/taxon/Mimosa%20pudica" {
			t.Fatalf("unexpected escaped path: %s", r.URL.EscapedPath())
		}
		_, _ = w.Write([]byte(`{"scientificName":"Mimosa pudica"}`))
	}))

	raw, err := c.FetchTaxon(context.Background(), "Mimosa pudica")
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(raw) {
		t.Fatal("invalid JSON")
	}
}

func TestAPIError(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	_, err := c.FetchTaxon(context.Background(), "Missing taxon")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d", apiErr.StatusCode)
	}
}

func TestFetchTaxonInto(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"scientificName":"Mimosa pudica","family":"Fabaceae"}`))
	}))
	var got struct {
		ScientificName string `json:"scientificName"`
		Family         string `json:"family"`
	}
	if err := c.FetchTaxonInto(context.Background(), "Mimosa pudica", &got); err != nil {
		t.Fatal(err)
	}
	if got.ScientificName != "Mimosa pudica" || got.Family != "Fabaceae" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

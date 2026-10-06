package ffb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestFetchTaxaPreservesOrder(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/taxon/")
		_, _ = fmt.Fprintf(w, `{"name":%q}`, name)
	}))

	taxa := []string{"A", "B", "C", "D"}
	got := c.FetchTaxa(context.Background(), taxa, 3)
	for i := range got {
		if got[i].Taxon != taxa[i] {
			t.Fatalf("index %d got %q want %q", i, got[i].Taxon, taxa[i])
		}
		if got[i].Err != nil {
			t.Fatalf("index %d: %v", i, got[i].Err)
		}
	}
}

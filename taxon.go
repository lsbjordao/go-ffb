package ffb

import (
	"context"
	"encoding/json"
	"fmt"
)

// FetchTaxon returns the exact JSON payload emitted by the Flora e Funga do Brasil API.
// Keeping the raw JSON preserves compatibility with the JavaScript package, whose
// fetchTaxon method returns response.data without imposing a schema.
func (c *Client) FetchTaxon(ctx context.Context, taxon string) (json.RawMessage, error) {
	body, err := c.get(ctx, "taxon", taxon)
	if err != nil {
		return nil, err
	}
	out := make(json.RawMessage, len(body))
	copy(out, body)
	return out, nil
}

// FetchTaxonInto decodes a taxon response into a caller-provided Go value.
// This is useful when applications maintain their own stable struct for the API payload.
func (c *Client) FetchTaxonInto(ctx context.Context, taxon string, dst any) error {
	body, err := c.FetchTaxon(ctx, taxon)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("ffb: decode taxon response: %w", err)
	}
	return nil
}

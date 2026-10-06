package ffb

import (
	"context"
	"encoding/json"
	"sync"
)

type TaxonResult struct {
	Taxon string
	Data  json.RawMessage
	Err   error
}

// FetchTaxa fetches taxon payloads concurrently while preserving the input order.
func (c *Client) FetchTaxa(ctx context.Context, taxa []string, workers int) []TaxonResult {
	if workers < 1 {
		workers = 1
	}
	if workers > len(taxa) && len(taxa) > 0 {
		workers = len(taxa)
	}
	results := make([]TaxonResult, len(taxa))
	if len(taxa) == 0 {
		return results
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				taxon := taxa[idx]
				data, err := c.FetchTaxon(ctx, taxon)
				results[idx] = TaxonResult{Taxon: taxon, Data: data, Err: err}
			}
		}()
	}

	for idx := range taxa {
		select {
		case <-ctx.Done():
			for i := idx; i < len(taxa); i++ {
				results[i] = TaxonResult{Taxon: taxa[i], Err: ctx.Err()}
			}
			close(jobs)
			wg.Wait()
			return results
		case jobs <- idx:
		}
	}
	close(jobs)
	wg.Wait()
	return results
}

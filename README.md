# go-ffb

Go SDK for the **Flora e Funga do Brasil** web services maintained by Jardim Botânico do Rio de Janeiro.

It follows the same conceptual API as [`@lsbjordao/ffb-js`](https://www.npmjs.com/package/@lsbjordao/ffb-js), while using idiomatic Go: `context.Context`, typed string slices, structured HTTP errors, standard-library HTTP/JSON support, optional decoding, and concurrent batch fetching.

## Install

```bash
go get github.com/lsbjordao/go-ffb
```

## Library usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    ffb "github.com/lsbjordao/go-ffb"
)

func main() {
    client, err := ffb.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    families, err := client.FetchFamilies(ctx)
    if err != nil {
        log.Fatal(err)
    }

    genera, err := client.FetchGeneraOfFamily(ctx, families[0])
    if err != nil {
        log.Fatal(err)
    }

    species, err := client.FetchSpeciesOfGenus(ctx, genera[0])
    if err != nil {
        log.Fatal(err)
    }

    taxon, err := client.FetchTaxon(ctx, species[0])
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(taxon))
}
```

## API mapping from `ffb-js`

| JavaScript | Go |
|---|---|
| `fetchFamilies()` | `FetchFamilies(ctx)` |
| `fetchGeneraOfFamily(family)` | `FetchGeneraOfFamily(ctx, family)` |
| `fetchSpeciesOfGenus(genus)` | `FetchSpeciesOfGenus(ctx, genus)` |
| `fetchTaxon(taxon)` | `FetchTaxon(ctx, taxon)` |

The list endpoints are decoded as `[]string`. `FetchTaxon` intentionally returns `json.RawMessage`, matching the JavaScript package's schema-free `response.data` behavior and avoiding assumptions about fields that may evolve upstream.

For typed application-level decoding:

```go
type Taxon struct {
    ScientificName string `json:"scientificName"`
    Family         string `json:"family"`
}

var taxon Taxon
err := client.FetchTaxonInto(ctx, "Mimosa pudica", &taxon)
```

## Concurrent batch fetching

```go
taxa := []string{
    "Mimosa pudica",
    "Mimosa debilis",
    "Mimosa gracilis",
}

results := client.FetchTaxa(ctx, taxa, 4)
for _, result := range results {
    if result.Err != nil {
        log.Printf("%s: %v", result.Taxon, result.Err)
        continue
    }
    fmt.Printf("%s: %s\n", result.Taxon, result.Data)
}
```

`FetchTaxa` limits concurrency to the requested worker count and preserves input order.

## CLI

Install:

```bash
go install github.com/lsbjordao/go-ffb/cmd/ffb@latest
```

Examples:

```bash
ffb families
ffb genera Fabaceae
ffb species Mimosa
ffb taxon "Mimosa pudica"
```

The CLI emits pretty-printed JSON to stdout, so it composes well with tools such as `jq`.

## Error handling

Non-2xx API responses are returned as `*ffb.APIError`:

```go
var apiErr *ffb.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.StatusCode)
}
```

Network, context cancellation, invalid JSON, and decoding errors remain distinguishable through Go's error chain.

## Design goals

- same workflow and endpoint semantics as `@lsbjordao/ffb-js`
- no runtime dependency outside the Go standard library
- safe URL construction for scientific names and other path parameters
- explicit cancellation and deadlines via `context.Context`
- preserve raw upstream taxon JSON for data pipelines
- easy integration with ETL, `goJobs`, JSON/NDJSON, Parquet or DuckDB layers

## Upstream service

Default base URL:

```text
https://servicos.jbrj.gov.br/v2/flora/
```

The base URL and HTTP client are configurable, which also makes the package easy to test.

## License

ISC

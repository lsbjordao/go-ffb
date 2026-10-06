package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ffb "github.com/lsbjordao/go-ffb"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
  ffb families
  ffb genera <family>
  ffb species <genus>
  ffb taxon <scientific-name>`)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	client, err := ffb.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ctx := context.Background()

	switch os.Args[1] {
	case "families":
		values, err := client.FetchFamilies(ctx)
		if err != nil {
			fail(err)
		}
		if err := printJSON(values); err != nil {
			fail(err)
		}
	case "genera":
		if len(os.Args) != 3 {
			usage()
			os.Exit(2)
		}
		values, err := client.FetchGeneraOfFamily(ctx, os.Args[2])
		if err != nil {
			fail(err)
		}
		if err := printJSON(values); err != nil {
			fail(err)
		}
	case "species":
		if len(os.Args) != 3 {
			usage()
			os.Exit(2)
		}
		values, err := client.FetchSpeciesOfGenus(ctx, os.Args[2])
		if err != nil {
			fail(err)
		}
		if err := printJSON(values); err != nil {
			fail(err)
		}
	case "taxon":
		if len(os.Args) != 3 {
			usage()
			os.Exit(2)
		}
		raw, err := client.FetchTaxon(ctx, os.Args[2])
		if err != nil {
			fail(err)
		}
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			fail(err)
		}
		if err := printJSON(v); err != nil {
			fail(err)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

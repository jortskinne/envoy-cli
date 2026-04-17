package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// WritePlaceholderReport writes placeholder results to w in the given format.
func WritePlaceholderReport(w io.Writer, results []PlaceholderResult, format string) error {
	switch format {
	case "json":
		return writePlaceholderJSON(w, results)
	default:
		return writePlaceholderText(w, results)
	}
}

func writePlaceholderText(w io.Writer, results []PlaceholderResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No placeholder values detected.")
		return err
	}
	_, err := fmt.Fprintf(w, "Found %d placeholder(s):\n", len(results))
	if err != nil {
		return err
	}
	for _, r := range results {
		_, err = fmt.Fprintf(w, "  %-30s = %s  (pattern: %s)\n", r.Key, r.Value, r.Pattern)
		if err != nil {
			return err
		}
	}
	return nil
}

func writePlaceholderJSON(w io.Writer, results []PlaceholderResult) error {
	out := results
	if out == nil {
		out = []PlaceholderResult{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

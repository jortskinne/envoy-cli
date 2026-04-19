package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

type resolveReportEntry struct {
	Key     string `json:"key"`
	Source  string `json:"source"`
	Missing bool   `json:"missing,omitempty"`
}

// WriteResolveReport writes a human-readable or JSON report of resolve results.
func WriteResolveReport(w io.Writer, results []ResolveResult, format string) error {
	if format == "json" {
		return writeResolveJSON(w, results)
	}
	return writeResolveText(w, results)
}

func writeResolveText(w io.Writer, results []ResolveResult) error {
	missingCount := 0
	for _, r := range results {
		if r.Missing {
			missingCount++
			fmt.Fprintf(w, "  MISSING  %s\n", r.Entry.Key)
		} else if r.Source == "os" {
			fmt.Fprintf(w, "  OS       %s\n", r.Entry.Key)
		}
	}
	if missingCount == 0 {
		fmt.Fprintln(w, "All keys resolved.")
	} else {
		fmt.Fprintf(w, "%d key(s) missing values.\n", missingCount)
	}
	return nil
}

func writeResolveJSON(w io.Writer, results []ResolveResult) error {
	var entries []resolveReportEntry
	for _, r := range results {
		if r.Source != "file" || r.Missing {
			entries = append(entries, resolveReportEntry{
				Key:     r.Entry.Key,
				Source:  r.Source,
				Missing: r.Missing,
			})
		}
	}
	if entries == nil {
		entries = []resolveReportEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

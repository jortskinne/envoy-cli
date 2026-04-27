package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// TagReport summarises the result of a Tag or ExtractTags operation.
type TagReport struct {
	Total   int               `json:"total"`
	Tagged  int               `json:"tagged"`
	Entries map[string]string `json:"entries"`
}

// BuildTagReport constructs a TagReport from an ExtractTags result and
// total entry count.
func BuildTagReport(tags map[string]string, total int) TagReport {
	return TagReport{
		Total:   total,
		Tagged:  len(tags),
		Entries: tags,
	}
}

// WriteTagReport writes a TagReport to w in the requested format
// ("text" or "json").
func WriteTagReport(w io.Writer, report TagReport, format string) error {
	switch format {
	case "json":
		return writeTagJSON(w, report)
	default:
		return writeTagText(w, report)
	}
}

func writeTagText(w io.Writer, r TagReport) error {
	fmt.Fprintf(w, "Tags: %d/%d entries tagged\n", r.Tagged, r.Total)
	if r.Tagged == 0 {
		fmt.Fprintln(w, "  (none)")
		return nil
	}
	for k, v := range r.Entries {
		fmt.Fprintf(w, "  %-30s %s\n", k, v)
	}
	return nil
}

func writeTagJSON(w io.Writer, r TagReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

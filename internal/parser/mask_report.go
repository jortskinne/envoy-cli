package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// MaskReport summarises which keys were masked in a MaskAll run.
type MaskReport struct {
	Total  int      `json:"total"`
	Masked int      `json:"masked"`
	Keys   []string `json:"masked_keys"`
}

// BuildMaskReport compares original and masked entry slices and returns a report.
func BuildMaskReport(original, masked []EnvEntry) MaskReport {
	r := MaskReport{Total: len(original)}
	for i, orig := range original {
		if i < len(masked) && masked[i].Value != orig.Value {
			r.Masked++
			r.Keys = append(r.Keys, orig.Key)
		}
	}
	if r.Keys == nil {
		r.Keys = []string{}
	}
	return r
}

// WriteMaskReport writes a human-readable or JSON mask report to w.
func WriteMaskReport(w io.Writer, r MaskReport, format string) error {
	switch format {
	case "json":
		return writeMaskJSON(w, r)
	default:
		return writeMaskText(w, r)
	}
}

func writeMaskText(w io.Writer, r MaskReport) error {
	fmt.Fprintf(w, "Mask Report: %d/%d keys masked\n", r.Masked, r.Total)
	for _, k := range r.Keys {
		fmt.Fprintf(w, "  - %s\n", k)
	}
	return nil
}

func writeMaskJSON(w io.Writer, r MaskReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

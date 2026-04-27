package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// PinReport is a serialisable summary of a pin operation.
type PinReport struct {
	Total   int         `json:"total"`
	Pinned  int         `json:"pinned"`
	Skipped int         `json:"skipped"`
	Entries []PinResult `json:"entries"`
}

// BuildPinReport aggregates results into a PinReport.
func BuildPinReport(results []PinResult) PinReport {
	r := PinReport{Total: len(results), Entries: results}
	for _, res := range results {
		if res.Pinned {
			r.Pinned++
		}
		if res.Skipped {
			r.Skipped++
		}
	}
	return r
}

// WritePinReport writes a pin report in the requested format.
func WritePinReport(w io.Writer, report PinReport, format string) error {
	switch format {
	case "json":
		return writePinJSON(w, report)
	default:
		return writePinText(w, report)
	}
}

func writePinText(w io.Writer, r PinReport) error {
	if r.Total == 0 {
		_, err := fmt.Fprintln(w, "No entries processed.")
		return err
	}
	for _, e := range r.Entries {
		status := "pinned"
		if e.Skipped {
			status = "skipped (already pinned)"
		}
		if _, err := fmt.Fprintf(w, "  %-30s %s\n", e.Key, status); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\nTotal: %d  Pinned: %d  Skipped: %d\n", r.Total, r.Pinned, r.Skipped)
	return err
}

func writePinJSON(w io.Writer, r PinReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

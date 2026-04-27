package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// FreezeReport summarises the result of a freeze operation.
type FreezeReport struct {
	Frozen  []string `json:"frozen"`
	Skipped []string `json:"skipped"`
	Total   int      `json:"total"`
}

// BuildFreezeReport compares original and updated entries to produce a report.
func BuildFreezeReport(original, updated []EnvEntry, tag string) FreezeReport {
	if tag == "" {
		tag = "@frozen"
	}
	orig := make(map[string]EnvEntry, len(original))
	for _, e := range original {
		orig[e.Key] = e
	}

	var frozen, skipped []string
	for _, e := range updated {
		wasAlready := IsFrozen(orig[e.Key], tag)
		nowFrozen := IsFrozen(e, tag)
		if nowFrozen && !wasAlready {
			frozen = append(frozen, e.Key)
		} else if !nowFrozen {
			skipped = append(skipped, e.Key)
		}
	}
	return FreezeReport{
		Frozen:  frozen,
		Skipped: skipped,
		Total:   len(updated),
	}
}

// WriteFreezeReport writes the report to w in the given format ("text" or "json").
func WriteFreezeReport(w io.Writer, r FreezeReport, format string) error {
	switch format {
	case "json":
		return writeFreezeJSON(w, r)
	default:
		return writeFreezeText(w, r)
	}
}

func writeFreezeText(w io.Writer, r FreezeReport) error {
	if len(r.Frozen) == 0 {
		_, err := fmt.Fprintln(w, "No entries frozen.")
		return err
	}
	for _, k := range r.Frozen {
		if _, err := fmt.Fprintf(w, "frozen: %s\n", k); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\n%d of %d entries frozen.\n", len(r.Frozen), r.Total)
	return err
}

func writeFreezeJSON(w io.Writer, r FreezeReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

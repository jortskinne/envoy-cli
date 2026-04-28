package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteDiffValuesReport writes a ValueDiff slice to w in the requested format.
// Supported formats: "text" (default), "json".
func WriteDiffValuesReport(diffs []ValueDiff, format string, w io.Writer) error {
	switch format {
	case "json":
		return writeDiffValuesJSON(diffs, w)
	default:
		return writeDiffValuesText(diffs, w)
	}
}

func writeDiffValuesText(diffs []ValueDiff, w io.Writer) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No value differences found.")
		return err
	}
	for _, d := range diffs {
		sensTag := ""
		if d.Sensitive {
			sensTag = " [sensitive]"
		}
		_, err := fmt.Fprintf(w, "~ %s%s\n  base:  %s\n  other: %s\n",
			d.Key, sensTag, d.BaseVal, d.OtherVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDiffValuesJSON(diffs []ValueDiff, w io.Writer) error {
	type jsonEntry struct {
		Key       string `json:"key"`
		Base      string `json:"base"`
		Other     string `json:"other"`
		Sensitive bool   `json:"sensitive"`
	}
	out := make([]jsonEntry, 0, len(diffs))
	for _, d := range diffs {
		out = append(out, jsonEntry{
			Key:       d.Key,
			Base:      d.BaseVal,
			Other:     d.OtherVal,
			Sensitive: d.Sensitive,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

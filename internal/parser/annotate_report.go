package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteAnnotateReport writes annotations to w in the given format ("text" or "json").
func WriteAnnotateReport(w io.Writer, annotations map[string]string, format string) error {
	switch format {
	case "json":
		return writeAnnotateJSON(w, annotations)
	default:
		return writeAnnotateText(w, annotations)
	}
}

func writeAnnotateText(w io.Writer, annotations map[string]string) error {
	if len(annotations) == 0 {
		_, err := fmt.Fprintln(w, "No annotations found.")
		return err
	}
	for key, comment := range annotations {
		if _, err := fmt.Fprintf(w, "%-30s # %s\n", key, comment); err != nil {
			return err
		}
	}
	return nil
}

func writeAnnotateJSON(w io.Writer, annotations map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(annotations)
}

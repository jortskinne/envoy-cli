package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteTypeCheckReport(w io.Writer, issues []TypeCheckIssue, format string) error {
	switch format {
	case "json":
		return writeTypeCheckJSON(w, issues)
	default:
		return writeTypeCheckText(w, issues)
	}
}

func writeTypeCheckText(w io.Writer, issues []TypeCheckIssue) error {
	if len(issues) == 0 {
		_, err := fmt.Fprintln(w, "typecheck: all values passed type validation")
		return err
	}
	for _, issue := range issues {
		_, err := fmt.Fprintf(w, "[FAIL] %s=%q expected=%s: %s\n", issue.Key, issue.Value, issue.Expected, issue.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTypeCheckJSON(w io.Writer, issues []TypeCheckIssue) error {
	type output struct {
		Issues []TypeCheckIssue `json:"issues"`
		Total  int              `json:"total"`
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output{Issues: issues, Total: len(issues)})
}

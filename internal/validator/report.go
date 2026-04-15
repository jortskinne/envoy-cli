package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteValidationReport writes a human-readable or JSON validation report to w.
func WriteValidationReport(w io.Writer, result ValidationResult, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return writeJSONValidationReport(w, result)
	default:
		return writeTextValidationReport(w, result)
	}
}

func writeTextValidationReport(w io.Writer, result ValidationResult) error {
	if !result.HasErrors() {
		_, err := fmt.Fprintln(w, "✓ Validation passed — no issues found.")
		return err
	}

	fmt.Fprintf(w, "✗ Validation failed — %d issue(s) found:\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(w, "  [ERROR] %s\n", e.Error())
	}
	return nil
}

func writeJSONValidationReport(w io.Writer, result ValidationResult) error {
	type jsonError struct {
		Key     string `json:"key"`
		Message string `json:"message"`
	}
	type jsonReport struct {
		Valid  bool        `json:"valid"`
		Errors []jsonError `json:"errors"`
	}

	report := jsonReport{
		Valid:  !result.HasErrors(),
		Errors: make([]jsonError, 0, len(result.Errors)),
	}
	for _, e := range result.Errors {
		report.Errors = append(report.Errors, jsonError{Key: e.Key, Message: e.Message})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

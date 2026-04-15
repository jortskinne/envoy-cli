package differ

import (
	"fmt"
	"io"
	"strings"
)

// ReportFormat controls the output format of a diff report.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatJSON ReportFormat = "json"
)

// ReportOptions configures how a diff report is rendered.
type ReportOptions struct {
	Format ReportFormat
	Color  bool
}

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// WriteReport writes a formatted diff report to the given writer.
func WriteReport(w io.Writer, result *DiffResult, opts ReportOptions) error {
	if opts.Format == FormatJSON {
		return writeJSONReport(w, result)
	}
	return writeTextReport(w, result, opts.Color)
}

func writeTextReport(w io.Writer, result *DiffResult, color bool) error {
	if !result.HasDiff() {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, e := range result.Entries {
		line := e.String()
		if color {
			line = colorize(e.Type, line)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeJSONReport(w io.Writer, result *DiffResult) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range result.Entries {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"type\": %q, \"old\": %q, \"new\": %q}",
			e.Key, string(e.Type), e.OldValue, e.NewValue,
		))
		if i < len(result.Entries)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func colorize(dt DiffType, line string) string {
	switch dt {
	case DiffAdded:
		return colorGreen + line + colorReset
	case DiffRemoved:
		return colorRed + line + colorReset
	case DiffChanged:
		return colorYellow + line + colorReset
	}
	return line
}

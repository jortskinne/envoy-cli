package audit

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteAuditReport writes the audit log to w in the given format ("text" or "json").
func WriteAuditReport(log AuditLog, format string, w io.Writer) error {
	switch format {
	case "json":
		return writeJSONAuditReport(log, w)
	default:
		return writeTextAuditReport(log, w)
	}
}

func writeTextAuditReport(log AuditLog, w io.Writer) error {
	if len(log.Entries) == 0 {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return err
	}

	fmt.Fprintf(w, "Audit Log — %s\n", log.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(w, "─────────────────────────────────────")

	for _, e := range log.Entries {
		maskedNote := ""
		if e.Masked {
			maskedNote = " [masked]"
		}
		switch e.Change {
		case ChangeAdded:
			fmt.Fprintf(w, "  + %-30s = %s%s\n", e.Key, e.NewValue, maskedNote)
		case ChangeRemoved:
			fmt.Fprintf(w, "  - %-30s = %s%s\n", e.Key, e.OldValue, maskedNote)
		case ChangeUpdated:
			fmt.Fprintf(w, "  ~ %-30s   %s → %s%s\n", e.Key, e.OldValue, e.NewValue, maskedNote)
		}
	}
	return nil
}

func writeJSONAuditReport(log AuditLog, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}

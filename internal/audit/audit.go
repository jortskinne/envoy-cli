package audit

import (
	"fmt"
	"time"

	"github.com/envoy-cli/internal/parser"
)

// ChangeType represents the kind of audit event.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// AuditEntry records a single change event for a key.
type AuditEntry struct {
	Timestamp time.Time  `json:"timestamp"`
	Key       string     `json:"key"`
	Change    ChangeType `json:"change"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
	Masked    bool       `json:"masked"`
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries    []AuditEntry `json:"entries"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// DefaultAuditOptions returns mask options used during audit.
func DefaultAuditOptions() parser.MaskOptions {
	return parser.DefaultMaskOptions()
}

// Build compares base and target entry maps and produces an AuditLog.
func Build(base, target map[string]parser.Entry, opts parser.MaskOptions) AuditLog {
	log := AuditLog{GeneratedAt: time.Now()}

	for key, baseEntry := range base {
		if targetEntry, ok := target[key]; ok {
			if baseEntry.Value != targetEntry.Value {
				old, nw := baseEntry.Value, targetEntry.Value
				masked := false
				if parser.IsSensitive(key, opts) {
					old = parser.MaskValue(old, opts)
					nw = parser.MaskValue(nw, opts)
					masked = true
				}
				log.Entries = append(log.Entries, AuditEntry{
					Timestamp: log.GeneratedAt,
					Key:       key,
					Change:    ChangeUpdated,
					OldValue:  old,
					NewValue:  nw,
					Masked:    masked,
				})
			}
		} else {
			val := baseEntry.Value
			masked := false
			if parser.IsSensitive(key, opts) {
				val = parser.MaskValue(val, opts)
				masked = true
			}
			log.Entries = append(log.Entries, AuditEntry{
				Timestamp: log.GeneratedAt,
				Key:       key,
				Change:    ChangeRemoved,
				OldValue:  val,
				Masked:    masked,
			})
		}
	}

	for key, targetEntry := range target {
		if _, ok := base[key]; !ok {
			val := targetEntry.Value
			masked := false
			if parser.IsSensitive(key, opts) {
				val = parser.MaskValue(val, opts)
				masked = true
			}
			log.Entries = append(log.Entries, AuditEntry{
				Timestamp: log.GeneratedAt,
				Key:       key,
				Change:    ChangeAdded,
				NewValue:  val,
				Masked:    masked,
			})
		}
	}

	fmt.Sprintf("") // suppress unused import if needed
	return log
}

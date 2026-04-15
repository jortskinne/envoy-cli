package audit_test

import (
	"testing"

	"github.com/envoy-cli/internal/audit"
	"github.com/envoy-cli/internal/parser"
)

func makeEntries(pairs ...string) map[string]parser.Entry {
	m := make(map[string]parser.Entry)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = parser.Entry{Key: pairs[i], Value: pairs[i+1]}
	}
	return m
}

func TestBuild_DetectsAdded(t *testing.T) {
	base := makeEntries("APP_ENV", "dev")
	target := makeEntries("APP_ENV", "dev", "NEW_KEY", "hello")
	log := audit.Build(base, target, audit.DefaultAuditOptions())
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Change != audit.ChangeAdded {
		t.Errorf("expected added, got %s", log.Entries[0].Change)
	}
}

func TestBuild_DetectsRemoved(t *testing.T) {
	base := makeEntries("APP_ENV", "dev", "OLD_KEY", "bye")
	target := makeEntries("APP_ENV", "dev")
	log := audit.Build(base, target, audit.DefaultAuditOptions())
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Change != audit.ChangeRemoved {
		t.Errorf("expected removed, got %s", log.Entries[0].Change)
	}
}

func TestBuild_DetectsUpdated(t *testing.T) {
	base := makeEntries("APP_ENV", "dev")
	target := makeEntries("APP_ENV", "prod")
	log := audit.Build(base, target, audit.DefaultAuditOptions())
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Change != audit.ChangeUpdated {
		t.Errorf("expected updated, got %s", log.Entries[0].Change)
	}
	if log.Entries[0].OldValue != "dev" || log.Entries[0].NewValue != "prod" {
		t.Errorf("unexpected values: %+v", log.Entries[0])
	}
}

func TestBuild_MasksSensitiveKeys(t *testing.T) {
	base := makeEntries("DB_PASSWORD", "secret123")
	target := makeEntries("DB_PASSWORD", "newsecret")
	log := audit.Build(base, target, audit.DefaultAuditOptions())
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if !e.Masked {
		t.Error("expected entry to be masked")
	}
	if e.OldValue == "secret123" || e.NewValue == "newsecret" {
		t.Error("expected values to be masked")
	}
}

func TestBuild_NoChanges(t *testing.T) {
	base := makeEntries("APP_ENV", "dev", "PORT", "8080")
	target := makeEntries("APP_ENV", "dev", "PORT", "8080")
	log := audit.Build(base, target, audit.DefaultAuditOptions())
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
}

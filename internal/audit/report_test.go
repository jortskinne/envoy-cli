package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/envoy-cli/internal/audit"
)

func sampleLog() audit.AuditLog {
	now := time.Now()
	return audit.AuditLog{
		GeneratedAt: now,
		Entries: []audit.AuditEntry{
			{Timestamp: now, Key: "APP_ENV", Change: audit.ChangeUpdated, OldValue: "dev", NewValue: "prod"},
			{Timestamp: now, Key: "NEW_KEY", Change: audit.ChangeAdded, NewValue: "hello"},
			{Timestamp: now, Key: "OLD_KEY", Change: audit.ChangeRemoved, OldValue: "bye"},
		},
	}
}

func TestWriteAuditReport_TextNoChanges(t *testing.T) {
	var buf bytes.Buffer
	log := audit.AuditLog{}
	if err := audit.WriteAuditReport(log, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", buf.String())
	}
}

func TestWriteAuditReport_TextWithEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := audit.WriteAuditReport(sampleLog(), "text", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Error("expected OLD_KEY in output")
	}
}

func TestWriteAuditReport_JSONValid(t *testing.T) {
	var buf bytes.Buffer
	if err := audit.WriteAuditReport(sampleLog(), "json", &buf); err != nil {
		t.Fatal(err)
	}
	var result audit.AuditLog
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}
}

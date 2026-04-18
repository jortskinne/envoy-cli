package parser

import (
	"testing"
)

func makeValidateEntries(pairs ...string) []EnvEntry {
	var entries []EnvEntry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestValidateKeys_NoConflicts(t *testing.T) {
	entries := makeValidateEntries("APP_HOST", "localhost", "APP_PORT", "8080")
	conflicts := ValidateKeys(entries, DefaultValidateKeysOptions())
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestValidateKeys_DetectsCaseConflict(t *testing.T) {
	entries := makeValidateEntries("APP_HOST", "localhost", "app_host", "other")
	conflicts := ValidateKeys(entries, DefaultValidateKeysOptions())
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].Kind != "case_conflict" {
		t.Errorf("expected case_conflict, got %s", conflicts[0].Kind)
	}
}

func TestValidateKeys_DetectsShadow(t *testing.T) {
	entries := makeValidateEntries("DB_URL", "first", "DB_URL", "second")
	conflicts := ValidateKeys(entries, DefaultValidateKeysOptions())
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].Kind != "shadow" {
		t.Errorf("expected shadow, got %s", conflicts[0].Kind)
	}
}

func TestValidateKeys_DisableCaseCheck(t *testing.T) {
	opts := DefaultValidateKeysOptions()
	opts.DetectCaseConflicts = false
	entries := makeValidateEntries("APP_HOST", "a", "app_host", "b")
	conflicts := ValidateKeys(entries, opts)
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts when case check disabled, got %d", len(conflicts))
	}
}

func TestValidateKeys_DisableShadowCheck(t *testing.T) {
	opts := DefaultValidateKeysOptions()
	opts.DetectShadows = false
	entries := makeValidateEntries("KEY", "v1", "KEY", "v2")
	conflicts := ValidateKeys(entries, opts)
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts when shadow check disabled, got %d", len(conflicts))
	}
}

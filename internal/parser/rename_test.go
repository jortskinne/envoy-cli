package parser

import (
	"testing"
)

func makeRenameEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestRename_Success(t *testing.T) {
	entries := makeRenameEntries()
	out, res, err := Rename(entries, "DB_HOST", "DATABASE_HOST", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Error("expected Renamed=true")
	}
	if out[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %s", out[0].Key)
	}
}

func TestRename_SameKey(t *testing.T) {
	entries := makeRenameEntries()
	_, res, err := Rename(entries, "DB_HOST", "DB_HOST", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected Skipped=true for identical keys")
	}
}

func TestRename_MissingKey_Error(t *testing.T) {
	entries := makeRenameEntries()
	opts := DefaultRenameOptions()
	opts.ErrorIfMissing = true
	_, _, err := Rename(entries, "MISSING_KEY", "NEW_KEY", opts)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRename_MissingKey_NoError(t *testing.T) {
	entries := makeRenameEntries()
	opts := DefaultRenameOptions()
	opts.ErrorIfMissing = false
	_, res, err := Rename(entries, "MISSING_KEY", "NEW_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected Skipped=true")
	}
}

func TestRename_TargetExists(t *testing.T) {
	entries := makeRenameEntries()
	_, res, err := Rename(entries, "DB_HOST", "DB_PORT", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected Skipped=true when target key exists")
	}
}

func TestRename_DryRun(t *testing.T) {
	entries := makeRenameEntries()
	opts := DefaultRenameOptions()
	opts.DryRun = true
	out, res, err := Rename(entries, "DB_HOST", "DATABASE_HOST", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Error("expected Renamed=true in dry run")
	}
	if out[0].Key != "DB_HOST" {
		t.Error("dry run should not modify entries")
	}
}

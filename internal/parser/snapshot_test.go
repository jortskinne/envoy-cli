package parser_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envoy-cli/internal/parser"
)

func makeSnapshotEntries() []parser.EnvEntry {
	return []parser.EnvEntry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := makeSnapshotEntries()
	if err := parser.SaveSnapshot(path, "test-label", entries); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	snap, err := parser.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if snap.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", snap.Label)
	}

	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	if snap.Timestamp.After(time.Now().Add(time.Second)) {
		t.Error("timestamp is in the future")
	}

	if len(snap.Entries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(snap.Entries))
	}

	for i, e := range entries {
		if snap.Entries[i].Key != e.Key || snap.Entries[i].Value != e.Value {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, snap.Entries[i], e)
		}
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := parser.LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)

	_, err := parser.LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSnapshotEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := makeSnapshotEntries()
	_ = parser.SaveSnapshot(path, "", entries)

	result, err := parser.SnapshotEntries(path)
	if err != nil {
		t.Fatalf("SnapshotEntries failed: %v", err)
	}

	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}

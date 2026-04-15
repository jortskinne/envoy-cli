package syncer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestSync_AddMissing(t *testing.T) {
	src := entries("APP_ENV", "production", "NEW_KEY", "hello")
	tgt := entries("APP_ENV", "staging")

	merged, res, err := Sync(src, tgt, DefaultSyncOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY added, got %v", res.Added)
	}
	if len(res.Updated) != 0 {
		t.Errorf("add-missing mode should not update, got %v", res.Updated)
	}
	if len(merged) != 2 {
		t.Errorf("expected 2 merged entries, got %d", len(merged))
	}
}

func TestSync_Overwrite(t *testing.T) {
	src := entries("APP_ENV", "production", "DB_HOST", "prod-db")
	tgt := entries("APP_ENV", "staging", "DB_HOST", "local-db")

	opts := SyncOptions{Mode: ModeOverwrite}
	_, res, err := Sync(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 2 {
		t.Errorf("expected 2 updates, got %v", res.Updated)
	}
}

func TestSync_FullRemovesExtra(t *testing.T) {
	src := entries("APP_ENV", "production")
	tgt := entries("APP_ENV", "staging", "STALE_KEY", "old")

	opts := SyncOptions{Mode: ModeFull}
	merged, res, err := Sync(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "STALE_KEY" {
		t.Errorf("expected STALE_KEY removed, got %v", res.Removed)
	}
	for _, e := range merged {
		if e.Key == "STALE_KEY" {
			t.Error("STALE_KEY should not appear in merged output")
		}
	}
}

func TestWriteEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.out")

	in := entries("FOO", "bar", "BAZ", "qux")
	if err := WriteEnvFile(path, in); err != nil {
		t.Fatalf("WriteEnvFile error: %v", err)
	}

	parsed, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(parsed) != len(in) {
		t.Errorf("expected %d entries, got %d", len(in), len(parsed))
	}
	_ = os.Remove(path)
}

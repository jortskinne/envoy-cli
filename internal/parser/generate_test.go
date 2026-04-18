package parser

import (
	"testing"
)

func TestGenerate_BasicKeys(t *testing.T) {
	keys := []string{"APP_NAME", "APP_PORT"}
	opts := DefaultGenerateOptions()
	entries, err := Generate(keys, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Value != "CHANGEME" {
			t.Errorf("expected placeholder for %s, got %q", e.Key, e.Value)
		}
	}
}

func TestGenerate_SensitiveKeysGetRandomValue(t *testing.T) {
	keys := []string{"APP_SECRET", "DB_PASSWORD"}
	opts := DefaultGenerateOptions()
	opts.Sensitive = []string{"APP_SECRET", "DB_PASSWORD"}
	entries, err := Generate(keys, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.Value == "CHANGEME" {
			t.Errorf("expected random value for sensitive key %s", e.Key)
		}
		if len(e.Value) != opts.Length {
			t.Errorf("expected length %d for %s, got %d", opts.Length, e.Key, len(e.Value))
		}
	}
}

func TestGenerate_WithPrefix(t *testing.T) {
	keys := []string{"HOST", "PORT"}
	opts := DefaultGenerateOptions()
	opts.Prefix = "APP_"
	entries, err := Generate(keys, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", entries[0].Key)
	}
	if entries[1].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", entries[1].Key)
	}
}

func TestGenerate_EmptyKeysReturnsError(t *testing.T) {
	_, err := Generate([]string{}, DefaultGenerateOptions())
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestGenerate_SkipsBlankKeys(t *testing.T) {
	keys := []string{"VALID_KEY", "   "}
	opts := DefaultGenerateOptions()
	entries, err := Generate(keys, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after skipping blank, got %d", len(entries))
	}
}

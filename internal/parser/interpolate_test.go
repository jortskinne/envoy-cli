package parser

import (
	"testing"
)

func makeInterpolateEntries(kvs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestInterpolate_ResolvesSimpleReference(t *testing.T) {
	entries := makeInterpolateEntries(
		"BASE_URL", "https://example.com",
		"API_URL", "${BASE_URL}/api",
	)
	opts := DefaultInterpolateOptions()
	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[1].Value; got != "https://example.com/api" {
		t.Errorf("expected %q, got %q", "https://example.com/api", got)
	}
}

func TestInterpolate_ResolvesUnbracedReference(t *testing.T) {
	entries := makeInterpolateEntries(
		"HOST", "localhost",
		"DSN", "postgres://$HOST/db",
	)
	opts := DefaultInterpolateOptions()
	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[1].Value; got != "postgres://localhost/db" {
		t.Errorf("expected %q, got %q", "postgres://localhost/db", got)
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	entries := makeInterpolateEntries(
		"SCHEME", "https",
		"HOST", "api.example.com",
		"BASE", "${SCHEME}://${HOST}",
		"FULL", "${BASE}/v1",
	)
	opts := DefaultInterpolateOptions()
	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[3].Value; got != "https://api.example.com/v1" {
		t.Errorf("expected %q, got %q", "https://api.example.com/v1", got)
	}
}

func TestInterpolate_LeavesUnknownReference(t *testing.T) {
	entries := makeInterpolateEntries(
		"URL", "${UNDEFINED}/path",
	)
	opts := DefaultInterpolateOptions()
	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[0].Value; got != "${UNDEFINED}/path" {
		t.Errorf("expected reference to remain, got %q", got)
	}
}

func TestInterpolate_FailOnMissing(t *testing.T) {
	entries := makeInterpolateEntries(
		"URL", "${MISSING}/path",
	)
	opts := DefaultInterpolateOptions()
	opts.FailOnMissing = true
	_, err := Interpolate(entries, opts)
	if err == nil {
		t.Error("expected error for missing variable, got nil")
	}
}

func TestInterpolate_DoesNotMutateOriginal(t *testing.T) {
	original := makeInterpolateEntries(
		"A", "hello",
		"B", "${A} world",
	)
	opts := DefaultInterpolateOptions()
	_, _ = Interpolate(original, opts)
	if original[1].Value != "${A} world" {
		t.Error("original entries were mutated")
	}
}

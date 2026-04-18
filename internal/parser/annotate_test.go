package parser

import (
	"testing"
)

func makeAnnotateEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy", Comment: ""},
		{Key: "APP_ENV", Value: "production", Comment: "existing comment"},
		{Key: "SECRET_KEY", Value: "abc123", Comment: ""},
	}
}

func TestAnnotate_AddsComment(t *testing.T) {
	entries := makeAnnotateEntries()
	annotations := map[string]string{"APP_NAME": "application name"}
	result := Annotate(entries, annotations, DefaultAnnotateOptions())
	if result[0].Comment != "application name" {
		t.Errorf("expected comment 'application name', got %q", result[0].Comment)
	}
}

func TestAnnotate_DoesNotOverwriteByDefault(t *testing.T) {
	entries := makeAnnotateEntries()
	annotations := map[string]string{"APP_ENV": "new comment"}
	result := Annotate(entries, annotations, DefaultAnnotateOptions())
	if result[1].Comment != "existing comment" {
		t.Errorf("expected original comment preserved, got %q", result[1].Comment)
	}
}

func TestAnnotate_OverwriteFlag(t *testing.T) {
	entries := makeAnnotateEntries()
	annotations := map[string]string{"APP_ENV": "replaced"}
	opts := DefaultAnnotateOptions()
	opts.Overwrite = true
	result := Annotate(entries, annotations, opts)
	if result[1].Comment != "replaced" {
		t.Errorf("expected 'replaced', got %q", result[1].Comment)
	}
}

func TestAnnotate_SkipsUnknownKeys(t *testing.T) {
	entries := makeAnnotateEntries()
	annotations := map[string]string{"UNKNOWN_KEY": "nope"}
	result := Annotate(entries, annotations, DefaultAnnotateOptions())
	for _, e := range result {
		if e.Key == "UNKNOWN_KEY" {
			t.Error("unexpected entry for UNKNOWN_KEY")
		}
	}
}

func TestExtractAnnotations_ReturnsComments(t *testing.T) {
	entries := makeAnnotateEntries()
	out := ExtractAnnotations(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(out))
	}
	if out["APP_ENV"] != "existing comment" {
		t.Errorf("unexpected comment: %q", out["APP_ENV"])
	}
}

func TestExtractAnnotations_EmptyEntries(t *testing.T) {
	out := ExtractAnnotations([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

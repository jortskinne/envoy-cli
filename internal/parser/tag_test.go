package parser

import (
	"testing"
)

func makeTagEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "LOG_LEVEL", Value: "info", Comment: "# existing"},
	}
}

func TestTag_AppliesLabel(t *testing.T) {
	entries := makeTagEntries()
	opts := DefaultTagOptions()
	opts.Tags = map[string]string{"APP_NAME": "core", "DB_PASSWORD": "sensitive"}

	result, err := Tag(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Comment != "#tag:core" {
		t.Errorf("expected #tag:core, got %q", result[0].Comment)
	}
	if result[1].Comment != "#tag:sensitive" {
		t.Errorf("expected #tag:sensitive, got %q", result[1].Comment)
	}
}

func TestTag_DoesNotOverwriteByDefault(t *testing.T) {
	entries := makeTagEntries()
	opts := DefaultTagOptions()
	opts.Tags = map[string]string{"LOG_LEVEL": "ops"}

	result, err := Tag(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[2].Comment != "# existing" {
		t.Errorf("expected existing comment preserved, got %q", result[2].Comment)
	}
}

func TestTag_OverwriteFlag(t *testing.T) {
	entries := makeTagEntries()
	opts := DefaultTagOptions()
	opts.Tags = map[string]string{"LOG_LEVEL": "ops"}
	opts.Overwrite = true

	result, err := Tag(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[2].Comment != "#tag:ops" {
		t.Errorf("expected #tag:ops, got %q", result[2].Comment)
	}
}

func TestTag_SkipsUnknownKeys(t *testing.T) {
	entries := makeTagEntries()
	opts := DefaultTagOptions()
	opts.Tags = map[string]string{"UNKNOWN_KEY": "ghost"}

	result, err := Tag(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range result {
		if e.Comment == "#tag:ghost" {
			t.Errorf("should not have tagged unknown key")
		}
	}
}

func TestExtractTags_ReturnsMappedLabels(t *testing.T) {
	entries := []EnvEntry{
		{Key: "APP_NAME", Value: "myapp", Comment: "#tag:core"},
		{Key: "DB_PASSWORD", Value: "secret", Comment: "#tag:sensitive"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
	tags := ExtractTags(entries, "")
	if tags["APP_NAME"] != "core" {
		t.Errorf("expected core, got %q", tags["APP_NAME"])
	}
	if tags["DB_PASSWORD"] != "sensitive" {
		t.Errorf("expected sensitive, got %q", tags["DB_PASSWORD"])
	}
	if _, ok := tags["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should not appear in tag map")
	}
}

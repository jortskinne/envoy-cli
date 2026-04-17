package parser

import (
	"testing"
)

func makeSanitizeEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "  APP_NAME  ", Value: "  myapp  "},
		{Key: "db-host", Value: "localhost"},
		{Key: "SECRET\x01KEY", Value: "val\x00ue"},
		{Key: "", Value: "orphan"},
	}
}

func TestSanitize_TrimWhitespace(t *testing.T) {
	entries := makeSanitizeEntries()
	opts := DefaultSanitizeOptions()
	opts.StripControlChars = false
	result := Sanitize(entries, opts)
	if result[0].Key != "APP_NAME" {
		t.Errorf("expected trimmed key, got %q", result[0].Key)
	}
	if result[0].Value != "myapp" {
		t.Errorf("expected trimmed value, got %q", result[0].Value)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	entries := makeSanitizeEntries()
	opts := DefaultSanitizeOptions()
	opts.TrimWhitespace = false
	result := Sanitize(entries, opts)
	for _, e := range result {
		for _, r := range e.Key + e.Value {
			if r < 32 || r == 127 {
				t.Errorf("found control char %d in %q=%q", r, e.Key, e.Value)
			}
		}
	}
}

func TestSanitize_RemoveEmptyKeys(t *testing.T) {
	entries := makeSanitizeEntries()
	opts := DefaultSanitizeOptions()
	opts.RemoveEmptyKeys = true
	result := Sanitize(entries, opts)
	for _, e := range result {
		if e.Key == "" {
			t.Error("expected empty key to be removed")
		}
	}
}

func TestSanitize_NormalizeKeys(t *testing.T) {
	entries := []EnvEntry{
		{Key: "db-host", Value: "localhost"},
		{Key: "app.name", Value: "myapp"},
	}
	opts := DefaultSanitizeOptions()
	opts.NormalizeKeys = true
	result := Sanitize(entries, opts)
	if result[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", result[0].Key)
	}
	if result[1].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %q", result[1].Key)
	}
}

func TestSanitize_DefaultOptionsPreservesAll(t *testing.T) {
	entries := []EnvEntry{
		{Key: "KEY", Value: "value"},
	}
	opts := DefaultSanitizeOptions()
	result := Sanitize(entries, opts)
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

package parser

import (
	"testing"
)

func makeTrimEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "  APP_NAME  ", Value: "  myapp  "},
		{Key: "DB_HOST", Value: "  localhost"},
		{Key: "SECRET_KEY", Value: "\"abc123\""},
		{Key: "  TOKEN", Value: "'bearer xyz'"},
	}
}

func TestTrim_KeysAndValues(t *testing.T) {
	entries := makeTrimEntries()
	result := Trim(entries, DefaultTrimOptions())

	if result[0].Key != "APP_NAME" {
		t.Errorf("expected trimmed key, got %q", result[0].Key)
	}
	if result[0].Value != "myapp" {
		t.Errorf("expected trimmed value, got %q", result[0].Value)
	}
	if result[1].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", result[1].Value)
	}
}

func TestTrim_QuotesNotRemovedByDefault(t *testing.T) {
	entries := makeTrimEntries()
	opts := DefaultTrimOptions()
	result := Trim(entries, opts)

	if result[2].Value != `"abc123"` {
		t.Errorf("expected quotes preserved, got %q", result[2].Value)
	}
}

func TestTrim_RemovesDoubleQuotes(t *testing.T) {
	entries := makeTrimEntries()
	opts := DefaultTrimOptions()
	opts.TrimQuotes = true
	result := Trim(entries, opts)

	if result[2].Value != "abc123" {
		t.Errorf("expected unquoted value, got %q", result[2].Value)
	}
}

func TestTrim_RemovesSingleQuotes(t *testing.T) {
	entries := makeTrimEntries()
	opts := DefaultTrimOptions()
	opts.TrimQuotes = true
	result := Trim(entries, opts)

	if result[3].Value != "bearer xyz" {
		t.Errorf("expected unquoted value, got %q", result[3].Value)
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	entries := makeTrimEntries()
	origKey := entries[0].Key
	Trim(entries, DefaultTrimOptions())
	if entries[0].Key != origKey {
		t.Error("original entries should not be mutated")
	}
}

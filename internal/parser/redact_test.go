package parser

import (
	"testing"
)

func makeRedactEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "PORT", Value: "8080"},
		{Key: "SECRET_TOKEN", Value: "tok"},
	}
}

func TestRedact_SensitiveKeysAutoRedacted(t *testing.T) {
	entries := makeRedactEntries()
	opts := DefaultRedactOptions()
	result := Redact(entries, opts)

	for _, e := range result {
		if IsSensitive(e.Key) && e.Value != "[REDACTED]" {
			t.Errorf("expected %s to be redacted, got %q", e.Key, e.Value)
		}
	}
}

func TestRedact_NonSensitiveKeysUnchanged(t *testing.T) {
	entries := makeRedactEntries()
	opts := DefaultRedactOptions()
	result := Redact(entries, opts)

	for _, e := range result {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("APP_NAME should not be redacted, got %q", e.Value)
		}
		if e.Key == "PORT" && e.Value != "8080" {
			t.Errorf("PORT should not be redacted, got %q", e.Value)
		}
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	entries := makeRedactEntries()
	opts := DefaultRedactOptions()
	opts.RedactSensitive = false
	opts.Keys = []string{"APP_NAME", "PORT"}
	result := Redact(entries, opts)

	for _, e := range result {
		switch e.Key {
		case "APP_NAME", "PORT":
			if e.Value != "[REDACTED]" {
				t.Errorf("expected %s to be redacted", e.Key)
			}
		case "DB_PASSWORD", "API_KEY", "SECRET_TOKEN":
			if e.Value == "[REDACTED]" {
				t.Errorf("%s should not be redacted when RedactSensitive=false", e.Key)
			}
		}
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	entries := makeRedactEntries()
	opts := DefaultRedactOptions()
	opts.Placeholder = "***"
	result := Redact(entries, opts)

	for _, e := range result {
		if IsSensitive(e.Key) && e.Value != "***" {
			t.Errorf("expected custom placeholder for %s", e.Key)
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	entries := makeRedactEntries()
	orig := entries[1].Value
	opts := DefaultRedactOptions()
	Redact(entries, opts)
	if entries[1].Value != orig {
		t.Error("Redact mutated original entries")
	}
}

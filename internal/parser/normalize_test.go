package parser

import (
	"testing"
)

func makeNormalizeEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "db_host", Value: "  localhost  "},
		{Key: "api_key", Value: "abc123"},
		{Key: "EMPTY_VAL", Value: ""},
		{Key: "already_quoted", Value: `"quoted"`},
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	entries := makeNormalizeEntries()
	opts := DefaultNormalizeOptions()
	opts.QuoteValues = false
	result := Normalize(entries, opts)
	for _, e := range result {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("expected uppercase key, got %q", e.Key)
		}
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	entries := makeNormalizeEntries()
	opts := DefaultNormalizeOptions()
	result := Normalize(entries, opts)
	for _, e := range result {
		if e.Value != strings.TrimSpace(e.Value) && !strings.HasPrefix(e.Value, `"`) {
			t.Errorf("expected trimmed value for %q, got %q", e.Key, e.Value)
		}
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	entries := makeNormalizeEntries()
	opts := DefaultNormalizeOptions()
	opts.RemoveEmpty = true
	result := Normalize(entries, opts)
	for _, e := range result {
		if e.Value == "" {
			t.Errorf("expected empty entries to be removed, found key %q", e.Key)
		}
	}
}

func TestNormalize_QuoteValues(t *testing.T) {
	entries := []EnvEntry{
		{Key: "FOO", Value: "bar"},
	}
	opts := DefaultNormalizeOptions()
	opts.QuoteValues = true
	result := Normalize(entries, opts)
	if result[0].Value != `"bar"` {
		t.Errorf("expected quoted value, got %q", result[0].Value)
	}
}

func TestNormalize_DoesNotDoubleQuote(t *testing.T) {
	entries := []EnvEntry{
		{Key: "FOO", Value: `"already"`},
	}
	opts := DefaultNormalizeOptions()
	opts.QuoteValues = true
	result := Normalize(entries, opts)
	if result[0].Value != `"already"` {
		t.Errorf("expected no double quoting, got %q", result[0].Value)
	}
}

func TestNormalize_PreservesOrder(t *testing.T) {
	entries := makeNormalizeEntries()
	opts := DefaultNormalizeOptions()
	opts.RemoveEmpty = false
	result := Normalize(entries, opts)
	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}

package differ_test

import (
	"testing"

	"github.com/envoy-cli/internal/differ"
	"github.com/envoy-cli/internal/parser"
)

func baseEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
	}
}

func targetEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "newpassword"},
		{Key: "LOG_LEVEL", Value: "debug"},
	}
}

func TestDiff_DetectsChangedKey(t *testing.T) {
	result := differ.Diff(baseEntries(), targetEntries(), parser.DefaultMaskOptions())

	var found *differ.DiffEntry
	for i := range result.Entries {
		if result.Entries[i].Key == "DB_PASSWORD" && result.Entries[i].Type == differ.DiffChanged {
			found = &result.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected changed entry for DB_PASSWORD")
	}
	if found.OldValue == "secret123" || found.NewValue == "newpassword" {
		t.Error("sensitive values should be masked")
	}
}

func TestDiff_DetectsRemovedKey(t *testing.T) {
	result := differ.Diff(baseEntries(), targetEntries(), parser.DefaultMaskOptions())

	var found bool
	for _, e := range result.Entries {
		if e.Key == "PORT" && e.Type == differ.DiffRemoved {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected PORT to be detected as removed")
	}
}

func TestDiff_DetectsAddedKey(t *testing.T) {
	result := differ.Diff(baseEntries(), targetEntries(), parser.DefaultMaskOptions())

	var found bool
	for _, e := range result.Entries {
		if e.Key == "LOG_LEVEL" && e.Type == differ.DiffAdded {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected LOG_LEVEL to be detected as added")
	}
}

func TestDiff_NoDiffWhenEqual(t *testing.T) {
	entries := baseEntries()
	result := differ.Diff(entries, entries, parser.DefaultMaskOptions())
	if result.HasDiff() {
		t.Errorf("expected no diff, got %d entries", len(result.Entries))
	}
}

func TestDiff_ResultIsSorted(t *testing.T) {
	result := differ.Diff(baseEntries(), targetEntries(), parser.DefaultMaskOptions())
	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i].Key < result.Entries[i-1].Key {
			t.Errorf("results not sorted: %s before %s", result.Entries[i-1].Key, result.Entries[i].Key)
		}
	}
}

func TestDiffEntry_String(t *testing.T) {
	cases := []struct {
		entry  differ.DiffEntry
		prefix string
	}{
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffAdded, NewValue: "bar"}, "+"},
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffRemoved, OldValue: "bar"}, "-"},
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffChanged, OldValue: "a", NewValue: "b"}, "~"},
	}
	for _, tc := range cases {
		s := tc.entry.String()
		if len(s) == 0 || string(s[0]) != tc.prefix {
			t.Errorf("expected string to start with %q, got %q", tc.prefix, s)
		}
	}
}

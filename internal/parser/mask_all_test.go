package parser

import (
	"testing"
)

func makeMaskAllEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestMaskAll_SensitiveKeysAutoMasked(t *testing.T) {
	entries := makeMaskAllEntries()
	opts := DefaultMaskAllOptions()
	out := MaskAll(entries, opts)

	for _, e := range out {
		if IsSensitive(e.Key) && e.Value != "****" {
			t.Errorf("expected %q to be masked, got %q", e.Key, e.Value)
		}
	}
}

func TestMaskAll_NonSensitiveKeysUnchanged(t *testing.T) {
	entries := makeMaskAllEntries()
	out := MaskAll(entries, DefaultMaskAllOptions())

	for _, e := range out {
		if !IsSensitive(e.Key) {
			original := ""
			for _, orig := range entries {
				if orig.Key == e.Key {
					original = orig.Value
				}
			}
			if e.Value != original {
				t.Errorf("key %q should not be masked", e.Key)
			}
		}
	}
}

func TestMaskAll_ExplicitKeys(t *testing.T) {
	entries := makeMaskAllEntries()
	opts := DefaultMaskAllOptions()
	opts.Keys = []string{"APP_NAME"}
	out := MaskAll(entries, opts)

	for _, e := range out {
		if e.Key == "APP_NAME" && e.Value != "****" {
			t.Errorf("expected APP_NAME masked, got %q", e.Value)
		}
		if e.Key == "PORT" && e.Value != "8080" {
			t.Errorf("PORT should not be masked")
		}
	}
}

func TestMaskAll_RevealTrailing(t *testing.T) {
	entries := []EnvEntry{{Key: "API_KEY", Value: "abcdef"}}
	opts := DefaultMaskAllOptions()
	opts.RevealTrailing = 2
	out := MaskAll(entries, opts)
	if len(out[0].Value) == 0 {
		t.Fatal("expected non-empty masked value")
	}
	if out[0].Value == "abcdef" {
		t.Error("value should have been masked")
	}
}

func TestMaskAll_DoesNotMutateOriginal(t *testing.T) {
	entries := makeMaskAllEntries()
	origVal := entries[0].Value
	MaskAll(entries, DefaultMaskAllOptions())
	if entries[0].Value != origVal {
		t.Error("MaskAll must not mutate the original slice")
	}
}

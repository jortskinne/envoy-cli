package parser

import (
	"testing"
)

func makeFreezeEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "LOG_LEVEL", Value: "info", Comment: "existing comment"},
	}
}

func TestFreeze_ExplicitKeys(t *testing.T) {
	entries := makeFreezeEntries()
	opts := DefaultFreezeOptions()
	opts.Keys = []string{"DB_HOST", "APP_ENV"}

	result, err := Freeze(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !IsFrozen(result[0], opts.Tag) {
		t.Error("expected DB_HOST to be frozen")
	}
	if IsFrozen(result[1], opts.Tag) {
		t.Error("expected DB_PASS to not be frozen")
	}
	if !IsFrozen(result[2], opts.Tag) {
		t.Error("expected APP_ENV to be frozen")
	}
}

func TestFreeze_FreezeAll(t *testing.T) {
	entries := makeFreezeEntries()
	opts := DefaultFreezeOptions()
	opts.FreezeAll = true

	result, err := Freeze(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, e := range result {
		if !IsFrozen(e, opts.Tag) {
			t.Errorf("expected %s to be frozen", e.Key)
		}
	}
}

func TestFreeze_DoesNotDuplicateTag(t *testing.T) {
	entries := []EnvEntry{
		{Key: "DB_HOST", Value: "localhost", Comment: "@frozen"},
	}
	opts := DefaultFreezeOptions()
	opts.Keys = []string{"DB_HOST"}

	result, err := Freeze(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "@frozen"
	if result[0].Comment != expected {
		t.Errorf("expected comment %q, got %q", expected, result[0].Comment)
	}
}

func TestFreeze_PreservesExistingComment(t *testing.T) {
	entries := makeFreezeEntries()
	opts := DefaultFreezeOptions()
	opts.Keys = []string{"LOG_LEVEL"}

	result, err := Freeze(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[3].Comment != "existing comment @frozen" {
		t.Errorf("unexpected comment: %q", result[3].Comment)
	}
}

func TestFreeze_NoKeysAndNoFreezeAll_ReturnsError(t *testing.T) {
	entries := makeFreezeEntries()
	opts := DefaultFreezeOptions()

	_, err := Freeze(entries, opts)
	if err == nil {
		t.Error("expected error when no keys and FreezeAll=false")
	}
}

func TestFreeze_DoesNotMutateOriginal(t *testing.T) {
	entries := makeFreezeEntries()
	opts := DefaultFreezeOptions()
	opts.FreezeAll = true

	Freeze(entries, opts)

	for _, e := range entries {
		if IsFrozen(e, opts.Tag) {
			t.Errorf("original entry %s was mutated", e.Key)
		}
	}
}

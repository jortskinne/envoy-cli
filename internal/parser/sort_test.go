package parser

import (
	"testing"
)

func makeSortEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "ZEBRA", Value: "1"},
		{Key: "AWS_SECRET", Value: "2"},
		{Key: "APP_NAME", Value: "3"},
		{Key: "DB_HOST", Value: "4"},
		{Key: "AWS_KEY", Value: "5"},
		{Key: "APP_ENV", Value: "6"},
	}
}

func TestSort_Alpha(t *testing.T) {
	entries := makeSortEntries()
	opts := DefaultSortOptions()
	out := Sort(entries, opts)

	expected := []string{"APP_ENV", "APP_NAME", "AWS_KEY", "AWS_SECRET", "DB_HOST", "ZEBRA"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	entries := makeSortEntries()
	opts := SortOptions{Order: SortAlphaDesc}
	out := Sort(entries, opts)

	expected := []string{"ZEBRA", "DB_HOST", "AWS_SECRET", "AWS_KEY", "APP_NAME", "APP_ENV"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_ByGroup(t *testing.T) {
	entries := makeSortEntries()
	opts := SortOptions{Order: SortByGroup}
	out := Sort(entries, opts)

	// Groups: APP (APP_ENV, APP_NAME), AWS (AWS_KEY, AWS_SECRET), DB (DB_HOST), ZEBRA
	if out[0].Key != "APP_ENV" || out[1].Key != "APP_NAME" {
		t.Errorf("expected APP group first, got %q %q", out[0].Key, out[1].Key)
	}
	if out[2].Key != "AWS_KEY" || out[3].Key != "AWS_SECRET" {
		t.Errorf("expected AWS group second, got %q %q", out[2].Key, out[3].Key)
	}
	if out[4].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out[4].Key)
	}
	if out[5].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA last, got %q", out[5].Key)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	entries := makeSortEntries()
	origFirst := entries[0].Key
	_ = Sort(entries, DefaultSortOptions())
	if entries[0].Key != origFirst {
		t.Errorf("original slice was mutated")
	}
}

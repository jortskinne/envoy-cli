package parser

import (
	"testing"
)

func makeEntries(pairs ...string) []EnvEntry {
	if len(pairs)%2 != 0 {
		panic("makeEntries requires an even number of arguments")
	}
	out := make([]EnvEntry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMerge_AppendsNewKeys(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	overlay := makeEntries("PORT", "8080")
	result := Merge(base, overlay, DefaultMergeOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[1].Key != "PORT" || result[1].Value != "8080" {
		t.Errorf("unexpected entry: %+v", result[1])
	}
}

func TestMerge_DoesNotOverwriteByDefault(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	overlay := makeEntries("HOST", "remotehost")
	result := Merge(base, overlay, DefaultMergeOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Value != "localhost" {
		t.Errorf("expected original value to be preserved, got %q", result[0].Value)
	}
}

func TestMerge_OverwriteFlag(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	overlay := makeEntries("HOST", "remotehost")
	opts := DefaultMergeOptions()
	opts.Overwrite = true
	result := Merge(base, overlay, opts)
	if result[0].Value != "remotehost" {
		t.Errorf("expected overwritten value, got %q", result[0].Value)
	}
}

func TestMerge_SkipsEmptyValues(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	overlay := makeEntries("PORT", "")
	result := Merge(base, overlay, DefaultMergeOptions())
	if len(result) != 1 {
		t.Errorf("expected empty overlay value to be skipped, got %d entries", len(result))
	}
}

func TestMerge_AllowsEmptyValuesWhenFlagDisabled(t *testing.T) {
	base := makeEntries("HOST", "localhost")
	overlay := makeEntries("PORT", "")
	opts := DefaultMergeOptions()
	opts.SkipEmpty = false
	result := Merge(base, overlay, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries when SkipEmpty=false, got %d", len(result))
	}
}

func TestMerge_PreservesBaseOrder(t *testing.T) {
	base := makeEntries("A", "1", "B", "2", "C", "3")
	overlay := makeEntries("D", "4")
	result := Merge(base, overlay, DefaultMergeOptions())
	keys := []string{"A", "B", "C", "D"}
	for i, k := range keys {
		if result[i].Key != k {
			t.Errorf("position %d: expected key %q, got %q", i, k, result[i].Key)
		}
	}
}

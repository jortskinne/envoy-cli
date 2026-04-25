package parser

import (
	"testing"
)

func makeCoerceEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "ENABLED", Value: "yes"},
		{Key: "DEBUG", Value: "0"},
		{Key: "PORT", Value: "  8080  "},
		{Key: "RATIO", Value: "3.14"},
		{Key: "NAME", Value: "alice"},
	}
}

func TestCoerce_InfersBool(t *testing.T) {
	entries := makeCoerceEntries()
	out, results, err := Coerce(entries, DefaultCoerceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := toKeyMap(out)
	if m["ENABLED"] != "true" {
		t.Errorf("ENABLED: want \"true\", got %q", m["ENABLED"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("DEBUG: want \"false\", got %q", m["DEBUG"])
	}

	// at least ENABLED and DEBUG should be in results
	if len(results) < 2 {
		t.Errorf("expected at least 2 coerce results, got %d", len(results))
	}
}

func TestCoerce_NormalisesInt(t *testing.T) {
	entries := makeCoerceEntries()
	out, _, err := Coerce(entries, DefaultCoerceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toKeyMap(out)
	// PORT has surrounding spaces; after int coercion it should be trimmed
	if m["PORT"] != "8080" {
		t.Errorf("PORT: want \"8080\", got %q", m["PORT"])
	}
}

func TestCoerce_NormalisesFloat(t *testing.T) {
	entries := makeCoerceEntries()
	out, _, err := Coerce(entries, DefaultCoerceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toKeyMap(out)
	if m["RATIO"] != "3.14" {
		t.Errorf("RATIO: want \"3.14\", got %q", m["RATIO"])
	}
}

func TestCoerce_LeavesStringUnchanged(t *testing.T) {
	entries := makeCoerceEntries()
	out, results, err := Coerce(entries, DefaultCoerceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toKeyMap(out)
	if m["NAME"] != "alice" {
		t.Errorf("NAME: want \"alice\", got %q", m["NAME"])
	}
	for _, r := range results {
		if r.Key == "NAME" {
			t.Errorf("NAME should not appear in coerce results")
		}
	}
}

func TestCoerce_TargetTypeBool_Error(t *testing.T) {
	entries := []EnvEntry{{Key: "NAME", Value: "alice"}}
	opts := DefaultCoerceOptions()
	opts.TargetType = "bool"
	_, _, err := Coerce(entries, opts)
	if err == nil {
		t.Error("expected error coercing \"alice\" to bool")
	}
}

func TestCoerce_RestrictsToKeys(t *testing.T) {
	entries := makeCoerceEntries()
	opts := DefaultCoerceOptions()
	opts.Keys = []string{"ENABLED"}
	out, results, err := Coerce(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "ENABLED" {
		t.Errorf("expected exactly 1 result for ENABLED, got %+v", results)
	}
	m := toKeyMap(out)
	// DEBUG should be untouched (still "0")
	if m["DEBUG"] != "0" {
		t.Errorf("DEBUG should be unchanged, got %q", m["DEBUG"])
	}
}

// toKeyMap is a small helper shared across tests.
func toKeyMap(entries []EnvEntry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

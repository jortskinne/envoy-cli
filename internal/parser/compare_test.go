package parser

import (
	"testing"
)

func makeCompareEntries(kvs ...string) []EnvEntry {
	var entries []EnvEntry
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, EnvEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestCompare_MatchingKeys(t *testing.T) {
	base := makeCompareEntries("HOST", "localhost", "PORT", "5432")
	other := makeCompareEntries("HOST", "localhost", "PORT", "5432")
	results := Compare(base, other, DefaultCompareOptions())
	for _, r := range results {
		if r.Status != "match" {
			t.Errorf("expected match for key %s, got %s", r.Key, r.Status)
		}
	}
}

func TestCompare_DetectsMismatch(t *testing.T) {
	base := makeCompareEntries("DB_HOST", "localhost")
	other := makeCompareEntries("DB_HOST", "prod.db.internal")
	results := Compare(base, other, DefaultCompareOptions())
	if len(results) != 1 || results[0].Status != "mismatch" {
		t.Fatalf("expected mismatch, got %+v", results)
	}
}

func TestCompare_BaseOnly(t *testing.T) {
	base := makeCompareEntries("SECRET", "abc")
	other := []EnvEntry{}
	results := Compare(base, other, DefaultCompareOptions())
	if len(results) != 1 || results[0].Status != "base_only" {
		t.Fatalf("expected base_only, got %+v", results)
	}
}

func TestCompare_OtherOnly(t *testing.T) {
	base := []EnvEntry{}
	other := makeCompareEntries("NEW_KEY", "value")
	results := Compare(base, other, DefaultCompareOptions())
	if len(results) != 1 || results[0].Status != "other_only" {
		t.Fatalf("expected other_only, got %+v", results)
	}
}

func TestCompare_IgnoreWhitespace(t *testing.T) {
	base := makeCompareEntries("API_URL", "  https://api.example.com  ")
	other := makeCompareEntries("API_URL", "https://api.example.com")
	opts := DefaultCompareOptions()
	opts.IgnoreWhitespace = true
	results := Compare(base, other, opts)
	if len(results) != 1 || results[0].Status != "match" {
		t.Fatalf("expected match with whitespace ignored, got %+v", results)
	}
}

func TestCompare_WhitespaceSensitive(t *testing.T) {
	base := makeCompareEntries("API_URL", "  https://api.example.com  ")
	other := makeCompareEntries("API_URL", "https://api.example.com")
	opts := DefaultCompareOptions()
	opts.IgnoreWhitespace = false
	results := Compare(base, other, opts)
	if len(results) != 1 || results[0].Status != "mismatch" {
		t.Fatalf("expected mismatch with whitespace sensitive, got %+v", results)
	}
}

func TestCompare_MixedResults(t *testing.T) {
	base := makeCompareEntries("A", "1", "B", "2", "C", "3")
	other := makeCompareEntries("A", "1", "B", "99", "D", "4")
	results := Compare(base, other, DefaultCompareOptions())
	statusMap := make(map[string]string)
	for _, r := range results {
		statusMap[r.Key] = r.Status
	}
	if statusMap["A"] != "match" {
		t.Errorf("A should be match")
	}
	if statusMap["B"] != "mismatch" {
		t.Errorf("B should be mismatch")
	}
	if statusMap["C"] != "base_only" {
		t.Errorf("C should be base_only")
	}
	if statusMap["D"] != "other_only" {
		t.Errorf("D should be other_only")
	}
}

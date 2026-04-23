package parser

import (
	"testing"
)

func makeCastEntries() []Entry {
	return []Entry{
		{Key: "ENABLE_FEATURE", Value: "true"},
		{Key: "MAX_RETRIES", Value: "003"},
		{Key: "THRESHOLD", Value: "0.75"},
		{Key: "APP_NAME", Value: "  envoy  "},
		{Key: "ACTIVE", Value: "yes"},
		{Key: "DISABLED", Value: "no"},
	}
}

func TestCast_InfersBool(t *testing.T) {
	entries := makeCastEntries()
	results, err := Cast(entries, DefaultCastOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "ENABLE_FEATURE" {
			if r.InferredType != "bool" || r.Normalized != "true" {
				t.Errorf("expected bool/true, got %s/%s", r.InferredType, r.Normalized)
			}
		}
		if r.Key == "ACTIVE" {
			if r.InferredType != "bool" || r.Normalized != "true" {
				t.Errorf("expected bool/true for yes, got %s/%s", r.InferredType, r.Normalized)
			}
		}
		if r.Key == "DISABLED" {
			if r.InferredType != "bool" || r.Normalized != "false" {
				t.Errorf("expected bool/false for no, got %s/%s", r.InferredType, r.Normalized)
			}
		}
	}
}

func TestCast_NormalizesInt(t *testing.T) {
	entries := makeCastEntries()
	results, _ := Cast(entries, DefaultCastOptions())
	for _, r := range results {
		if r.Key == "MAX_RETRIES" {
			if r.InferredType != "int" {
				t.Errorf("expected int, got %s", r.InferredType)
			}
			if r.Normalized != "3" {
				t.Errorf("expected normalized value '3', got %q", r.Normalized)
			}
		}
	}
}

func TestCast_NormalizesFloat(t *testing.T) {
	entries := makeCastEntries()
	results, _ := Cast(entries, DefaultCastOptions())
	for _, r := range results {
		if r.Key == "THRESHOLD" {
			if r.InferredType != "float" {
				t.Errorf("expected float, got %s", r.InferredType)
			}
		}
	}
}

func TestCast_FallsBackToString(t *testing.T) {
	entries := makeCastEntries()
	results, _ := Cast(entries, DefaultCastOptions())
	for _, r := range results {
		if r.Key == "APP_NAME" {
			if r.InferredType != "string" {
				t.Errorf("expected string, got %s", r.InferredType)
			}
			if r.Normalized != "envoy" {
				t.Errorf("expected trimmed value 'envoy', got %q", r.Normalized)
			}
		}
	}
}

func TestCast_FiltersByKeys(t *testing.T) {
	entries := makeCastEntries()
	opts := DefaultCastOptions()
	opts.Keys = []string{"MAX_RETRIES"}
	results, _ := Cast(entries, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "MAX_RETRIES" {
		t.Errorf("unexpected key %q", results[0].Key)
	}
}

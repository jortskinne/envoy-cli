package parser

import (
	"strings"
	"testing"
)

func makeTransformEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "  myapp  "},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DEBUG", Value: "FALSE"},
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	out := Transform(entries, opts)
	if out[0].Value != "myapp" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestTransform_Uppercase(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.Uppercase = true
	out := Transform(entries, opts)
	if out[1].Value != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out[1].Value)
	}
}

func TestTransform_Lowercase(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.Lowercase = true
	out := Transform(entries, opts)
	if out[2].Value != "false" {
		t.Errorf("expected false, got %q", out[2].Value)
	}
}

func TestTransform_FilterByKeys(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.Uppercase = true
	opts.Keys = []string{"APP_ENV"}
	out := Transform(entries, opts)
	if out[1].Value != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out[1].Value)
	}
	if out[2].Value != "FALSE" {
		t.Errorf("DEBUG should be unchanged, got %q", out[2].Value)
	}
}

func TestTransform_CustomFn(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.Custom = func(v string) string {
		return strings.ReplaceAll(v, "app", "svc")
	}
	out := Transform(entries, opts)
	if out[0].Value != "mysvc" {
		t.Errorf("expected mysvc, got %q", out[0].Value)
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.Uppercase = true
	Transform(entries, opts)
	if entries[1].Value != "production" {
		t.Error("original entries should not be mutated")
	}
}

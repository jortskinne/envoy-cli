package parser

import (
	"testing"
)

func makeEnvDiffBase() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func makeEnvDiffOther() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "remotehost"},
		{Key: "LOG_LEVEL", Value: "debug"},
		{Key: "NEW_FEATURE", Value: "enabled"},
	}
}

func TestEnvDiff_DetectsAdded(t *testing.T) {
	res := EnvDiff(makeEnvDiffBase(), makeEnvDiffOther(), DefaultEnvDiffOptions())
	if len(res.Added) != 1 || res.Added[0].Key != "NEW_FEATURE" {
		t.Errorf("expected NEW_FEATURE in Added, got %+v", res.Added)
	}
}

func TestEnvDiff_DetectsRemoved(t *testing.T) {
	res := EnvDiff(makeEnvDiffBase(), makeEnvDiffOther(), DefaultEnvDiffOptions())
	if len(res.Removed) != 1 || res.Removed[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY in Removed, got %+v", res.Removed)
	}
}

func TestEnvDiff_DetectsChanged(t *testing.T) {
	res := EnvDiff(makeEnvDiffBase(), makeEnvDiffOther(), DefaultEnvDiffOptions())
	if len(res.Changed) != 2 {
		t.Fatalf("expected 2 changed entries, got %d", len(res.Changed))
	}
	keys := map[string]bool{}
	for _, c := range res.Changed {
		keys[c.Key] = true
	}
	if !keys["DB_HOST"] || !keys["LOG_LEVEL"] {
		t.Errorf("unexpected changed keys: %+v", res.Changed)
	}
}

func TestEnvDiff_DetectsUnchanged(t *testing.T) {
	res := EnvDiff(makeEnvDiffBase(), makeEnvDiffOther(), DefaultEnvDiffOptions())
	if len(res.Unchanged) != 1 || res.Unchanged[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME unchanged, got %+v", res.Unchanged)
	}
}

func TestEnvDiff_IgnoreCase(t *testing.T) {
	base := []EnvEntry{{Key: "APP_NAME", Value: "myapp"}}
	other := []EnvEntry{{Key: "app_name", Value: "myapp"}}
	opts := DefaultEnvDiffOptions()
	opts.IgnoreCase = true
	res := EnvDiff(base, other, opts)
	if len(res.Unchanged) != 1 {
		t.Errorf("expected unchanged with IgnoreCase, got added=%v removed=%v changed=%v", res.Added, res.Removed, res.Changed)
	}
}

func TestEnvDiff_TrimValues(t *testing.T) {
	base := []EnvEntry{{Key: "HOST", Value: "localhost"}}
	other := []EnvEntry{{Key: "HOST", Value: "  localhost  "}}
	opts := DefaultEnvDiffOptions()
	opts.TrimValues = true
	res := EnvDiff(base, other, opts)
	if len(res.Unchanged) != 1 {
		t.Errorf("expected unchanged after trim, got changed=%v", res.Changed)
	}
}

package parser

import (
	"testing"
)

func makeDefaultsEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "DEBUG", Value: ""},
	}
}

func TestApplyDefaults_AddsNewKeys(t *testing.T) {
	entries := makeDefaultsEntries()
	defaults := map[string]string{"PORT": "8080", "TIMEOUT": "30"}
	result := ApplyDefaults(entries, defaults, DefaultDefaultsOptions())
	keys := make(map[string]string)
	for _, e := range result {
		keys[e.Key] = e.Value
	}
	if keys["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", keys["PORT"])
	}
	if keys["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %s", keys["TIMEOUT"])
	}
}

func TestApplyDefaults_DoesNotOverwriteByDefault(t *testing.T) {
	entries := makeDefaultsEntries()
	defaults := map[string]string{"APP_NAME": "other", "LOG_LEVEL": "debug"}
	result := ApplyDefaults(entries, defaults, DefaultDefaultsOptions())
	for _, e := range result {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME unchanged, got %s", e.Value)
		}
		if e.Key == "LOG_LEVEL" && e.Value != "info" {
			t.Errorf("expected LOG_LEVEL unchanged, got %s", e.Value)
		}
	}
}

func TestApplyDefaults_OverwriteFlag(t *testing.T) {
	entries := makeDefaultsEntries()
	defaults := map[string]string{"APP_NAME": "newapp"}
	opts := DefaultDefaultsOptions()
	opts.Overwrite = true
	result := ApplyDefaults(entries, defaults, opts)
	for _, e := range result {
		if e.Key == "APP_NAME" && e.Value != "newapp" {
			t.Errorf("expected APP_NAME=newapp, got %s", e.Value)
		}
	}
}

func TestApplyDefaults_SkipsEmptyDefault(t *testing.T) {
	entries := makeDefaultsEntries()
	defaults := map[string]string{"NEW_KEY": ""}
	result := ApplyDefaults(entries, defaults, DefaultDefaultsOptions())
	for _, e := range result {
		if e.Key == "NEW_KEY" {
			t.Error("expected NEW_KEY to be skipped due to empty default")
		}
	}
}

func TestApplyDefaults_AllowsEmptyDefaultWhenSkipEmptyFalse(t *testing.T) {
	entries := makeDefaultsEntries()
	defaults := map[string]string{"NEW_KEY": ""}
	opts := DefaultDefaultsOptions()
	opts.SkipEmpty = false
	result := ApplyDefaults(entries, defaults, opts)
	found := false
	for _, e := range result {
		if e.Key == "NEW_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to be added when SkipEmpty=false")
	}
}

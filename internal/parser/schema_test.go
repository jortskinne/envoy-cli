package parser

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempSchema(t *testing.T, schema Schema) string {
	t.Helper()
	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadSchema_Valid(t *testing.T) {
	schema := Schema{
		"APP_ENV":    {Required: true, Description: "Application environment"},
		"DB_HOST":    {Required: true, DefaultValue: "localhost"},
		"LOG_LEVEL":  {Required: false, DefaultValue: "info"},
	}
	path := writeTempSchema(t, schema)

	loaded, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 3 {
		t.Errorf("expected 3 entries, got %d", len(loaded))
	}
	if !loaded["APP_ENV"].Required {
		t.Errorf("APP_ENV should be required")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSchema_RequiredKeys(t *testing.T) {
	schema := Schema{
		"APP_ENV":  {Required: true},
		"DB_HOST":  {Required: true},
		"LOG_LEVEL": {Required: false},
	}
	keys := schema.RequiredKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 required keys, got %d", len(keys))
	}
}

func TestSchema_ApplyDefaults(t *testing.T) {
	schema := Schema{
		"LOG_LEVEL": {DefaultValue: "info"},
		"DB_HOST":   {DefaultValue: "localhost"},
	}
	entries := []EnvEntry{
		{Key: "DB_HOST", Value: "myhost"},
	}

	result := schema.ApplyDefaults(entries)

	found := false
	for _, e := range result {
		if e.Key == "LOG_LEVEL" && e.Value == "info" {
			found = true
		}
		if e.Key == "DB_HOST" && e.Value != "myhost" {
			t.Errorf("DB_HOST should not be overwritten, got %q", e.Value)
		}
	}
	if !found {
		t.Errorf("expected LOG_LEVEL default to be applied")
	}
}

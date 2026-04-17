package parser

import (
	"bytes"
	"strings"
	"testing"
)

func makePlaceholderEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "API_KEY", Value: "<your-api-key>"},
		{Key: "DB_PASS", Value: "CHANGE_ME"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "TOKEN", Value: "[replace-this]"},
		{Key: "SECRET", Value: "TODO"},
	}
}

func TestFindPlaceholders_DetectsAngleBracket(t *testing.T) {
	results, err := FindPlaceholders(makePlaceholderEntries(), DefaultPlaceholderOptions())
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]bool{}
	for _, r := range results {
		keys[r.Key] = true
	}
	if !keys["API_KEY"] {
		t.Error("expected API_KEY to be detected as placeholder")
	}
}

func TestFindPlaceholders_DetectsChangeME(t *testing.T) {
	results, err := FindPlaceholders(makePlaceholderEntries(), DefaultPlaceholderOptions())
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]bool{}
	for _, r := range results {
		keys[r.Key] = true
	}
	for _, expected := range []string{"DB_PASS", "TOKEN", "SECRET"} {
		if !keys[expected] {
			t.Errorf("expected %s to be detected as placeholder", expected)
		}
	}
}

func TestFindPlaceholders_SkipsNormalValues(t *testing.T) {
	results, err := FindPlaceholders(makePlaceholderEntries(), DefaultPlaceholderOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Key == "APP_NAME" || r.Key == "DB_HOST" {
			t.Errorf("key %s should not be flagged as placeholder", r.Key)
		}
	}
}

func TestFindPlaceholders_InvalidPattern(t *testing.T) {
	opts := DefaultPlaceholderOptions()
	opts.Patterns = []string{`[invalid`}
	_, err := FindPlaceholders(makePlaceholderEntries(), opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestWritePlaceholderReport_Text(t *testing.T) {
	results := []PlaceholderResult{
		{Key: "API_KEY", Value: "<your-api-key>", Pattern: `^<.+>$`},
	}
	var buf bytes.Buffer
	if err := WritePlaceholderReport(&buf, results, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Error("expected API_KEY in text report")
	}
}

func TestWritePlaceholderReport_JSONEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := WritePlaceholderReport(&buf, nil, "json"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Error("expected JSON array in output")
	}
}

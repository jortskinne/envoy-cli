package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseFile_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDATABASE_URL=postgres://localhost/db\n")
	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(envFile.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(envFile.Entries))
	}
	if envFile.Entries[0].Key != "APP_ENV" || envFile.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", envFile.Entries[0])
	}
}

func TestParseFile_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nKEY=value\n")
	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(envFile.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(envFile.Entries))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envFile.Entries[0].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", envFile.Entries[0].Value)
	}
}

func TestParseFile_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := envFile.Entries[0]
	if entry.Value != "8080" {
		t.Errorf("expected value 8080, got %q", entry.Value)
	}
	if entry.Comment != "http port" {
		t.Errorf("expected comment 'http port', got %q", entry.Comment)
	}
}

func TestToMap(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	envFile, _ := ParseFile(path)
	m := envFile.ToMap()
	if m["FOO"].Value != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"].Value)
	}
	if m["BAZ"].Value != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"].Value)
	}
}

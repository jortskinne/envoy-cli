package parser

import (
	"encoding/json"
	"os"
	"testing"
)

func writeImportTempFile(t *testing.T, content, ext string) string {
	t.Helper()
	f, err := os.CreateTemp("", "import-*"+ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func makeImportBase() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "dev"},
	}
}

func TestImport_DotEnvAddsNewKeys(t *testing.T) {
	path := writeImportTempFile(t, "NEW_KEY=hello\nANOTHER=world\n", ".env")
	result, err := Import(makeImportBase(), path, DefaultImportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := toKeyMap(result)
	if keys["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", keys["NEW_KEY"])
	}
	if keys["ANOTHER"] != "world" {
		t.Errorf("expected ANOTHER=world, got %q", keys["ANOTHER"])
	}
}

func TestImport_DoesNotOverwriteByDefault(t *testing.T) {
	path := writeImportTempFile(t, "APP_NAME=override\n", ".env")
	result, err := Import(makeImportBase(), path, DefaultImportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := toKeyMap(result)
	if keys["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to remain 'myapp', got %q", keys["APP_NAME"])
	}
}

func TestImport_OverwriteFlag(t *testing.T) {
	path := writeImportTempFile(t, "APP_NAME=override\n", ".env")
	opts := DefaultImportOptions()
	opts.Overwrite = true
	result, err := Import(makeImportBase(), path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := toKeyMap(result)
	if keys["APP_NAME"] != "override" {
		t.Errorf("expected APP_NAME=override, got %q", keys["APP_NAME"])
	}
}

func TestImport_JSONFormat(t *testing.T) {
	m := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	data, _ := json.Marshal(m)
	path := writeImportTempFile(t, string(data), ".json")
	result, err := Import(makeImportBase(), path, DefaultImportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := toKeyMap(result)
	if keys["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", keys["DB_HOST"])
	}
}

func TestImport_SkipInvalidLines(t *testing.T) {
	path := writeImportTempFile(t, "VALID=yes\n!!!bad line\nOTHER=ok\n", ".env")
	opts := DefaultImportOptions()
	opts.SkipInvalid = true
	result, err := Import(makeImportBase(), path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := toKeyMap(result)
	if keys["VALID"] != "yes" {
		t.Errorf("expected VALID=yes")
	}
}

func TestImport_MissingFile(t *testing.T) {
	_, err := Import(makeImportBase(), "/nonexistent/path.env", DefaultImportOptions())
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func toKeyMap(entries []EnvEntry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

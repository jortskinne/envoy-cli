package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTemplateEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	entries := makeTemplateEntries()
	out, err := RenderTemplate("Hello from {{APP_NAME}} on port {{PORT}}", entries, DefaultTemplateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "Hello from envoy on port 8080" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderTemplate_UnresolvedStrictMode(t *testing.T) {
	_, err := RenderTemplate("host={{MISSING}}", makeTemplateEntries(), DefaultTemplateOptions())
	if err == nil {
		t.Fatal("expected error for unresolved placeholder")
	}
}

func TestRenderTemplate_UnresolvedLenient(t *testing.T) {
	opts := DefaultTemplateOptions()
	opts.StrictMode = false
	out, err := RenderTemplate("host={{MISSING}}", makeTemplateEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "host={{MISSING}}" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderTemplate_FallsBackToOSEnv(t *testing.T) {
	t.Setenv("OS_KEY", "from-os")
	out, err := RenderTemplate("val={{OS_KEY}}", []EnvEntry{}, DefaultTemplateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "val=from-os" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderTemplateFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	_ = os.WriteFile(p, []byte("db={{DB_HOST}}"), 0600)
	out, err := RenderTemplateFile(p, makeTemplateEntries(), DefaultTemplateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "db=localhost" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderTemplateFile_MissingFile(t *testing.T) {
	_, err := RenderTemplateFile("/no/such/file.txt", makeTemplateEntries(), DefaultTemplateOptions())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

package parser

import (
	"strings"
	"testing"
)

func makeConvertEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestConvert_DotEnv(t *testing.T) {
	entries := makeConvertEntries()
	out, err := Convert(entries, DefaultConvertOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME=envoy in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD=s3cr3t") {
		t.Errorf("expected unmasked password in dotenv output")
	}
}

func TestConvert_DotEnvMasked(t *testing.T) {
	entries := makeConvertEntries()
	opts := DefaultConvertOptions()
	opts.MaskSecrets = true
	out, err := Convert(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected password to be masked, got:\n%s", out)
	}
}

func TestConvert_Export(t *testing.T) {
	entries := makeConvertEntries()
	opts := DefaultConvertOptions()
	opts.Format = FormatExport
	out, err := Convert(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestConvert_JSON(t *testing.T) {
	entries := makeConvertEntries()
	opts := DefaultConvertOptions()
	opts.Format = FormatJSON
	out, err := Convert(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"APP_NAME\"") {
		t.Errorf("expected JSON key APP_NAME, got:\n%s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
}

func TestConvert_YAML(t *testing.T) {
	entries := makeConvertEntries()
	opts := DefaultConvertOptions()
	opts.Format = FormatYAML
	out, err := Convert(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME:") {
		t.Errorf("expected YAML key APP_NAME:, got:\n%s", out)
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	entries := makeConvertEntries()
	opts := DefaultConvertOptions()
	opts.Format = OutputFormat("toml")
	_, err := Convert(entries, opts)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

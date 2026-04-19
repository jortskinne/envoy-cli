package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildMaskReport_CountsMasked(t *testing.T) {
	orig := []EnvEntry{
		{Key: "API_KEY", Value: "secret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	masked := []EnvEntry{
		{Key: "API_KEY", Value: "****"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	r := BuildMaskReport(orig, masked)
	if r.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Total)
	}
	if r.Masked != 1 {
		t.Errorf("expected masked 1, got %d", r.Masked)
	}
	if len(r.Keys) != 1 || r.Keys[0] != "API_KEY" {
		t.Errorf("unexpected masked keys: %v", r.Keys)
	}
}

func TestBuildMaskReport_NoMasked(t *testing.T) {
	orig := []EnvEntry{{Key: "PORT", Value: "8080"}}
	masked := []EnvEntry{{Key: "PORT", Value: "8080"}}
	r := BuildMaskReport(orig, masked)
	if r.Masked != 0 {
		t.Error("expected 0 masked")
	}
	if len(r.Keys) != 0 {
		t.Error("expected empty keys slice")
	}
}

func TestWriteMaskReport_TextFormat(t *testing.T) {
	r := MaskReport{Total: 3, Masked: 2, Keys: []string{"API_KEY", "DB_PASSWORD"}}
	var buf bytes.Buffer
	if err := WriteMaskReport(&buf, r, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2/3") {
		t.Errorf("expected summary in output, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
}

func TestWriteMaskReport_JSONFormat(t *testing.T) {
	r := MaskReport{Total: 1, Masked: 1, Keys: []string{"SECRET"}}
	var buf bytes.Buffer
	if err := WriteMaskReport(&buf, r, "json"); err != nil {
		t.Fatal(err)
	}
	var out MaskReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Masked != 1 {
		t.Errorf("expected masked=1, got %d", out.Masked)
	}
}

package validator

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteValidationReport_TextNoErrors(t *testing.T) {
	var buf bytes.Buffer
	err := WriteValidationReport(&buf, ValidationResult{}, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %s", buf.String())
	}
}

func TestWriteValidationReport_TextWithErrors(t *testing.T) {
	result := ValidationResult{}
	result.Add("SECRET", "required key is missing")
	result.Add("lower_key", "key must be UPPER_SNAKE_CASE")

	var buf bytes.Buffer
	err := WriteValidationReport(&buf, result, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output")
	}
	if !strings.Contains(out, "lower_key") {
		t.Errorf("expected lower_key in output")
	}
}

func TestWriteValidationReport_JSONValid(t *testing.T) {
	var buf bytes.Buffer
	err := WriteValidationReport(&buf, ValidationResult{}, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["valid"] != true {
		t.Errorf("expected valid=true, got %v", out["valid"])
	}
}

func TestWriteValidationReport_JSONWithErrors(t *testing.T) {
	result := ValidationResult{}
	result.Add("DB_PASS", "required key is missing")

	var buf bytes.Buffer
	err := WriteValidationReport(&buf, result, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["valid"] != false {
		t.Errorf("expected valid=false")
	}
	errs, ok := out["errors"].([]interface{})
	if !ok || len(errs) != 1 {
		t.Errorf("expected 1 error in JSON output")
	}
}

func TestWriteValidationReport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := WriteValidationReport(&buf, ValidationResult{}, "xml")
	if err == nil {
		t.Errorf("expected error for unsupported format 'xml', got nil")
	}
}

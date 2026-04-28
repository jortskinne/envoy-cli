package parser

import (
	"bytes"
	"strings"
	"testing"
)

func makeDiffValuesBase() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
	}
}

func makeDiffValuesOther() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "newpassword"},
		{Key: "PORT", Value: "9090"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestDiffValues_DetectsChangedValues(t *testing.T) {
	diffs := DiffValues(makeDiffValuesBase(), makeDiffValuesOther(), DefaultDiffValuesOptions())
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	keys := map[string]bool{}
	for _, d := range diffs {
		keys[d.Key] = true
	}
	if !keys["DB_PASSWORD"] || !keys["PORT"] {
		t.Errorf("expected DB_PASSWORD and PORT in diffs, got %v", keys)
	}
}

func TestDiffValues_IgnoreCase(t *testing.T) {
	base := []EnvEntry{{Key: "MODE", Value: "Production"}}
	other := []EnvEntry{{Key: "MODE", Value: "production"}}
	opts := DefaultDiffValuesOptions()
	opts.IgnoreCase = true
	diffs := DiffValues(base, other, opts)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs with IgnoreCase, got %d", len(diffs))
	}
}

func TestDiffValues_MaskSensitive(t *testing.T) {
	opts := DefaultDiffValuesOptions()
	opts.MaskSensitive = true
	diffs := DiffValues(makeDiffValuesBase(), makeDiffValuesOther(), opts)
	for _, d := range diffs {
		if d.Sensitive {
			if d.BaseVal != "***" || d.OtherVal != "***" {
				t.Errorf("expected masked values for %s", d.Key)
			}
		}
	}
}

func TestDiffValues_SkipsMissingKeys(t *testing.T) {
	base := []EnvEntry{{Key: "ONLY_BASE", Value: "x"}}
	other := []EnvEntry{{Key: "ONLY_OTHER", Value: "y"}}
	diffs := DiffValues(base, other, DefaultDiffValuesOptions())
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs for non-overlapping keys, got %d", len(diffs))
	}
}

func TestWriteDiffValuesReport_TextNoDiffs(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteDiffValuesReport(nil, "text", &buf)
	if !strings.Contains(buf.String(), "No value differences") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteDiffValuesReport_JSONFormat(t *testing.T) {
	diffs := []ValueDiff{{Key: "PORT", BaseVal: "8080", OtherVal: "9090"}}
	var buf bytes.Buffer
	_ = WriteDiffValuesReport(diffs, "json", &buf)
	if !strings.Contains(buf.String(), "\"key\"") {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

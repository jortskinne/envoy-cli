package parser

import (
	"testing"
)

func makeLintEntries(pairs ...string) []EnvEntry {
	var entries []EnvEntry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestLint_NoIssuesForValidEntries(t *testing.T) {
	entries := makeLintEntries("APP_NAME", "myapp", "DB_HOST", "localhost")
	issues := Lint(entries, DefaultLintOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_DetectsLowercaseKey(t *testing.T) {
	entries := makeLintEntries("app_name", "myapp")
	issues := Lint(entries, DefaultLintOptions())
	if len(issues) == 0 {
		t.Fatal("expected an issue for lowercase key")
	}
	if issues[0].Severity != LintError {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestLint_DetectsEmptyValue(t *testing.T) {
	entries := makeLintEntries("API_KEY", "")
	opts := DefaultLintOptions()
	issues := Lint(entries, opts)
	found := false
	for _, i := range issues {
		if i.Key == "API_KEY" && i.Severity == LintWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for empty value on API_KEY")
	}
}

func TestLint_DetectsDuplicateKey(t *testing.T) {
	entries := []EnvEntry{
		{Key: "DB_HOST", Value: "host1"},
		{Key: "DB_HOST", Value: "host2"},
	}
	issues := Lint(entries, DefaultLintOptions())
	found := false
	for _, i := range issues {
		if i.Key == "DB_HOST" && i.Severity == LintWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for duplicate key DB_HOST")
	}
}

func TestLint_RespectDisabledRules(t *testing.T) {
	entries := makeLintEntries("lower_key", "")
	opts := LintOptions{
		DisallowEmptyValues: false,
		EnforceUpperSnake:   false,
		WarnDuplicateKeys:   false,
	}
	issues := Lint(entries, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues with all rules disabled, got %d", len(issues))
	}
}

func TestLint_MultipleIssuesSameEntry(t *testing.T) {
	entries := makeLintEntries("bad-key", "")
	issues := Lint(entries, DefaultLintOptions())
	// expect at least two issues: non-upper-snake + empty value
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(issues))
	}
}

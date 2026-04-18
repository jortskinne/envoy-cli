package parser

import (
	"testing"
)

func makeTypecheckEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
		{Key: "API_URL", Value: "https://example.com"},
		{Key: "ADMIN_EMAIL", Value: "admin@example.com"},
		{Key: "BAD_INT", Value: "notanint"},
		{Key: "BAD_BOOL", Value: "yes"},
		{Key: "BAD_URL", Value: "ftp://bad"},
		{Key: "BAD_EMAIL", Value: "notanemail"},
		{Key: "CODE", Value: "ABC-123"},
	}
}

func TestTypeCheck_ValidInt(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "PORT", Type: "int"}},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestTypeCheck_InvalidInt(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "BAD_INT", Type: "int"}},
	})
	if len(issues) != 1 || issues[0].Key != "BAD_INT" {
		t.Fatalf("expected 1 issue for BAD_INT, got %+v", issues)
	}
}

func TestTypeCheck_ValidBool(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "DEBUG", Type: "bool"}},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestTypeCheck_InvalidBool(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "BAD_BOOL", Type: "bool"}},
	})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestTypeCheck_ValidURL(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "API_URL", Type: "url"}},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestTypeCheck_InvalidEmail(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "BAD_EMAIL", Type: "email"}},
	})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestTypeCheck_RegexMatch(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "CODE", Type: "regex", Pattern: `^[A-Z]+-\d+$`}},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestTypeCheck_SkipsMissingKey(t *testing.T) {
	issues := TypeCheck(makeTypecheckEntries(), TypeCheckOptions{
		Rules: []TypeRule{{Key: "NONEXISTENT", Type: "int"}},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues for missing key, got %d", len(issues))
	}
}

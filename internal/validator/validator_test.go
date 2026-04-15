package validator

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	result := Validate(entries("FOO", "bar"), ValidateOptions{
		RequiredKeys: []string{"FOO", "MISSING_KEY"},
	})
	if !result.HasErrors() {
		t.Fatal("expected errors, got none")
	}
	if len(result.Errors) != 1 || result.Errors[0].Key != "MISSING_KEY" {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestValidate_RequiredKeyEmpty(t *testing.T) {
	result := Validate(entries("SECRET", ""), ValidateOptions{
		RequiredKeys: []string{"SECRET"},
	})
	if !result.HasErrors() {
		t.Fatal("expected error for empty required key")
	}
	if result.Errors[0].Key != "SECRET" {
		t.Errorf("expected error on SECRET, got %v", result.Errors)
	}
}

func TestValidate_DisallowEmptyValues(t *testing.T) {
	result := Validate(entries("FOO", "ok", "BAR", ""), ValidateOptions{
		DisallowEmptyValues: true,
	})
	if len(result.Errors) != 1 || result.Errors[0].Key != "BAR" {
		t.Errorf("expected single error on BAR, got %v", result.Errors)
	}
}

func TestValidate_EnforceUpperSnake_Valid(t *testing.T) {
	result := Validate(entries("VALID_KEY", "v", "ALSO_VALID_123", "v"), ValidateOptions{
		EnforceUpperSnake: true,
	})
	if result.HasErrors() {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
}

func TestValidate_EnforceUpperSnake_Invalid(t *testing.T) {
	result := Validate(entries("camelCase", "v", "lower_key", "v"), ValidateOptions{
		EnforceUpperSnake: true,
	})
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidate_NoErrors(t *testing.T) {
	result := Validate(entries("APP_ENV", "production", "PORT", "8080"), ValidateOptions{
		RequiredKeys:        []string{"APP_ENV", "PORT"},
		DisallowEmptyValues: true,
		EnforceUpperSnake:   true,
	})
	if result.HasErrors() {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
}

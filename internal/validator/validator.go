package validator

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Add(key, message string) {
	r.Errors = append(r.Errors, ValidationError{Key: key, Message: message})
}

// ValidateOptions controls which checks are performed.
type ValidateOptions struct {
	// RequiredKeys lists keys that must be present and non-empty.
	RequiredKeys []string
	// DisallowEmptyValues fails any key whose value is empty.
	DisallowEmptyValues bool
	// AllowedKeyPattern is a simple prefix/suffix check (e.g. must be UPPER_SNAKE).
	EnforceUpperSnake bool
}

// Validate runs all configured checks against the provided entries.
func Validate(entries []parser.Entry, opts ValidateOptions) ValidationResult {
	result := ValidationResult{}

	keySet := make(map[string]parser.Entry, len(entries))
	for _, e := range entries {
		keySet[e.Key] = e
	}

	// Check required keys.
	for _, req := range opts.RequiredKeys {
		e, ok := keySet[req]
		if !ok {
			result.Add(req, "required key is missing")
			continue
		}
		if strings.TrimSpace(e.Value) == "" {
			result.Add(req, "required key is present but empty")
		}
	}

	for _, e := range entries {
		// Disallow empty values.
		if opts.DisallowEmptyValues && strings.TrimSpace(e.Value) == "" {
			result.Add(e.Key, "value must not be empty")
		}

		// Enforce UPPER_SNAKE_CASE.
		if opts.EnforceUpperSnake && !isUpperSnake(e.Key) {
			result.Add(e.Key, "key must be UPPER_SNAKE_CASE")
		}
	}

	return result
}

// isUpperSnake returns true if the key contains only uppercase letters, digits, and underscores.
func isUpperSnake(key string) bool {
	if key == "" {
		return false
	}
	for _, ch := range key {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}

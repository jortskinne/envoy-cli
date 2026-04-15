package parser

import "fmt"

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintError   LintSeverity = "error"
	LintWarning LintSeverity = "warning"
)

// LintIssue describes a single lint finding for an env entry.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintOptions controls which lint rules are applied.
type LintOptions struct {
	DisallowEmptyValues  bool
	EnforceUpperSnake    bool
	DisallowLeadingSpace bool
	WarnDuplicateKeys    bool
}

// DefaultLintOptions returns a LintOptions with sensible defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		DisallowEmptyValues:  true,
		EnforceUpperSnake:    true,
		DisallowLeadingSpace: true,
		WarnDuplicateKeys:    true,
	}
}

// Lint runs all configured lint rules against the provided entries and returns
// a slice of LintIssue. An empty slice means no issues were found.
func Lint(entries []EnvEntry, opts LintOptions) []LintIssue {
	var issues []LintIssue
	seen := make(map[string]int)

	for _, e := range entries {
		seen[e.Key]++

		if opts.EnforceUpperSnake && !isUpperSnakeKey(e.Key) {
			issues = append(issues, LintIssue{
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q is not UPPER_SNAKE_CASE", e.Key),
				Severity: LintError,
			})
		}

		if opts.DisallowEmptyValues && e.Value == "" {
			issues = append(issues, LintIssue{
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q has an empty value", e.Key),
				Severity: LintWarning,
			})
		}

		if opts.DisallowLeadingSpace && len(e.Key) > 0 && e.Key[0] == ' ' {
			issues = append(issues, LintIssue{
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q has a leading space", e.Key),
				Severity: LintError,
			})
		}
	}

	if opts.WarnDuplicateKeys {
		for key, count := range seen {
			if count > 1 {
				issues = append(issues, LintIssue{
					Key:      key,
					Message:  fmt.Sprintf("key %q is defined %d times", key, count),
					Severity: LintWarning,
				})
			}
		}
	}

	return issues
}

// isUpperSnakeKey returns true if s consists only of uppercase letters, digits,
// and underscores, and is non-empty.
func isUpperSnakeKey(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

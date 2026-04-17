package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// PlaceholderOptions controls how placeholder detection works.
type PlaceholderOptions struct {
	// Patterns are regex patterns that identify placeholder values.
	Patterns []string
	// IncludeSensitive includes sensitive keys in the result.
	IncludeSensitive bool
}

// PlaceholderResult holds a key found to have a placeholder value.
type PlaceholderResult struct {
	Key     string
	Value   string
	Pattern string
}

// DefaultPlaceholderOptions returns sensible defaults.
func DefaultPlaceholderOptions() PlaceholderOptions {
	return PlaceholderOptions{
		Patterns: []string{
			`^<.+>$`,
			`^\[.+\]$`,
			`^CHANGE_ME$`,
			`^TODO$`,
			`^REPLACE_ME$`,
			`^your[_-].+`,
		},
		IncludeSensitive: true,
	}
}

// FindPlaceholders scans entries and returns those whose values match placeholder patterns.
func FindPlaceholders(entries []Entry, opts PlaceholderOptions) ([]PlaceholderResult, error) {
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid placeholder pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	var results []PlaceholderResult
	for _, e := range entries {
		if e.Comment {
			continue
		}
		val := strings.TrimSpace(e.Value)
		for _, re := range compiled {
			if re.MatchString(val) {
				results = append(results, PlaceholderResult{
					Key:     e.Key,
					Value:   val,
					Pattern: re.String(),
				})
				break
			}
		}
	}
	return results, nil
}

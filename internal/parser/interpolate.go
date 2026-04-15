package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolationPattern matches ${VAR} and $VAR style references
var interpolationPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// InterpolateOptions controls interpolation behaviour.
type InterpolateOptions struct {
	// FailOnMissing returns an error if a referenced variable is not found.
	FailOnMissing bool
	// MaxDepth limits recursive interpolation passes to prevent cycles.
	MaxDepth int
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		FailOnMissing: false,
		MaxDepth:      10,
	}
}

// Interpolate resolves variable references within entry values using the
// provided entries as the lookup table. It mutates a copy of the slice and
// returns the resolved entries.
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	// Build a lookup map from the current (possibly partially resolved) values.
	resolved := make([]Entry, len(entries))
	copy(resolved, entries)

	for pass := 0; pass < opts.MaxDepth; pass++ {
		lookup := buildLookup(resolved)
		changed := false

		for i, e := range resolved {
			newVal, err := expandValue(e.Value, lookup, opts.FailOnMissing)
			if err != nil {
				return nil, fmt.Errorf("interpolation error for key %q: %w", e.Key, err)
			}
			if newVal != e.Value {
				resolved[i].Value = newVal
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	return resolved, nil
}

func buildLookup(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func expandValue(value string, lookup map[string]string, failOnMissing bool) (string, error) {
	var expandErr error
	result := interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		// Extract variable name from either ${VAR} or $VAR form.
		submatches := interpolationPattern.FindStringSubmatch(match)
		varName := submatches[1]
		if varName == "" {
			varName = submatches[2]
		}
		if val, ok := lookup[varName]; ok {
			return val
		}
		if failOnMissing {
			expandErr = fmt.Errorf("undefined variable %q", varName)
			return match
		}
		return strings.TrimSpace(match) // leave unresolved references as-is
	})
	return result, expandErr
}

package parser

import (
	"regexp"
	"strings"
)

// GrepOptions controls how Grep filters entries.
type GrepOptions struct {
	Pattern     string
	SearchKeys  bool
	SearchValues bool
	CaseSensitive bool
	Invert      bool
}

// DefaultGrepOptions returns sensible defaults.
func DefaultGrepOptions() GrepOptions {
	return GrepOptions{
		SearchKeys:   true,
		SearchValues: true,
		CaseSensitive: false,
		Invert:       false,
	}
}

// Grep returns entries whose key or value matches the given pattern.
func Grep(entries []Entry, opts GrepOptions) ([]Entry, error) {
	pattern := opts.Pattern
	if !opts.CaseSensitive {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var result []Entry
	for _, e := range entries {
		matched := false
		if opts.SearchKeys && re.MatchString(e.Key) {
			matched = true
		}
		if opts.SearchValues && re.MatchString(e.Value) {
			matched = true
		}
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result = append(result, e)
		}
	}
	_ = strings.ToLower // suppress unused import
	return result, nil
}

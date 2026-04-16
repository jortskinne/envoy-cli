package parser

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TrimOptions controls how trimming is applied to env entries.
type TrimOptions struct {
	TrimKeys    bool
	TrimValues  bool
	TrimQuotes  bool
}

// DefaultTrimOptions returns sensible defaults.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		TrimKeys:   true,
		TrimValues: true,
		TrimQuotes: false,
	}
}

// Trim cleans up whitespace (and optionally quotes) from env entries.
func Trim(entries []EnvEntry, opts TrimOptions) []EnvEntry {
	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.TrimQuotes {
			e.Value = trimQuotes(e.Value)
		}
		result = append(result, e)
	}
	return result
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// EnvEntry is re-exported alias guard — actual type lives in env_parser.go.
var _ = parser.EnvEntry{}

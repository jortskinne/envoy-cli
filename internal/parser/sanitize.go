package parser

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// SanitizeOptions controls sanitization behavior.
type SanitizeOptions struct {
	StripControlChars bool
	TrimWhitespace    bool
	RemoveEmptyKeys   bool
	NormalizeKeys     bool // uppercase + underscores
}

// DefaultSanitizeOptions returns sensible defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		StripControlChars: true,
		TrimWhitespace:    true,
		RemoveEmptyKeys:   false,
		NormalizeKeys:     false,
	}
}

// Sanitize cleans up env entries according to the provided options.
func Sanitize(entries []EnvEntry, opts SanitizeOptions) []EnvEntry {
	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if opts.TrimWhitespace {
			e.Key = strings.TrimSpace(e.Key)
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.StripControlChars {
			e.Key = stripControl(e.Key)
			e.Value = stripControl(e.Value)
		}
		if opts.NormalizeKeys {
			e.Key = normalizeKey(e.Key)
		}
		if opts.RemoveEmptyKeys && e.Key == "" {
			continue
		}
		result = append(result, e)
	}
	return result
}

func stripControl(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 32 && r != 127 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeKey(s string) string {
	s = strings.ToUpper(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// silence unused import
var _ = parser.EnvEntry{}

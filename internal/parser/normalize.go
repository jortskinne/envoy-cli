package parser

import (
	"strings"
)

// NormalizeOptions controls normalization behavior.
type NormalizeOptions struct {
	UppercaseKeys   bool
	TrimValues      bool
	RemoveEmpty     bool
	QuoteValues     bool
}

// DefaultNormalizeOptions returns sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys: true,
		TrimValues:    true,
		RemoveEmpty:   false,
		QuoteValues:   false,
	}
}

// Normalize applies normalization rules to a slice of EnvEntry.
func Normalize(entries []EnvEntry, opts NormalizeOptions) []EnvEntry {
	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if opts.UppercaseKeys {
			e.Key = strings.ToUpper(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.RemoveEmpty && e.Value == "" {
			continue
		}
		if opts.QuoteValues && e.Value != "" {
			if !strings.HasPrefix(e.Value, `"`) {
				e.Value = `"` + e.Value + `"`
			}
		}
		result = append(result, e)
	}
	return result
}

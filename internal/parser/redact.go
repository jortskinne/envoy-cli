package parser

import (
	"strings"
)

// RedactOptions controls how entries are redacted.
type RedactOptions struct {
	// Keys is an explicit list of keys to redact regardless of sensitivity.
	Keys []string
	// RedactSensitive auto-redacts keys detected as sensitive.
	RedactSensitive bool
	// Placeholder is the string used to replace redacted values.
	Placeholder string
}

// DefaultRedactOptions returns sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		RedactSensitive: true,
		Placeholder:     "[REDACTED]",
	}
}

// Redact replaces values of sensitive or specified keys with a placeholder.
func Redact(entries []EnvEntry, opts RedactOptions) []EnvEntry {
	if opts.Placeholder == "" {
		opts.Placeholder = "[REDACTED]"
	}

	explicit := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[strings.ToUpper(k)] = true
	}

	result := make([]EnvEntry, len(entries))
	for i, e := range entries {
		copy := e
		if explicit[strings.ToUpper(e.Key)] || (opts.RedactSensitive && IsSensitive(e.Key)) {
			copy.Value = opts.Placeholder
		}
		result[i] = copy
	}
	return result
}

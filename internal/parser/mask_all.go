package parser

import "strings"

// MaskAllOptions controls behaviour of MaskAll.
type MaskAllOptions struct {
	// Placeholder replaces every sensitive value.
	Placeholder string
	// RevealTrailing shows the last N characters of a value (0 = hide all).
	RevealTrailing int
	// Keys is an explicit list of keys to mask; if empty, IsSensitive is used.
	Keys []string
}

// DefaultMaskAllOptions returns sensible defaults.
func DefaultMaskAllOptions() MaskAllOptions {
	return MaskAllOptions{
		Placeholder:    "****",
		RevealTrailing: 0,
	}
}

// MaskAll returns a copy of entries with all sensitive (or explicitly listed)
// values replaced according to opts.
func MaskAll(entries []EnvEntry, opts MaskAllOptions) []EnvEntry {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	out := make([]EnvEntry, len(entries))
	for i, e := range entries {
		copy := e
		shouldMask := keySet[strings.ToUpper(e.Key)] || (len(keySet) == 0 && IsSensitive(e.Key))
		if shouldMask {
			copy.Value = MaskValue(e.Value, MaskOptions{
				Placeholder:    opts.Placeholder,
				RevealTrailing: opts.RevealTrailing,
			})
		}
		out[i] = copy
	}
	return out
}

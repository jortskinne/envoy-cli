package parser

import "strings"

// DefaultPrefixOptions returns sensible defaults for prefix operations.
func DefaultPrefixOptions() PrefixOptions {
	return PrefixOptions{
		Overwrite: false,
		DryRun:    false,
	}
}

// PrefixOptions controls how prefix add/remove behaves.
type PrefixOptions struct {
	Overwrite bool
	DryRun    bool
}

// AddPrefix prepends prefix to every entry key that does not already have it.
func AddPrefix(entries []EnvEntry, prefix string, opts PrefixOptions) []EnvEntry {
	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if !strings.HasPrefix(e.Key, prefix) {
			newKey := prefix + e.Key
			if !opts.DryRun {
				e.Key = newKey
			}
		} else if opts.Overwrite {
			// already has prefix — leave as-is; overwrite flag is a no-op here
		}
		result = append(result, e)
	}
	return result
}

// RemovePrefix strips prefix from every entry key that starts with it.
func RemovePrefix(entries []EnvEntry, prefix string, opts PrefixOptions) []EnvEntry {
	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			stripped := strings.TrimPrefix(e.Key, prefix)
			if stripped == "" {
				// skip keys that become empty after stripping
				continue
			}
			if !opts.DryRun {
				e.Key = stripped
			}
		}
		result = append(result, e)
	}
	return result
}

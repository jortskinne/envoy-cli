package parser

import (
	"strings"
)

// TransformFn is a function applied to an entry's value.
type TransformFn func(string) string

// TransformOptions controls how entries are transformed.
type TransformOptions struct {
	// Keys to transform; empty means all keys.
	Keys []string
	// Uppercase converts values to uppercase.
	Uppercase bool
	// Lowercase converts values to lowercase.
	Lowercase bool
	// TrimSpace trims leading/trailing whitespace from values.
	TrimSpace bool
	// Custom is an optional user-supplied transform function.
	Custom TransformFn
}

// DefaultTransformOptions returns safe defaults.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{TrimSpace: true}
}

// Transform applies value transformations to matching entries.
func Transform(entries []EnvEntry, opts TransformOptions) []EnvEntry {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	result := make([]EnvEntry, len(entries))
	for i, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			result[i] = e
			continue
		}
		v := e.Value
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.Uppercase {
			v = strings.ToUpper(v)
		} else if opts.Lowercase {
			v = strings.ToLower(v)
		}
		if opts.Custom != nil {
			v = opts.Custom(v)
		}
		result[i] = EnvEntry{Key: e.Key, Value: v, Comment: e.Comment}
	}
	return result
}

package parser

import "fmt"

// SetKeyOptions controls behaviour of SetKey.
type SetKeyOptions struct {
	// Overwrite allows replacing an existing key's value.
	Overwrite bool
	// DryRun returns the resulting entries without side-effects.
	DryRun bool
}

// DefaultSetKeyOptions returns sensible defaults.
func DefaultSetKeyOptions() SetKeyOptions {
	return SetKeyOptions{
		Overwrite: false,
		DryRun:    false,
	}
}

// SetKey adds or updates a single key=value pair in entries.
// If the key already exists and Overwrite is false, an error is returned.
func SetKey(entries []EnvEntry, key, value string, opts SetKeyOptions) ([]EnvEntry, error) {
	if key == "" {
		return nil, fmt.Errorf("key must not be empty")
	}

	result := make([]EnvEntry, len(entries))
	copy(result, entries)

	for i, e := range result {
		if e.Key == key {
			if !opts.Overwrite {
				return nil, fmt.Errorf("key %q already exists; use overwrite flag to replace it", key)
			}
			result[i].Value = value
			return result, nil
		}
	}

	// Key not found — append it.
	result = append(result, EnvEntry{Key: key, Value: value})
	return result, nil
}

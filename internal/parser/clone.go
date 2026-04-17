package parser

import "fmt"

// CloneOptions controls how entries are cloned between environments.
type CloneOptions struct {
	// Prefix filters entries to only those starting with this prefix.
	Prefix string
	// StripPrefix removes the prefix from cloned keys.
	StripPrefix bool
	// Overwrite allows overwriting existing keys in the destination.
	Overwrite bool
}

// DefaultCloneOptions returns sensible defaults.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		Prefix:      "",
		StripPrefix: false,
		Overwrite:   false,
	}
}

// Clone copies entries from src into dst according to opts.
// Returns the merged slice and a count of how many entries were cloned.
func Clone(dst, src []Entry, opts CloneOptions) ([]Entry, int, error) {
	destIndex := make(map[string]int, len(dst))
	for i, e := range dst {
		destIndex[e.Key] = i
	}

	result := make([]Entry, len(dst))
	copy(result, dst)

	cloned := 0
	for _, e := range src {
		if opts.Prefix != "" && len(e.Key) < len(opts.Prefix) {
			continue
		}
		if opts.Prefix != "" && e.Key[:len(opts.Prefix)] != opts.Prefix {
			continue
		}

		newKey := e.Key
		if opts.StripPrefix && opts.Prefix != "" {
			newKey = e.Key[len(opts.Prefix):]
			if newKey == "" {
				return nil, 0, fmt.Errorf("stripping prefix from key %q results in empty key", e.Key)
			}
		}

		cloned++
		if idx, exists := destIndex[newKey]; exists {
			if opts.Overwrite {
				result[idx] = Entry{Key: newKey, Value: e.Value}
			}
		} else {
			result = append(result, Entry{Key: newKey, Value: e.Value})
			destIndex[newKey] = len(result) - 1
		}
	}

	return result, cloned, nil
}

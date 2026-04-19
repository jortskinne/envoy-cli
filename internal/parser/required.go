package parser

import "fmt"

// RequiredOptions configures required key checking.
type RequiredOptions struct {
	AllowEmpty bool // if true, keys present but empty are considered satisfied
}

// DefaultRequiredOptions returns sensible defaults.
func DefaultRequiredOptions() RequiredOptions {
	return RequiredOptions{
		AllowEmpty: false,
	}
}

// RequiredResult holds the outcome of a required-key check.
type RequiredResult struct {
	Key     string
	Present bool
	Empty   bool
}

// CheckRequired verifies that all keys in required are present (and non-empty
// unless AllowEmpty is set) within entries. It returns one result per required
// key and a non-nil error if any check fails.
func CheckRequired(entries []EnvEntry, required []string, opts RequiredOptions) ([]RequiredResult, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	var results []RequiredResult
	var missing []string

	for _, key := range required {
		val, ok := lookup[key]
		res := RequiredResult{Key: key, Present: ok, Empty: ok && val == ""}
		results = append(results, res)

		if !ok {
			missing = append(missing, key)
		} else if !opts.AllowEmpty && val == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return results, fmt.Errorf("missing or empty required keys: %v", missing)
	}
	return results, nil
}

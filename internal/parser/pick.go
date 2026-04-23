package parser

import "fmt"

// DefaultPickOptions returns a PickOptions with safe defaults.
func DefaultPickOptions() PickOptions {
	return PickOptions{
		StrictMode: false,
	}
}

// PickOptions controls the behaviour of Pick.
type PickOptions struct {
	// Keys is the ordered list of keys to extract.
	Keys []string
	// StrictMode causes Pick to return an error if any requested key is absent.
	StrictMode bool
}

// Pick returns a new slice containing only the entries whose keys appear in
// opts.Keys, preserving the order specified in opts.Keys.
// When StrictMode is true, an error is returned for every key that is not
// found in entries.
func Pick(entries []EnvEntry, opts PickOptions) ([]EnvEntry, error) {
	lookup := make(map[string]EnvEntry, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e
	}

	var result []EnvEntry
	var missing []string

	for _, key := range opts.Keys {
		if key == "" {
			continue
		}
		e, ok := lookup[key]
		if !ok {
			if opts.StrictMode {
				missing = append(missing, key)
			}
			continue
		}
		result = append(result, e)
	}

	if len(missing) > 0 {
		return result, fmt.Errorf("pick: missing keys: %v", missing)
	}
	return result, nil
}

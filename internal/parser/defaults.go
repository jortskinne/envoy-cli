package parser

// DefaultsOptions controls how default values are applied.
type DefaultsOptions struct {
	// Overwrite replaces existing non-empty values with defaults.
	Overwrite bool
	// SkipEmpty skips applying a default when the default value itself is empty.
	SkipEmpty bool
}

// DefaultDefaultsOptions returns sensible defaults.
func DefaultDefaultsOptions() DefaultsOptions {
	return DefaultsOptions{
		Overwrite: false,
		SkipEmpty: true,
	}
}

// ApplyDefaults merges default values into entries.
// Keys present in defaults but missing from entries are added.
// Keys already present in entries are only overwritten if opts.Overwrite is true.
func ApplyDefaults(entries []EnvEntry, defaults map[string]string, opts DefaultsOptions) []EnvEntry {
	existing := make(map[string]int, len(entries))
	for i, e := range entries {
		existing[e.Key] = i
	}

	result := make([]EnvEntry, len(entries))
	copy(result, entries)

	for key, val := range defaults {
		if opts.SkipEmpty && val == "" {
			continue
		}
		if idx, found := existing[key]; found {
			if opts.Overwrite {
				result[idx].Value = val
			}
		} else {
			result = append(result, EnvEntry{Key: key, Value: val})
			existing[key] = len(result) - 1
		}
	}
	return result
}

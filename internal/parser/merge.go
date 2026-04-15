package parser

// MergeOptions controls how two sets of env entries are merged.
type MergeOptions struct {
	// Overwrite replaces existing keys in base with values from overlay.
	Overwrite bool
	// SkipEmpty ignores overlay entries whose value is empty.
	SkipEmpty bool
}

// DefaultMergeOptions returns sensible defaults for merging.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Overwrite: false,
		SkipEmpty: true,
	}
}

// Merge combines base and overlay EnvEntry slices according to opts.
// The returned slice preserves the order of base entries, with new
// keys from overlay appended at the end.
func Merge(base, overlay []EnvEntry, opts MergeOptions) []EnvEntry {
	// Build an index of base keys for fast lookup.
	index := make(map[string]int, len(base))
	for i, e := range base {
		index[e.Key] = i
	}

	result := make([]EnvEntry, len(base))
	copy(result, base)

	for _, oe := range overlay {
		if oe.Key == "" {
			continue
		}
		if opts.SkipEmpty && oe.Value == "" {
			continue
		}

		if idx, exists := index[oe.Key]; exists {
			if opts.Overwrite {
				result[idx] = oe
			}
		} else {
			result = append(result, oe)
			index[oe.Key] = len(result) - 1
		}
	}

	return result
}

package parser

// SliceOptions controls how Slice trims an entry list to a sub-range.
type SliceOptions struct {
	// Start is the zero-based index of the first entry to include (inclusive).
	Start int
	// End is the zero-based index of the last entry to include (exclusive).
	// A value of -1 means "until the end of the slice".
	End int
	// Keys restricts slicing to entries whose keys are in this list, applied
	// before positional slicing.
	Keys []string
}

// DefaultSliceOptions returns a SliceOptions that returns all entries.
func DefaultSliceOptions() SliceOptions {
	return SliceOptions{
		Start: 0,
		End:   -1,
	}
}

// Slice returns a sub-range of entries. If Keys is non-empty only those entries
// are considered before the positional Start/End window is applied.
func Slice(entries []EnvEntry, opts SliceOptions) ([]EnvEntry, error) {
	pool := entries

	if len(opts.Keys) > 0 {
		keySet := make(map[string]struct{}, len(opts.Keys))
		for _, k := range opts.Keys {
			keySet[k] = struct{}{}
		}
		filtered := make([]EnvEntry, 0, len(opts.Keys))
		for _, e := range entries {
			if _, ok := keySet[e.Key]; ok {
				filtered = append(filtered, e)
			}
		}
		pool = filtered
	}

	n := len(pool)
	start := opts.Start
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}

	end := opts.End
	if end < 0 || end > n {
		end = n
	}
	if end < start {
		end = start
	}

	result := make([]EnvEntry, end-start)
	copy(result, pool[start:end])
	return result, nil
}

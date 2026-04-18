package parser

// ReorderOptions controls how entries are reordered to match a reference key list.
type ReorderOptions struct {
	// AppendExtra appends keys not in the reference list at the end.
	// If false, extra keys are dropped.
	AppendExtra bool
}

// DefaultReorderOptions returns sensible defaults.
func DefaultReorderOptions() ReorderOptions {
	return ReorderOptions{
		AppendExtra: true,
	}
}

// Reorder reorders entries to match the order defined by keys.
// Keys not present in entries are skipped.
// Entries whose keys are not in keys are appended at the end if AppendExtra is true.
func Reorder(entries []EnvEntry, keys []string, opts ReorderOptions) []EnvEntry {
	index := make(map[string]EnvEntry, len(entries))
	for _, e := range entries {
		index[e.Key] = e
	}

	seen := make(map[string]bool)
	result := make([]EnvEntry, 0, len(entries))

	for _, k := range keys {
		if e, ok := index[k]; ok {
			result = append(result, e)
			seen[k] = true
		}
	}

	if opts.AppendExtra {
		for _, e := range entries {
			if !seen[e.Key] {
				result = append(result, e)
			}
		}
	}

	return result
}

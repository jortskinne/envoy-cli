package parser

import "strings"

// ShrinkOptions controls how shrinking is performed.
type ShrinkOptions struct {
	// RemoveComments strips comment-only entries.
	RemoveComments bool
	// RemoveEmpty removes entries with empty values.
	RemoveEmpty bool
	// DedupeKeys removes duplicate keys, keeping the last occurrence.
	DedupeKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
}

// DefaultShrinkOptions returns sensible defaults.
func DefaultShrinkOptions() ShrinkOptions {
	return ShrinkOptions{
		RemoveComments: true,
		RemoveEmpty:    false,
		DedupeKeys:     true,
		TrimValues:     true,
	}
}

// Shrink reduces an env entry slice by applying the configured cleanup passes.
// It returns a new slice and a count of how many entries were removed.
func Shrink(entries []EnvEntry, opts ShrinkOptions) ([]EnvEntry, int) {
	original := len(entries)
	result := make([]EnvEntry, 0, len(entries))

	// Pass 1: remove comment-only entries.
	for _, e := range entries {
		if opts.RemoveComments && strings.HasPrefix(strings.TrimSpace(e.Key), "#") {
			continue
		}
		result = append(result, e)
	}

	// Pass 2: trim values.
	if opts.TrimValues {
		for i := range result {
			result[i].Value = strings.TrimSpace(result[i].Value)
		}
	}

	// Pass 3: remove empty values.
	if opts.RemoveEmpty {
		filtered := result[:0]
		for _, e := range result {
			if e.Value != "" {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	// Pass 4: deduplicate keys (keep last).
	if opts.DedupeKeys {
		seen := make(map[string]int, len(result))
		for i, e := range result {
			seen[e.Key] = i
		}
		deduped := make([]EnvEntry, 0, len(seen))
		added := make(map[string]bool, len(seen))
		// Traverse in reverse to preserve last-wins semantics, then reverse back.
		for i := len(result) - 1; i >= 0; i-- {
			e := result[i]
			if seen[e.Key] == i && !added[e.Key] {
				deduped = append(deduped, e)
				added[e.Key] = true
			}
		}
		// Reverse to restore original order.
		for l, r := 0, len(deduped)-1; l < r; l, r = l+1, r-1 {
			deduped[l], deduped[r] = deduped[r], deduped[l]
		}
		result = deduped
	}

	return result, original - len(result)
}

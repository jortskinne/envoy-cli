package parser

import "strings"

// MergeValuesStrategy controls how conflicting values are resolved during a merge.
type MergeValuesStrategy string

const (
	// MergeStrategyBase keeps the value from the base entries on conflict.
	MergeStrategyBase MergeValuesStrategy = "base"
	// MergeStrategyOther prefers the value from the other (incoming) entries on conflict.
	MergeStrategyOther MergeValuesStrategy = "other"
	// MergeStrategyConcat concatenates base and other values with a separator.
	MergeStrategyConcat MergeValuesStrategy = "concat"
)

// MergeValuesOptions configures the MergeValues operation.
type MergeValuesOptions struct {
	// Strategy determines how conflicting keys are resolved.
	Strategy MergeValuesStrategy
	// ConcatSeparator is used when Strategy == MergeStrategyConcat.
	ConcatSeparator string
	// Keys restricts merging to specific keys; if empty, all keys are considered.
	Keys []string
	// IgnoreEmpty skips entries from the other set whose value is empty.
	IgnoreEmpty bool
}

// DefaultMergeValuesOptions returns sensible defaults for MergeValues.
func DefaultMergeValuesOptions() MergeValuesOptions {
	return MergeValuesOptions{
		Strategy:        MergeStrategyOther,
		ConcatSeparator: ",",
	}
}

// MergeValues combines two slices of EnvEntry, resolving conflicts according to
// the provided options. Keys that exist only in one set are always included.
// The order of the returned slice is: base entries first (in original order),
// followed by any keys that appear only in other (in their original order).
func MergeValues(base, other []EnvEntry, opts MergeValuesOptions) []EnvEntry {
	allowedKeys := buildMergeKeySet(opts.Keys)

	// Index other entries for quick lookup.
	otherIndex := make(map[string]EnvEntry, len(other))
	for _, e := range other {
		otherIndex[e.Key] = e
	}

	// Track which other keys have been consumed.
	consumed := make(map[string]bool, len(other))

	result := make([]EnvEntry, 0, len(base)+len(other))

	for _, baseEntry := range base {
		otherEntry, exists := otherIndex[baseEntry.Key]

		// If a key filter is active and this key is not in it, keep base value unchanged.
		if len(allowedKeys) > 0 && !allowedKeys[baseEntry.Key] {
			result = append(result, baseEntry)
			if exists {
				consumed[baseEntry.Key] = true
			}
			continue
		}

		if !exists {
			// Key only in base — keep as-is.
			result = append(result, baseEntry)
			continue
		}

		consumed[baseEntry.Key] = true

		if opts.IgnoreEmpty && strings.TrimSpace(otherEntry.Value) == "" {
			result = append(result, baseEntry)
			continue
		}

		merged := resolveConflict(baseEntry, otherEntry, opts)
		result = append(result, merged)
	}

	// Append keys that exist only in other.
	for _, otherEntry := range other {
		if consumed[otherEntry.Key] {
			continue
		}
		if len(allowedKeys) > 0 && !allowedKeys[otherEntry.Key] {
			continue
		}
		if opts.IgnoreEmpty && strings.TrimSpace(otherEntry.Value) == "" {
			continue
		}
		result = append(result, otherEntry)
	}

	return result
}

// resolveConflict applies the chosen strategy to a pair of conflicting entries.
func resolveConflict(base, other EnvEntry, opts MergeValuesOptions) EnvEntry {
	switch opts.Strategy {
	case MergeStrategyBase:
		return base
	case MergeStrategyConcat:
		sep := opts.ConcatSeparator
		if sep == "" {
			sep = ","
		}
		merged := base
		if base.Value != "" && other.Value != "" {
			merged.Value = base.Value + sep + other.Value
		} else if other.Value != "" {
			merged.Value = other.Value
		}
		return merged
	default: // MergeStrategyOther
		result := base
		result.Value = other.Value
		return result
	}
}

// buildMergeKeySet converts a slice of key names into a lookup map.
func buildMergeKeySet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

package parser

// OmitOptions controls which entries are removed from the result.
type OmitOptions struct {
	// Keys is an explicit list of keys to omit.
	Keys []string

	// Prefix removes all entries whose key starts with the given prefix.
	Prefix string

	// SensitiveOnly removes all entries detected as sensitive by the masking rules.
	SensitiveOnly bool

	// EmptyValues removes entries whose value is the empty string.
	EmptyValues bool
}

// DefaultOmitOptions returns an OmitOptions with no filters applied.
func DefaultOmitOptions() OmitOptions {
	return OmitOptions{}
}

// Omit removes entries from src according to the supplied options and returns
// the filtered slice. The original slice is never modified.
//
// Filters are applied with OR semantics: an entry is omitted if it matches
// ANY of the active criteria.
func Omit(src []EnvEntry, opts OmitOptions) []EnvEntry {
	keySet := buildOmitKeySet(opts.Keys)

	result := make([]EnvEntry, 0, len(src))
	for _, entry := range src {
		if shouldOmit(entry, opts, keySet) {
			continue
		}
		result = append(result, entry)
	}
	return result
}

// buildOmitKeySet converts a slice of key names into a fast-lookup map.
func buildOmitKeySet(keys []string) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	return set
}

// shouldOmit returns true when the entry matches at least one omit criterion.
func shouldOmit(entry EnvEntry, opts OmitOptions, keySet map[string]struct{}) bool {
	// Explicit key list.
	if _, ok := keySet[entry.Key]; ok {
		return true
	}

	// Prefix filter.
	if opts.Prefix != "" && len(entry.Key) >= len(opts.Prefix) &&
		entry.Key[:len(opts.Prefix)] == opts.Prefix {
		return true
	}

	// Sensitive detection.
	if opts.SensitiveOnly && IsSensitive(entry.Key) {
		return true
	}

	// Empty value filter.
	if opts.EmptyValues && entry.Value == "" {
		return true
	}

	return false
}

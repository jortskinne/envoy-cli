package parser

import "strings"

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	Keys      []string // if set, only include these keys
	Prefix    string   // if set, only include keys with this prefix
	Exclude   []string // keys to exclude
	Sensitive bool     // if true, only include sensitive keys
}

// DefaultFilterOptions returns permissive defaults.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{}
}

// Filter returns a subset of entries based on the provided options.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	allowKey := buildKeySet(opts.Keys)
	excludeKey := buildKeySet(opts.Exclude)

	var result []Entry
	for _, e := range entries {
		if len(allowKey) > 0 {
			if _, ok := allowKey[e.Key]; !ok {
				continue
			}
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if _, ok := excludeKey[e.Key]; ok {
			continue
		}
		if opts.Sensitive && !IsSensitive(e.Key, DefaultMaskOptions()) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func buildKeySet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

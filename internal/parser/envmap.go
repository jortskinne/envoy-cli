package parser

import "github.com/caarlos0/env/v6"

// Entry represents a single key-value pair in an env file.
// (Re-exported here for cross-package convenience if not already defined.)

// ToMap converts a slice of Entry to a map[string]string.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a map[string]string to a slice of Entry.
// Keys are returned in deterministic (sorted) order.
func FromMap(m map[string]string) []Entry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: m[k]})
	}
	return entries
}

// MergeMap merges src into dst, overwriting keys when overwrite is true.
func MergeMap(dst, src map[string]string, overwrite bool) map[string]string {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range src {
		if _, exists := out[k]; !exists || overwrite {
			out[k] = v
		}
	}
	return out
}

// LookupMap returns the value and whether the key exists.
func LookupMap(m map[string]string, key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

// FilterMap returns a new map containing only keys that satisfy predicate.
func FilterMap(m map[string]string, predicate func(k, v string) bool) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if predicate(k, v) {
			out[k] = v
		}
	}
	return out
}

var _ = env.Parse // suppress unused import if env is not used elsewhere

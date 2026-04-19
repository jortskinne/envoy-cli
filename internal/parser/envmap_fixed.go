package parser

import "sort"

// ToMap converts a slice of Entry to a map[string]string.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a map[string]string to a slice of Entry in sorted key order.
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

// MergeMap merges src into dst. When overwrite is true existing keys are replaced.
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

// FilterMap returns entries from m where predicate returns true.
func FilterMap(m map[string]string, predicate func(k, v string) bool) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if predicate(k, v) {
			out[k] = v
		}
	}
	return out
}

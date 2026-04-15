package differ

import (
	"fmt"
	"sort"

	"github.com/envoy-cli/internal/parser"
)

// DiffType represents the type of difference found between two env files.
type DiffType string

const (
	DiffAdded   DiffType = "added"
	DiffRemoved DiffType = "removed"
	DiffChanged DiffType = "changed"
)

// DiffEntry represents a single difference between two env files.
type DiffEntry struct {
	Key      string
	Type     DiffType
	OldValue string
	NewValue string
}

// String returns a human-readable representation of a DiffEntry.
func (d DiffEntry) String() string {
	switch d.Type {
	case DiffAdded:
		return fmt.Sprintf("+ %s=%s", d.Key, d.NewValue)
	case DiffRemoved:
		return fmt.Sprintf("- %s=%s", d.Key, d.OldValue)
	case DiffChanged:
		return fmt.Sprintf("~ %s: %s -> %s", d.Key, d.OldValue, d.NewValue)
	}
	return ""
}

// DiffResult holds all differences between two env files.
type DiffResult struct {
	Entries []DiffEntry
}

// HasDiff returns true if there are any differences.
func (r *DiffResult) HasDiff() bool {
	return len(r.Entries) > 0
}

// Diff compares two sets of parsed env entries and returns a DiffResult.
// maskOpts controls whether sensitive values are masked in the output.
func Diff(base, target []parser.Entry, maskOpts parser.MaskOptions) *DiffResult {
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	targetMap := make(map[string]string, len(target))
	for _, e := range target {
		targetMap[e.Key] = e.Value
	}

	var entries []DiffEntry

	for key, baseVal := range baseMap {
		if targetVal, ok := targetMap[key]; !ok {
			entries = append(entries, DiffEntry{
				Key:      key,
				Type:     DiffRemoved,
				OldValue: maskedValue(key, baseVal, maskOpts),
			})
		} else if baseVal != targetVal {
			entries = append(entries, DiffEntry{
				Key:      key,
				Type:     DiffChanged,
				OldValue: maskedValue(key, baseVal, maskOpts),
				NewValue: maskedValue(key, targetVal, maskOpts),
			})
		}
	}

	for key, targetVal := range targetMap {
		if _, ok := baseMap[key]; !ok {
			entries = append(entries, DiffEntry{
				Key:      key,
				Type:     DiffAdded,
				NewValue: maskedValue(key, targetVal, maskOpts),
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &DiffResult{Entries: entries}
}

func maskedValue(key, value string, opts parser.MaskOptions) string {
	if parser.IsSensitive(key, opts) {
		return parser.MaskValue(value, opts)
	}
	return value
}

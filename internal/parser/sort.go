package parser

import (
	"sort"

	"github.com/envoy-cli/internal/parser"
)

// SortOrder defines how entries should be sorted.
type SortOrder string

const (
	SortAlpha      SortOrder = "alpha"
	SortAlphaDesc  SortOrder = "alpha-desc"
	SortByGroup    SortOrder = "group"
)

// SortOptions configures sorting behaviour.
type SortOptions struct {
	Order      SortOrder
	GroupByPrefix bool // group keys sharing the same prefix (e.g. DB_, AWS_)
}

// DefaultSortOptions returns sensible defaults.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:         SortAlpha,
		GroupByPrefix: false,
	}
}

// Sort returns a new slice of EnvEntry values ordered according to opts.
func Sort(entries []EnvEntry, opts SortOptions) []EnvEntry {
	out := make([]EnvEntry, len(entries))
	copy(out, entries)

	switch opts.Order {
	case SortAlphaDesc:
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key > out[j].Key
		})
	case SortByGroup:
		sort.SliceStable(out, func(i, j int) bool {
			gi := groupPrefix(out[i].Key)
			gj := groupPrefix(out[j].Key)
			if gi != gj {
				return gi < gj
			}
			return out[i].Key < out[j].Key
		})
	default: // SortAlpha
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}

// groupPrefix returns the first underscore-delimited segment of a key,
// or the whole key if no underscore is present.
func groupPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' {
			return key[:i]
		}
	}
	return key
}

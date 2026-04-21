package parser

import "strings"

// EnvDiffOptions controls behaviour of the EnvDiff function.
type EnvDiffOptions struct {
	// IgnoreCase treats keys as case-insensitive when comparing.
	IgnoreCase bool
	// TrimValues strips leading/trailing whitespace from values before comparing.
	TrimValues bool
}

// DefaultEnvDiffOptions returns sensible defaults.
func DefaultEnvDiffOptions() EnvDiffOptions {
	return EnvDiffOptions{
		IgnoreCase: false,
		TrimValues: true,
	}
}

// EnvDiffResult holds the outcome of comparing two env entry slices.
type EnvDiffResult struct {
	Added    []EnvEntry // present in other, missing from base
	Removed  []EnvEntry // present in base, missing from other
	Changed  []EnvChange // key exists in both but value differs
	Unchanged []EnvEntry // identical in both
}

// EnvChange represents a key whose value differs between base and other.
type EnvChange struct {
	Key      string
	BaseVal  string
	OtherVal string
}

// EnvDiff compares two slices of EnvEntry and returns a structured diff result.
func EnvDiff(base, other []EnvEntry, opts EnvDiffOptions) EnvDiffResult {
	normKey := func(k string) string {
		if opts.IgnoreCase {
			return strings.ToLower(k)
		}
		return k
	}
	normVal := func(v string) string {
		if opts.TrimValues {
			return strings.TrimSpace(v)
		}
		return v
	}

	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[normKey(e.Key)] = normVal(e.Value)
	}

	otherMap := make(map[string]string, len(other))
	for _, e := range other {
		otherMap[normKey(e.Key)] = normVal(e.Value)
	}

	var result EnvDiffResult

	for _, e := range other {
		nk := normKey(e.Key)
		nv := normVal(e.Value)
		if baseVal, ok := baseMap[nk]; ok {
			if baseVal == nv {
				result.Unchanged = append(result.Unchanged, e)
			} else {
				result.Changed = append(result.Changed, EnvChange{
					Key:      e.Key,
					BaseVal:  baseVal,
					OtherVal: nv,
				})
			}
		} else {
			result.Added = append(result.Added, e)
		}
	}

	for _, e := range base {
		if _, ok := otherMap[normKey(e.Key)]; !ok {
			result.Removed = append(result.Removed, e)
		}
	}

	return result
}

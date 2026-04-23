package parser

import "strings"

// CrossCheckOptions configures the CrossCheck operation.
type CrossCheckOptions struct {
	// IgnoreCase treats keys as case-insensitive when matching.
	IgnoreCase bool
	// RequireAllBase ensures every key in base exists in other.
	RequireAllBase bool
	// RequireAllOther ensures every key in other exists in base.
	RequireAllOther bool
}

// DefaultCrossCheckOptions returns sensible defaults.
func DefaultCrossCheckOptions() CrossCheckOptions {
	return CrossCheckOptions{
		IgnoreCase:      false,
		RequireAllBase:  true,
		RequireAllOther: false,
	}
}

// CrossCheckResult holds the result of a cross-check between two env sets.
type CrossCheckResult struct {
	Key       string
	Status    string // "missing_in_other", "missing_in_base", "type_mismatch", "ok"
	BaseValue string
	OtherValue string
}

// CrossCheck compares two sets of entries and reports key-level discrepancies.
func CrossCheck(base, other []EnvEntry, opts CrossCheckOptions) []CrossCheckResult {
	normalize := func(k string) string {
		if opts.IgnoreCase {
			return strings.ToUpper(k)
		}
		return k
	}

	otherMap := make(map[string]string)
	for _, e := range other {
		otherMap[normalize(e.Key)] = e.Value
	}

	baseMap := make(map[string]string)
	for _, e := range base {
		baseMap[normalize(e.Key)] = e.Value
	}

	var results []CrossCheckResult

	if opts.RequireAllBase {
		for _, e := range base {
			nk := normalize(e.Key)
			if oval, found := otherMap[nk]; !found {
				results = append(results, CrossCheckResult{
					Key:       e.Key,
					Status:    "missing_in_other",
					BaseValue: e.Value,
				})
			} else {
				results = append(results, CrossCheckResult{
					Key:        e.Key,
					Status:     "ok",
					BaseValue:  e.Value,
					OtherValue: oval,
				})
			}
		}
	}

	if opts.RequireAllOther {
		for _, e := range other {
			nk := normalize(e.Key)
			if _, found := baseMap[nk]; !found {
				results = append(results, CrossCheckResult{
					Key:        e.Key,
					Status:     "missing_in_base",
					OtherValue: e.Value,
				})
			}
		}
	}

	return results
}

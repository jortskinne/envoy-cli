package parser

import "fmt"

// CompareOptions configures the Compare behaviour.
type CompareOptions struct {
	// IgnoreCase treats key lookup as case-insensitive.
	IgnoreCase bool
	// IgnoreWhitespace trims values before comparing.
	IgnoreWhitespace bool
}

// DefaultCompareOptions returns sensible defaults.
func DefaultCompareOptions() CompareOptions {
	return CompareOptions{
		IgnoreCase:       false,
		IgnoreWhitespace: true,
	}
}

// CompareResult holds the outcome for a single key.
type CompareResult struct {
	Key      string
	BaseVal  string
	OtherVal string
	Status   string // "match", "mismatch", "base_only", "other_only"
}

// Compare performs a key-by-key comparison between base and other entry slices.
// It returns one CompareResult per unique key found in either slice.
func Compare(base, other []EnvEntry, opts CompareOptions) []CompareResult {
	normalizeKey := func(k string) string {
		if opts.IgnoreCase {
			return fmt.Sprintf("%s", []byte(k)) // keep as-is; handled via map key
		}
		return k
	}
	normalizeVal := func(v string) string {
		if opts.IgnoreWhitespace {
			return trimWhitespace(v)
		}
		return v
	}

	baseMap := make(map[string]string)
	for _, e := range base {
		baseMap[normalizeKey(e.Key)] = normalizeVal(e.Value)
	}

	otherMap := make(map[string]string)
	for _, e := range other {
		otherMap[normalizeKey(e.Key)] = normalizeVal(e.Value)
	}

	seen := make(map[string]bool)
	var results []CompareResult

	for k, bv := range baseMap {
		seen[k] = true
		if ov, ok := otherMap[k]; ok {
			status := "match"
			if bv != ov {
				status = "mismatch"
			}
			results = append(results, CompareResult{Key: k, BaseVal: bv, OtherVal: ov, Status: status})
		} else {
			results = append(results, CompareResult{Key: k, BaseVal: bv, OtherVal: "", Status: "base_only"})
		}
	}

	for k, ov := range otherMap {
		if !seen[k] {
			results = append(results, CompareResult{Key: k, BaseVal: "", OtherVal: ov, Status: "other_only"})
		}
	}

	return results
}

// trimWhitespace removes leading/trailing spaces and tabs.
func trimWhitespace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

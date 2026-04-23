package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CastOptions controls how values are cast.
type CastOptions struct {
	// Keys is the explicit list of keys to cast. If empty, all entries are attempted.
	Keys []string
	// StrictMode causes Cast to return an error if a cast fails.
	StrictMode bool
}

// DefaultCastOptions returns sensible defaults.
func DefaultCastOptions() CastOptions {
	return CastOptions{
		StrictMode: false,
	}
}

// CastResult holds the inferred type and normalised string value for one entry.
type CastResult struct {
	Key          string
	Original     string
	Normalized   string
	InferredType string // "bool", "int", "float", "string"
}

// Cast attempts to infer and normalize the type of each entry value.
// Booleans are normalized to "true"/"false", integers and floats are
// re-serialised via strconv to strip leading zeros, etc.
func Cast(entries []Entry, opts CastOptions) ([]CastResult, error) {
	keySet := buildKeySet(opts.Keys)

	var results []CastResult
	for _, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			continue
		}

		r, err := castValue(e.Key, e.Value)
		if err != nil && opts.StrictMode {
			return nil, fmt.Errorf("cast: key %q: %w", e.Key, err)
		}
		results = append(results, r)
	}
	return results, nil
}

func castValue(key, raw string) (CastResult, error) {
	v := strings.TrimSpace(raw)

	// bool
	lower := strings.ToLower(v)
	if lower == "true" || lower == "yes" || lower == "1" {
		return CastResult{Key: key, Original: raw, Normalized: "true", InferredType: "bool"}, nil
	}
	if lower == "false" || lower == "no" || lower == "0" {
		return CastResult{Key: key, Original: raw, Normalized: "false", InferredType: "bool"}, nil
	}

	// int
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return CastResult{Key: key, Original: raw, Normalized: strconv.FormatInt(i, 10), InferredType: "int"}, nil
	}

	// float
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return CastResult{Key: key, Original: raw, Normalized: strconv.FormatFloat(f, 'f', -1, 64), InferredType: "float"}, nil
	}

	// fallback string
	return CastResult{Key: key, Original: raw, Normalized: v, InferredType: "string"}, nil
}

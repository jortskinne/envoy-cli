package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultCoerceOptions returns sensible defaults for Coerce.
func DefaultCoerceOptions() CoerceOptions {
	return CoerceOptions{
		BoolTrue:  []string{"1", "yes", "true", "on"},
		BoolFalse: []string{"0", "no", "false", "off"},
		Keys:      nil, // nil means all keys
	}
}

// CoerceOptions controls how values are coerced.
type CoerceOptions struct {
	// BoolTrue is the canonical string written for truthy booleans.
	BoolTrue []string
	// BoolFalse is the canonical string written for falsy booleans.
	BoolFalse []string
	// Keys restricts coercion to specific keys; nil means all.
	Keys []string
	// TargetType forces every matched key to a specific type ("bool", "int", "float", "string").
	// When empty the type is inferred.
	TargetType string
}

// CoerceResult records a single coercion outcome.
type CoerceResult struct {
	Key      string
	OldValue string
	NewValue string
	Type     string
}

// Coerce normalises env entry values to canonical typed representations.
// It returns the updated entries and a log of every change made.
func Coerce(entries []EnvEntry, opts CoerceOptions) ([]EnvEntry, []CoerceResult, error) {
	keySet := buildCoerceKeySet(opts.Keys)

	out := make([]EnvEntry, 0, len(entries))
	var results []CoerceResult

	for _, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			out = append(out, e)
			continue
		}

		newVal, typ, err := coerceValue(e.Value, opts)
		if err != nil {
			return nil, nil, fmt.Errorf("coerce %q: %w", e.Key, err)
		}

		if newVal != e.Value {
			results = append(results, CoerceResult{
				Key:      e.Key,
				OldValue: e.Value,
				NewValue: newVal,
				Type:     typ,
			})
			e.Value = newVal
		}
		out = append(out, e)
	}
	return out, results, nil
}

func coerceValue(v string, opts CoerceOptions) (string, string, error) {
	target := strings.ToLower(opts.TargetType)

	switch target {
	case "bool":
		return coerceBool(v, opts)
	case "int":
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return "", "", fmt.Errorf("cannot coerce %q to int", v)
		}
		return strconv.FormatInt(n, 10), "int", nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", "", fmt.Errorf("cannot coerce %q to float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), "float", nil
	case "string":
		return v, "string", nil
	default:
		// infer
		if res, typ, err := coerceBool(v, opts); err == nil && typ == "bool" {
			return res, typ, nil
		}
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return strconv.FormatInt(n, 10), "int", nil
		}
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return strconv.FormatFloat(f, 'f', -1, 64), "float", nil
		}
		return v, "string", nil
	}
}

func coerceBool(v string, opts CoerceOptions) (string, string, error) {
	lower := strings.ToLower(strings.TrimSpace(v))
	for _, t := range opts.BoolTrue {
		if lower == strings.ToLower(t) {
			return "true", "bool", nil
		}
	}
	for _, f := range opts.BoolFalse {
		if lower == strings.ToLower(f) {
			return "false", "bool", nil
		}
	}
	return v, "", fmt.Errorf("not a bool")
}

func buildCoerceKeySet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

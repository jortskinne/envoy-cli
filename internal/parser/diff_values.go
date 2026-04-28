package parser

import "strings"

// DiffValuesOptions controls the behaviour of DiffValues.
type DiffValuesOptions struct {
	// IgnoreCase compares values case-insensitively.
	IgnoreCase bool
	// TrimSpace strips surrounding whitespace before comparing.
	TrimSpace bool
	// MaskSensitive replaces sensitive values in the result with "***".
	MaskSensitive bool
}

// DefaultDiffValuesOptions returns sensible defaults.
func DefaultDiffValuesOptions() DiffValuesOptions {
	return DiffValuesOptions{
		IgnoreCase:    false,
		TrimSpace:     true,
		MaskSensitive: false,
	}
}

// ValueDiff represents a single key whose value differs between two env sets.
type ValueDiff struct {
	Key      string
	BaseVal  string
	OtherVal string
	// Sensitive is true when the key looks like a secret.
	Sensitive bool
}

// DiffValues compares values for keys that exist in both base and other,
// returning only the entries where the values differ.
func DiffValues(base, other []EnvEntry, opts DiffValuesOptions) []ValueDiff {
	otherMap := make(map[string]string, len(other))
	for _, e := range other {
		otherMap[e.Key] = e.Value
	}

	var diffs []ValueDiff
	for _, e := range base {
		oval, exists := otherMap[e.Key]
		if !exists {
			continue
		}

		bv, ov := e.Value, oval
		if opts.TrimSpace {
			bv = strings.TrimSpace(bv)
			ov = strings.TrimSpace(ov)
		}

		cmpBase, cmpOther := bv, ov
		if opts.IgnoreCase {
			cmpBase = strings.ToLower(bv)
			cmpOther = strings.ToLower(ov)
		}

		if cmpBase == cmpOther {
			continue
		}

		sensitive := IsSensitive(e.Key, DefaultMaskOptions())
		displayBase, displayOther := bv, ov
		if opts.MaskSensitive && sensitive {
			displayBase = "***"
			displayOther = "***"
		}

		diffs = append(diffs, ValueDiff{
			Key:       e.Key,
			BaseVal:   displayBase,
			OtherVal:  displayOther,
			Sensitive: sensitive,
		})
	}
	return diffs
}

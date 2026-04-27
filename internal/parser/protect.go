package parser

import "fmt"

// DefaultProtectOptions returns a ProtectOptions with safe defaults.
func DefaultProtectOptions() ProtectOptions {
	return ProtectOptions{
		Keys:       []string{},
		AllowEmpty: false,
		DryRun:     false,
	}
}

// ProtectOptions configures how keys are protected (made read-only / locked).
type ProtectOptions struct {
	// Keys is the explicit list of keys to protect.
	Keys []string
	// AllowEmpty allows protecting keys that currently have empty values.
	AllowEmpty bool
	// DryRun returns the result without modifying anything.
	DryRun bool
}

// ProtectResult describes the outcome of a Protect operation.
type ProtectResult struct {
	Protected []string
	Skipped   []string
}

// Protect marks the given keys as locked by prepending a "# PROTECTED" comment
// immediately before the key's line. Keys that are already protected or that
// have empty values (when AllowEmpty is false) are skipped.
func Protect(entries []EnvEntry, opts ProtectOptions) ([]EnvEntry, ProtectResult, error) {
	if len(opts.Keys) == 0 {
		return nil, ProtectResult{}, fmt.Errorf("protect: no keys specified")
	}

	target := buildKeySet(opts.Keys)
	result := ProtectResult{}

	out := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if !target[e.Key] {
			out = append(out, e)
			continue
		}
		if !opts.AllowEmpty && e.Value == "" {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}
		// Already protected?
		if e.Comment == "PROTECTED" {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}
		e.Comment = "PROTECTED"
		result.Protected = append(result.Protected, e.Key)
		out = append(out, e)
	}
	return out, result, nil
}

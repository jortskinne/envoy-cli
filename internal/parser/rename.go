package parser

import "fmt"

// RenameOptions controls rename behaviour.
type RenameOptions struct {
	// DryRun reports what would change without modifying entries.
	DryRun bool
	// ErrorIfMissing returns an error when the old key does not exist.
	ErrorIfMissing bool
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		DryRun:         false,
		ErrorIfMissing: true,
	}
}

// RenameResult describes the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Skipped bool
	Reason  string
}

// Rename renames oldKey to newKey in entries, returning updated entries and a result.
func Rename(entries []Entry, oldKey, newKey string, opts RenameOptions) ([]Entry, RenameResult, error) {
	result := RenameResult{OldKey: oldKey, NewKey: newKey}

	if oldKey == newKey {
		result.Skipped = true
		result.Reason = "old and new key are identical"
		return entries, result, nil
	}

	foundOld := false
	for _, e := range entries {
		if e.Key == oldKey {
			foundOld = true
		}
		if e.Key == newKey {
			result.Skipped = true
			result.Reason = fmt.Sprintf("key %q already exists", newKey)
			return entries, result, nil
		}
	}

	if !foundOld {
		if opts.ErrorIfMissing {
			return entries, result, fmt.Errorf("key %q not found", oldKey)
		}
		result.Skipped = true
		result.Reason = fmt.Sprintf("key %q not found", oldKey)
		return entries, result, nil
	}

	if opts.DryRun {
		result.Renamed = true
		result.Reason = "dry run"
		return entries, result, nil
	}

	out := make([]Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Key == oldKey {
			out[i].Key = newKey
			break
		}
	}
	result.Renamed = true
	return out, result, nil
}

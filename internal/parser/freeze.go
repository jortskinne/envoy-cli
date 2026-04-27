package parser

import "fmt"

// FreezeOptions controls which entries are frozen (made read-only via comment tag).
type FreezeOptions struct {
	// Keys is an explicit list of keys to freeze.
	Keys []string
	// FreezeAll freezes every entry when true.
	FreezeAll bool
	// Tag is the inline comment marker written to signal a frozen entry.
	Tag string
	// DryRun returns modified entries without writing.
	DryRun bool
}

// DefaultFreezeOptions returns sensible defaults.
func DefaultFreezeOptions() FreezeOptions {
	return FreezeOptions{
		Tag: "@frozen",
	}
}

// Freeze marks the specified entries with a frozen tag in their Comment field.
// Frozen entries are treated as read-only by commands that respect the tag.
func Freeze(entries []EnvEntry, opts FreezeOptions) ([]EnvEntry, error) {
	if len(opts.Keys) == 0 && !opts.FreezeAll {
		return nil, fmt.Errorf("freeze: no keys specified and FreezeAll is false")
	}

	tag := opts.Tag
	if tag == "" {
		tag = "@frozen"
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	result := make([]EnvEntry, len(entries))
	copy(result, entries)

	for i, e := range result {
		if opts.FreezeAll || keySet[e.Key] {
			result[i] = markFrozen(e, tag)
		}
	}
	return result, nil
}

// IsFrozen reports whether an entry carries the frozen tag.
func IsFrozen(e EnvEntry, tag string) bool {
	if tag == "" {
		tag = "@frozen"
	}
	return containsTag(e.Comment, tag)
}

func markFrozen(e EnvEntry, tag string) EnvEntry {
	if containsTag(e.Comment, tag) {
		return e
	}
	if e.Comment == "" {
		e.Comment = tag
	} else {
		e.Comment = e.Comment + " " + tag
	}
	return e
}

func containsTag(comment, tag string) bool {
	for i := 0; i <= len(comment)-len(tag); i++ {
		if comment[i:i+len(tag)] == tag {
			return true
		}
	}
	return false
}

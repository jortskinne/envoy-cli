package parser

import "fmt"

// PromoteOptions controls how keys are promoted between environments.
type PromoteOptions struct {
	// Overwrite allows existing keys in target to be overwritten.
	Overwrite bool
	// Keys limits promotion to specific keys. If empty, all keys are promoted.
	Keys []string
	// DryRun returns the result without writing.
	DryRun bool
}

// DefaultPromoteOptions returns sensible defaults.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		Overwrite: false,
		Keys:      nil,
		DryRun:    false,
	}
}

// PromoteResult holds the outcome of a promotion.
type PromoteResult struct {
	Merged  []EnvEntry
	Added   []string
	Skipped []string
}

// Promote copies entries from source into target according to opts.
func Promote(source, target []EnvEntry, opts PromoteOptions) (PromoteResult, error) {
	keyFilter := buildKeySet(opts.Keys)

	targetMap := make(map[string]int, len(target))
	for i, e := range target {
		targetMap[e.Key] = i
	}

	result := PromoteResult{Merged: append([]EnvEntry{}, target...)}

	for _, src := range source {
		if len(keyFilter) > 0 && !keyFilter[src.Key] {
			continue
		}
		if idx, exists := targetMap[src.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, src.Key)
				continue
			}
			result.Merged[idx].Value = src.Value
			result.Added = append(result.Added, src.Key)
		} else {
			result.Merged = append(result.Merged, src)
			targetMap[src.Key] = len(result.Merged) - 1
			result.Added = append(result.Added, src.Key)
		}
	}

	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			if _, found := buildKeySet(result.Added)[k]; !found {
				if _, skipped := buildKeySet(result.Skipped)[k]; !skipped {
					return result, fmt.Errorf("key %q not found in source", k)
				}
			}
		}
	}

	return result, nil
}

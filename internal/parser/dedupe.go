package parser

import "github.com/envoy-cli/envoy-cli/internal/parser"

// DedupeOptions controls deduplication behavior.
type DedupeOptions struct {
	// KeepFirst retains the first occurrence; if false, keeps the last.
	KeepFirst bool
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{KeepFirst: true}
}

// DedupeResult holds the output of a deduplication pass.
type DedupeResult struct {
	Entries  []EnvEntry
	Removed  []EnvEntry // duplicate entries that were dropped
}

// Dedupe removes duplicate keys from entries according to opts.
func Dedupe(entries []EnvEntry, opts DedupeOptions) DedupeResult {
	seen := make(map[string]int) // key -> index in out
	out := make([]EnvEntry, 0, len(entries))
	var removed []EnvEntry

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			if opts.KeepFirst {
				// current entry is the duplicate — drop it
				removed = append(removed, e)
			} else {
				// replace the earlier entry, mark old as removed
				removed = append(removed, out[idx])
				out[idx] = e
			}
		} else {
			seen[e.Key] = len(out)
			out = append(out, e)
		}
	}

	return DedupeResult{Entries: out, Removed: removed}
}

package parser

import "fmt"

// PinOptions controls how pinning behaves.
type PinOptions struct {
	Keys      []string // explicit keys to pin; empty means all
	Overwrite bool     // overwrite existing pin tag
	DryRun    bool     // do not modify; just report
}

// DefaultPinOptions returns sensible defaults.
func DefaultPinOptions() PinOptions {
	return PinOptions{
		Keys:      nil,
		Overwrite: false,
		DryRun:    false,
	}
}

// PinResult describes a single pinned entry.
type PinResult struct {
	Key     string
	Value   string
	Pinned  bool
	Skipped bool // already pinned and Overwrite=false
}

const pinTag = "pinned"

// Pin marks entries with a "# @pinned" comment so downstream tools
// can refuse to overwrite them. Returns the updated entries and a
// per-key result log.
func Pin(entries []EnvEntry, opts PinOptions) ([]EnvEntry, []PinResult, error) {
	keySet := buildKeySet(opts.Keys)

	results := make([]PinResult, 0, len(entries))
	out := make([]EnvEntry, 0, len(entries))

	for _, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			out = append(out, e)
			continue
		}
		if e.Key == "" {
			out = append(out, e)
			continue
		}

		alreadyPinned := containsTag(e.Comment, pinTag)
		if alreadyPinned && !opts.Overwrite {
			results = append(results, PinResult{Key: e.Key, Value: e.Value, Skipped: true})
			out = append(out, e)
			continue
		}

		if !opts.DryRun {
			if !alreadyPinned {
				if e.Comment == "" {
					e.Comment = fmt.Sprintf("# @%s", pinTag)
				} else {
					e.Comment = fmt.Sprintf("%s @%s", e.Comment, pinTag)
				}
			}
		}
		results = append(results, PinResult{Key: e.Key, Value: e.Value, Pinned: true})
		out = append(out, e)
	}

	return out, results, nil
}

// IsPinned reports whether the entry carries the pinned tag.
func IsPinned(e EnvEntry) bool {
	return containsTag(e.Comment, pinTag)
}

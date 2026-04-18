package parser

import (
	"fmt"
	"strings"
)

// KeyConflict describes a conflict found between two sets of entries.
type KeyConflict struct {
	Key     string
	BaseVal string
	OtherVal string
	Kind    string // "type_mismatch", "case_conflict", "shadow"
}

// DefaultValidateKeysOptions returns sensible defaults.
func DefaultValidateKeysOptions() ValidateKeysOptions {
	return ValidateKeysOptions{
		DetectCaseConflicts: true,
		DetectShadows:       true,
	}
}

// ValidateKeysOptions controls which checks are performed.
type ValidateKeysOptions struct {
	DetectCaseConflicts bool
	DetectShadows       bool
}

// ValidateKeys checks for case conflicts and shadowed keys within a single
// slice of EnvEntry values and returns a list of conflicts found.
func ValidateKeys(entries []EnvEntry, opts ValidateKeysOptions) []KeyConflict {
	var conflicts []KeyConflict

	// Build canonical (lower-case) → first-seen key map
	canonicMap := make(map[string]string) // lower -> original
	seen := make(map[string]int)          // original key -> count

	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		lower := strings.ToLower(e.Key)

		if opts.DetectCaseConflicts {
			if existing, ok := canonicMap[lower]; ok && existing != e.Key {
				conflicts = append(conflicts, KeyConflict{
					Key:  e.Key,
					Kind: "case_conflict",
					BaseVal: fmt.Sprintf("conflicts with %q", existing),
				})
			} else {
				canonicMap[lower] = e.Key
			}
		}

		if opts.DetectShadows {
			seen[e.Key]++
			if seen[e.Key] == 2 {
				conflicts = append(conflicts, KeyConflict{
					Key:  e.Key,
					Kind: "shadow",
					BaseVal: fmt.Sprintf("key %q defined more than once", e.Key),
				})
			}
		}
	}

	return conflicts
}

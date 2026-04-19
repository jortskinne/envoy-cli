package parser

import (
	"fmt"
	"os"
	"strings"
)

// ResolveOptions controls how env entries are resolved against the OS environment.
type ResolveOptions struct {
	// OverrideWithOS replaces entry values with matching OS env vars if present.
	OverrideWithOS bool
	// FillMissing populates entries with empty values from OS env vars.
	FillMissing bool
	// FailOnMissing returns an error if a key has no value after resolution.
	FailOnMissing bool
}

func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		OverrideWithOS: false,
		FillMissing:    true,
		FailOnMissing:  false,
	}
}

// ResolveResult holds the resolved entry and metadata.
type ResolveResult struct {
	Entry   Entry
	Source  string // "file", "os", or "empty"
	Missing bool
}

// Resolve merges entries with OS environment variables based on options.
func Resolve(entries []Entry, opts ResolveOptions) ([]ResolveResult, error) {
	var results []ResolveResult
	var missing []string

	for _, e := range entries {
		result := ResolveResult{Entry: e, Source: "file"}

		osVal, osFound := os.LookupEnv(e.Key)

		if osFound && opts.OverrideWithOS {
			result.Entry.Value = osVal
			result.Source = "os"
		} else if e.Value == "" && osFound && opts.FillMissing {
			result.Entry.Value = osVal
			result.Source = "os"
		} else if e.Value == "" && !osFound {
			result.Missing = true
			result.Source = "empty"
			missing = append(missing, e.Key)
		}

		results = append(results, result)
	}

	if opts.FailOnMissing && len(missing) > 0 {
		return results, fmt.Errorf("missing values for keys: %s", strings.Join(missing, ", "))
	}

	return results, nil
}

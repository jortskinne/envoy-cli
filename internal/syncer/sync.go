package syncer

import (
	"fmt"
	"os"

	"github.com/envoy-cli/internal/parser"
)

// SyncMode controls how missing or extra keys are handled during sync.
type SyncMode string

const (
	// ModeAddMissing only adds keys present in source but missing in target.
	ModeAddMissing SyncMode = "add-missing"
	// ModeOverwrite adds missing keys and overwrites changed values.
	ModeOverwrite SyncMode = "overwrite"
	// ModeFull adds missing, overwrites changed, and removes extra keys.
	ModeFull SyncMode = "full"
)

// SyncOptions configures the behaviour of a sync operation.
type SyncOptions struct {
	Mode     SyncMode
	DryRun   bool
	MaskKeys bool
}

// DefaultSyncOptions returns a conservative default configuration.
func DefaultSyncOptions() SyncOptions {
	return SyncOptions{
		Mode:     ModeAddMissing,
		DryRun:   false,
		MaskKeys: true,
	}
}

// SyncResult summarises the changes applied (or that would be applied) during a sync.
type SyncResult struct {
	Added     []string
	Updated   []string
	Removed   []string
	Unchanged []string
}

// HasChanges reports whether any additions, updates, or removals occurred.
func (r SyncResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Updated) > 0 || len(r.Removed) > 0
}

// Sync merges source entries into target entries according to opts.
// It returns the merged entry slice and a SyncResult describing what changed.
func Sync(source, target []parser.Entry, opts SyncOptions) ([]parser.Entry, SyncResult, error) {
	result := SyncResult{}

	sourceMap := make(map[string]parser.Entry, len(source))
	for _, e := range source {
		sourceMap[e.Key] = e
	}

	targetMap := make(map[string]parser.Entry, len(target))
	for _, e := range target {
		targetMap[e.Key] = e
	}

	merged := make([]parser.Entry, 0, len(target))

	// Iterate existing target entries.
	for _, te := range target {
		se, inSource := sourceMap[te.Key]
		switch {
		case !inSource && opts.Mode == ModeFull:
			// Remove keys not present in source.
			result.Removed = append(result.Removed, te.Key)
		case inSource && se.Value != te.Value && (opts.Mode == ModeOverwrite || opts.Mode == ModeFull):
			// Overwrite changed value.
			result.Updated = append(result.Updated, te.Key)
			merged = append(merged, se)
		default:
			result.Unchanged = append(result.Unchanged, te.Key)
			merged = append(merged, te)
		}
	}

	// Add keys present in source but missing in target.
	for _, se := range source {
		if _, exists := targetMap[se.Key]; !exists {
			result.Added = append(result.Added, se.Key)
			merged = append(merged, se)
		}
	}

	return merged, result, nil
}

// WriteEnvFile serialises entries to path in KEY=VALUE format.
func WriteEnvFile(path string, entries []parser.Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("syncer: create file %q: %w", path, err)
	}
	defer f.Close()

	for _, e := range entries {
		line := fmt.Sprintf("%s=%s\n", e.Key, e.Value)
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("syncer: write entry %q: %w", e.Key, err)
		}
	}
	return nil
}

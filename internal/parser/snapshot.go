package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time    `json:"timestamp"`
	Label     string       `json:"label,omitempty"`
	Entries   []EnvEntry   `json:"entries"`
}

// SaveSnapshot writes the given entries as a JSON snapshot to the specified path.
func SaveSnapshot(path, label string, entries []EnvEntry) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: failed to marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: failed to write file: %w", err)
	}

	return nil
}

// LoadSnapshot reads a JSON snapshot from the given path.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to read file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: failed to parse JSON: %w", err)
	}

	return &snap, nil
}

// SnapshotEntries is a convenience helper that returns the entries from a
// snapshot file, suitable for passing directly into Diff or Validate.
func SnapshotEntries(path string) ([]EnvEntry, error) {
	snap, err := LoadSnapshot(path)
	if err != nil {
		return nil, err
	}
	return snap.Entries, nil
}

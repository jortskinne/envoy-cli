package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// SchemaEntry defines the expected shape of a single env key.
type SchemaEntry struct {
	Required    bool   `json:"required"`
	DefaultValue string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
}

// Schema maps key names to their schema definitions.
type Schema map[string]SchemaEntry

// LoadSchema reads a JSON schema file and returns a parsed Schema.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parsing schema JSON: %w", err)
	}

	return schema, nil
}

// RequiredKeys returns all keys in the schema marked as required.
func (s Schema) RequiredKeys() []string {
	keys := make([]string, 0, len(s))
	for k, v := range s {
		if v.Required {
			keys = append(keys, k)
		}
	}
	return keys
}

// ApplyDefaults fills in default values from the schema for any missing keys
// in the provided entries slice, returning a new slice with defaults appended.
func (s Schema) ApplyDefaults(entries []EnvEntry) []EnvEntry {
	existing := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		existing[e.Key] = struct{}{}
	}

	result := make([]EnvEntry, len(entries))
	copy(result, entries)

	for key, def := range s {
		if _, found := existing[key]; !found && def.DefaultValue != "" {
			result = append(result, EnvEntry{Key: key, Value: def.DefaultValue})
		}
	}

	return result
}

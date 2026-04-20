package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultImportOptions returns sensible defaults for Import.
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Overwrite:    false,
		SkipInvalid:  false,
		Format:       "", // auto-detect
	}
}

// ImportOptions controls how external files are imported.
type ImportOptions struct {
	Overwrite   bool
	SkipInvalid bool
	Format      string // "dotenv", "json", or "" for auto
}

// Import reads entries from src (dotenv or JSON) and merges them into base.
// New keys are always added; existing keys are only overwritten if Overwrite is true.
func Import(base []EnvEntry, srcPath string, opts ImportOptions) ([]EnvEntry, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("import: cannot read file %q: %w", srcPath, err)
	}

	format := opts.Format
	if format == "" {
		switch strings.ToLower(filepath.Ext(srcPath)) {
		case ".json":
			format = "json"
		default:
			format = "dotenv"
		}
	}

	var incoming []EnvEntry
	switch format {
	case "json":
		incoming, err = importJSON(data)
	case "dotenv":
		incoming, err = importDotEnv(data, opts.SkipInvalid)
	default:
		return nil, fmt.Errorf("import: unsupported format %q", format)
	}
	if err != nil {
		return nil, err
	}

	existing := make(map[string]int, len(base))
	for i, e := range base {
		existing[e.Key] = i
	}

	result := make([]EnvEntry, len(base))
	copy(result, base)

	for _, entry := range incoming {
		if idx, found := existing[entry.Key]; found {
			if opts.Overwrite {
				result[idx].Value = entry.Value
			}
		} else {
			result = append(result, entry)
			existing[entry.Key] = len(result) - 1
		}
	}
	return result, nil
}

func importJSON(data []byte) ([]EnvEntry, error) {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("import: invalid JSON: %w", err)
	}
	var entries []EnvEntry
	for k, v := range m {
		entries = append(entries, EnvEntry{Key: k, Value: v})
	}
	return entries, nil
}

func importDotEnv(data []byte, skipInvalid bool) ([]EnvEntry, error) {
	lines := strings.Split(string(data), "\n")
	var entries []EnvEntry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		e, err := parseLine(line)
		if err != nil {
			if skipInvalid {
				continue
			}
			return nil, fmt.Errorf("import: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

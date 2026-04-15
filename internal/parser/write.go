package parser

import (
	"fmt"
	"os"
	"strings"
)

// WriteEnvFile serialises a slice of EnvEntry back to a .env file.
// Comments and blank lines are not preserved from the original parse;
// each entry is written as KEY=VALUE with optional inline comment.
func WriteEnvFile(path string, entries []EnvEntry) error {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		// Re-quote values that contain spaces or special characters.
		if needsQuoting(val) {
			val = `"` + strings.ReplaceAll(val, `"`, `\"`) + `"`
		}
		line := fmt.Sprintf("%s=%s", e.Key, val)
		if e.Comment != "" {
			line += " # " + e.Comment
		}
		sb.WriteString(line + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// needsQuoting returns true if the value should be wrapped in double quotes
// when written to a .env file.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch == ' ' || ch == '\t' || ch == '#' || ch == '$' || ch == '\'' {
			return true
		}
	}
	return false
}

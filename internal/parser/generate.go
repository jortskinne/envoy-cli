package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateOptions controls how placeholder .env entries are generated.
type GenerateOptions struct {
	Length      int
	Sensitive   []string // keys that get random secret values
	Placeholder string   // default value for non-sensitive keys
	Prefix      string
}

func DefaultGenerateOptions() GenerateOptions {
	return GenerateOptions{
		Length:      32,
		Sensitive:   []string{},
		Placeholder: "CHANGEME",
		Prefix:      "",
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// Generate produces a slice of EnvEntry from a list of key names.
// Sensitive keys receive a random secret; others get the placeholder value.
func Generate(keys []string, opts GenerateOptions) ([]EnvEntry, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("generate: no keys provided")
	}

	sensitiveSet := make(map[string]bool, len(opts.Sensitive))
	for _, k := range opts.Sensitive {
		sensitiveSet[strings.ToUpper(k)] = true
	}

	entries := make([]EnvEntry, 0, len(keys))
	for _, k := range keys {
		key := opts.Prefix + strings.ToUpper(strings.TrimSpace(k))
		if key == "" {
			continue
		}
		var val string
		if sensitiveSet[key] || sensitiveSet[strings.ToUpper(k)] {
			val = randomString(opts.Length)
		} else {
			val = opts.Placeholder
		}
		entries = append(entries, EnvEntry{Key: key, Value: val})
	}
	return entries, nil
}

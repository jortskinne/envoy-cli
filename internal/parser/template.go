package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateOptions controls template rendering behaviour.
type TemplateOptions struct {
	// StrictMode returns an error for any unresolved placeholder.
	StrictMode bool
}

// DefaultTemplateOptions returns sensible defaults.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{StrictMode: true}
}

var templatePlaceholder = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// RenderTemplate replaces {{KEY}} placeholders in src with values from
// entries. If StrictMode is enabled any unresolved placeholder is an error.
func RenderTemplate(src string, entries []EnvEntry, opts TemplateOptions) (string, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	var renderErr error
	result := templatePlaceholder.ReplaceAllStringFunc(src, func(match string) string {
		sub := templatePlaceholder.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := sub[1]
		if val, ok := lookup[key]; ok {
			return val
		}
		// Fall back to OS environment.
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		if opts.StrictMode && renderErr == nil {
			renderErr = fmt.Errorf("template: unresolved placeholder %q", key)
		}
		return match
	})

	if renderErr != nil {
		return "", renderErr
	}
	return result, nil
}

// RenderTemplateFile reads a template file, resolves placeholders, and
// returns the rendered string.
func RenderTemplateFile(path string, entries []EnvEntry, opts TemplateOptions) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template: cannot read file %q: %w", path, err)
	}
	return RenderTemplate(strings.TrimRight(string(raw), "\n"), entries, opts)
}

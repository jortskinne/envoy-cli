package parser

import (
	"strings"
)

// TagEntry represents an env entry with an associated tag label.
type TagEntry struct {
	Key   string
	Value string
	Tag   string
}

// DefaultTagOptions returns sensible defaults for Tag.
func DefaultTagOptions() TagOptions {
	return TagOptions{
		Overwrite: false,
		Separator: "#tag:",
	}
}

// TagOptions controls how Tag behaves.
type TagOptions struct {
	// Tags maps key -> tag label to apply.
	Tags map[string]string
	// Overwrite replaces an existing inline tag comment.
	Overwrite bool
	// Separator is the inline comment prefix used to store tags.
	Separator string
}

// Tag applies string labels to matching env entries via inline comments.
// Tags are stored as inline comments using the configured separator.
func Tag(entries []EnvEntry, opts TagOptions) ([]EnvEntry, error) {
	if opts.Separator == "" {
		opts.Separator = DefaultTagOptions().Separator
	}

	result := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		label, ok := opts.Tags[e.Key]
		if !ok {
			result = append(result, e)
			continue
		}
		if e.Comment != "" && !opts.Overwrite {
			result = append(result, e)
			continue
		}
		e.Comment = opts.Separator + label
		result = append(result, e)
	}
	return result, nil
}

// ExtractTags reads inline tag comments from entries and returns a map of
// key -> tag label for every entry that carries a tag comment.
func ExtractTags(entries []EnvEntry, separator string) map[string]string {
	if separator == "" {
		separator = DefaultTagOptions().Separator
	}
	out := make(map[string]string)
	for _, e := range entries {
		if strings.HasPrefix(e.Comment, separator) {
			out[e.Key] = strings.TrimPrefix(e.Comment, separator)
		}
	}
	return out
}

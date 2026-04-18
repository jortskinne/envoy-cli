package parser

import (
	"strings"
)

// Annotation holds a comment attached to an env entry.
type Annotation struct {
	Key     string
	Comment string
}

// AnnotateOptions controls annotation behaviour.
type AnnotateOptions struct {
	// Overwrite existing inline comments.
	Overwrite bool
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{Overwrite: false}
}

// Annotate applies comments from the annotations map to matching entries.
// The key in annotations is the env key; the value is the comment text.
func Annotate(entries []Entry, annotations map[string]string, opts AnnotateOptions) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		if comment, ok := annotations[e.Key]; ok {
			if copy.Comment == "" || opts.Overwrite {
				copy.Comment = strings.TrimSpace(comment)
			}
		}
		result[i] = copy
	}
	return result
}

// ExtractAnnotations returns a map of key -> comment for all entries that
// have a non-empty inline comment.
func ExtractAnnotations(entries []Entry) map[string]string {
	out := make(map[string]string)
	for _, e := range entries {
		if e.Comment != "" {
			out[e.Key] = e.Comment
		}
	}
	return out
}

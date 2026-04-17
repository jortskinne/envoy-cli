package parser

import (
	"strings"
)

// GroupedOutput holds sections of env entries with optional section headers.
type GroupedOutput struct {
	Sections []Section
}

type Section struct {
	Header  string
	Entries []EnvEntry
}

// DefaultGroupOptions returns default options for GroupByPrefix.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Separator: "_",
		CommentHeader: true,
	}
}

type GroupOptions struct {
	Separator     string
	CommentHeader bool
}

// GroupByPrefix groups env entries by their key prefix (before first separator).
// Entries without a separator are placed in an "OTHER" section.
func GroupByPrefix(entries []EnvEntry, opts GroupOptions) GroupedOutput {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	order := []string{}
	sections := map[string][]EnvEntry{}

	for _, e := range entries {
		prefix := "OTHER"
		if idx := strings.Index(e.Key, opts.Separator); idx > 0 {
			prefix = e.Key[:idx]
		}
		if _, exists := sections[prefix]; !exists {
			order = append(order, prefix)
		}
		sections[prefix] = append(sections[prefix], e)
	}

	out := GroupedOutput{}
	for _, prefix := range order {
		out.Sections = append(out.Sections, Section{
			Header:  prefix,
			Entries: sections[prefix],
		})
	}
	return out
}

// FlattenGrouped converts a GroupedOutput back to a flat slice,
// optionally inserting comment header entries.
func FlattenGrouped(g GroupedOutput, commentHeader bool) []EnvEntry {
	var result []EnvEntry
	for _, sec := range g.Sections {
		if commentHeader {
			result = append(result, EnvEntry{
				Key:       "",
				Value:     "",
				Comment:   "# --- " + sec.Header + " ---",
				IsComment: true,
			})
		}
		result = append(result, sec.Entries...)
	}
	return result
}

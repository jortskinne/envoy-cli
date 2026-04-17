package parser

import "strings"

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	Separator string // separator used to split keys (default: "__")
	Prefix    string // only flatten keys with this prefix
	Lowercase bool   // convert resulting keys to lowercase
}

// DefaultFlattenOptions returns sensible defaults.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "__",
	}
}

// Flatten takes entries whose keys contain a separator and groups them into
// a dot-notation style flat map representation, returning new EnvEntry slice.
// Keys without the separator are passed through unchanged.
func Flatten(entries []EnvEntry, opts FlattenOptions) []EnvEntry {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	out := make([]EnvEntry, 0, len(entries))
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			out = append(out, e)
			continue
		}

		parts := strings.SplitN(e.Key, opts.Separator, 2)
		if len(parts) < 2 {
			out = append(out, e)
			continue
		}

		newKey := parts[0] + "." + parts[1]
		if opts.Lowercase {
			newKey = strings.ToLower(newKey)
		}

		out = append(out, EnvEntry{
			Key:     newKey,
			Value:   e.Value,
			Comment: e.Comment,
		})
	}
	return out
}

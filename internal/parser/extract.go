package parser

import "strings"

// ExtractOptions controls extraction behaviour.
type ExtractOptions struct {
	Keys      []string // explicit keys to extract
	Prefix    string   // extract all keys with this prefix
	StripPrefix bool   // remove the prefix from extracted keys
}

func DefaultExtractOptions() ExtractOptions {
	return ExtractOptions{}
}

// Extract returns a subset of entries matching the given keys or prefix.
func Extract(entries []Entry, opts ExtractOptions) ([]Entry, error) {
	keySet := buildExtractKeySet(opts.Keys)
	var result []Entry
	for _, e := range entries {
		matched := false
		if len(keySet) > 0 {
			if _, ok := keySet[e.Key]; ok {
				matched = true
			}
		}
		if opts.Prefix != "" && strings.HasPrefix(e.Key, opts.Prefix) {
			matched = true
		}
		if !matched {
			continue
		}
		out := e
		if opts.StripPrefix && opts.Prefix != "" {
			out.Key = strings.TrimPrefix(e.Key, opts.Prefix)
		}
		result = append(result, out)
	}
	return result, nil
}

func buildExtractKeySet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

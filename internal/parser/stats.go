package parser

import "strings"

// EnvEntry is imported from env_parser.go

type Stats struct {
	Total     int
	Empty     int
	Sensitive int
	Commented int
	Prefixes  map[string]int
}

type StatsOptions struct {
	SensitiveKeys []string
}

func DefaultStatsOptions() StatsOptions {
	return StatsOptions{
		SensitiveKeys: DefaultMaskOptions().SensitiveKeys,
	}
}

func ComputeStats(entries []EnvEntry, opts StatsOptions) Stats {
	s := Stats{
		Prefixes: make(map[string]int),
	}
	s.Total = len(entries)
	for _, e := range entries {
		if e.Value == "" {
			s.Empty++
		}
		if IsSensitive(e.Key, opts.SensitiveKeys) {
			s.Sensitive++
		}
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			prefix := e.Key[:idx]
			s.Prefixes[prefix]++
		}
	}
	return s
}

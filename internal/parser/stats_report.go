package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func WriteStatsReport(w io.Writer, s Stats, format string) error {
	switch format {
	case "json":
		return writeStatsJSON(w, s)
	default:
		return writeStatsText(w, s)
	}
}

func writeStatsText(w io.Writer, s Stats) error {
	fmt.Fprintf(w, "Total keys   : %d\n", s.Total)
	fmt.Fprintf(w, "Empty values : %d\n", s.Empty)
	fmt.Fprintf(w, "Sensitive    : %d\n", s.Sensitive)
	if len(s.Prefixes) > 0 {
		fmt.Fprintln(w, "Prefixes:")
		keys := make([]string, 0, len(s.Prefixes))
		for k := range s.Prefixes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  %-20s %d\n", k, s.Prefixes[k])
		}
	}
	return nil
}

func writeStatsJSON(w io.Writer, s Stats) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

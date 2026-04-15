package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var (
	compareIgnoreWhitespace bool
	compareIgnoreCase       bool
	compareOutput           string
)

func init() {
	compareCmd := &cobra.Command{
		Use:   "compare <base> <other>",
		Short: "Compare two .env files key-by-key",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	compareCmd.Flags().BoolVar(&compareIgnoreWhitespace, "ignore-whitespace", true, "Trim values before comparing")
	compareCmd.Flags().BoolVar(&compareIgnoreCase, "ignore-case", false, "Treat keys as case-insensitive")
	compareCmd.Flags().StringVarP(&compareOutput, "output", "o", "text", "Output format: text or json")
	RootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	baseEntries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}
	otherEntries, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading other file: %w", err)
	}

	opts := parser.CompareOptions{
		IgnoreWhitespace: compareIgnoreWhitespace,
		IgnoreCase:       compareIgnoreCase,
	}
	results := parser.Compare(baseEntries, otherEntries, opts)

	if compareOutput == "json" {
		return writeCompareJSON(results)
	}
	return writeCompareText(results)
}

func writeCompareText(results []parser.CompareResult) error {
	w := os.Stdout
	for _, r := range results {
		switch r.Status {
		case "match":
			fmt.Fprintf(w, "  = %s\n", r.Key)
		case "mismatch":
			fmt.Fprintf(w, "  ~ %s  (%q -> %q)\n", r.Key, r.BaseVal, r.OtherVal)
		case "base_only":
			fmt.Fprintf(w, "  - %s  (only in base)\n", r.Key)
		case "other_only":
			fmt.Fprintf(w, "  + %s  (only in other)\n", r.Key)
		}
	}
	return nil
}

func writeCompareJSON(results []parser.CompareResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var filterCmd = &cobra.Command{
	Use:   "filter <file>",
	Short: "Filter .env entries by key, prefix, or sensitivity",
	Args:  cobra.ExactArgs(1),
	RunE:  runFilter,
}

func init() {
	filterCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to include")
	filterCmd.Flags().String("prefix", "", "Only include keys with this prefix")
	filterCmd.Flags().StringSlice("exclude", nil, "Keys to exclude")
	filterCmd.Flags().Bool("sensitive", false, "Only include sensitive keys")
	filterCmd.Flags().String("out", "", "Output file (default: stdout)")
	RootCmd.AddCommand(filterCmd)
}

func runFilter(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultFilterOptions()

	if keys, _ := cmd.Flags().GetStringSlice("keys"); len(keys) > 0 {
		opts.Keys = keys
	}
	if prefix, _ := cmd.Flags().GetString("prefix"); prefix != "" {
		opts.Prefix = prefix
	}
	if excl, _ := cmd.Flags().GetStringSlice("exclude"); len(excl) > 0 {
		opts.Exclude = excl
	}
	opts.Sensitive, _ = cmd.Flags().GetBool("sensitive")

	filtered := parser.Filter(entries, opts)

	var sb strings.Builder
	for _, e := range filtered {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}

	outPath, _ := cmd.Flags().GetString("out")
	if outPath != "" {
		return os.WriteFile(outPath, []byte(sb.String()), 0644)
	}
	fmt.Print(sb.String())
	return nil
}

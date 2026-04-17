package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/parser"
)

var extractCmd = &cobra.Command{
	Use:   "extract <file>",
	Short: "Extract a subset of keys from an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runExtract,
}

func init() {
	extractCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to extract")
	extractCmd.Flags().String("prefix", "", "Extract all keys with this prefix")
	extractCmd.Flags().Bool("strip-prefix", false, "Strip the prefix from extracted keys")
	extractCmd.Flags().String("output", "", "Write result to file instead of stdout")
	rootCmd.AddCommand(extractCmd)
}

func runExtract(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	prefix, _ := cmd.Flags().GetString("prefix")
	strip, _ := cmd.Flags().GetBool("strip-prefix")
	output, _ := cmd.Flags().GetString("output")

	opts := parser.DefaultExtractOptions()
	opts.Keys = keys
	opts.Prefix = prefix
	opts.StripPrefix = strip

	result, err := parser.Extract(entries, opts)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for _, e := range result {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}

	if output != "" {
		return os.WriteFile(output, []byte(sb.String()), 0644)
	}
	fmt.Print(sb.String())
	return nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var placeholderCmd = &cobra.Command{
	Use:   "placeholder <file>",
	Short: "Detect placeholder values in an .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlaceholder,
}

func init() {
	placeholderCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	placeholderCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
	placeholderCmd.Flags().StringSliceP("pattern", "p", nil, "Additional regex patterns to detect placeholders")
	RootCmd.AddCommand(placeholderCmd)
}

func runPlaceholder(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	extraPatterns, _ := cmd.Flags().GetStringSlice("pattern")

	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultPlaceholderOptions()
	if len(extraPatterns) > 0 {
		opts.Patterns = append(opts.Patterns, extraPatterns...)
	}

	results, err := parser.FindPlaceholders(entries, opts)
	if err != nil {
		return fmt.Errorf("placeholder detection error: %w", err)
	}

	w := cmd.OutOrStdout()
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("could not create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return parser.WritePlaceholderReport(w, results, format)
}

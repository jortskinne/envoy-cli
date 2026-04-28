package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var diffValuesCmd = &cobra.Command{
	Use:   "diff-values <base> <other>",
	Short: "Show keys whose values differ between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiffValues,
}

func init() {
	diffValuesCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	diffValuesCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
	diffValuesCmd.Flags().Bool("ignore-case", false, "Compare values case-insensitively")
	diffValuesCmd.Flags().Bool("mask", false, "Mask sensitive values in output")
	rootCmd.AddCommand(diffValuesCmd)
}

func runDiffValues(cmd *cobra.Command, args []string) error {
	baseFile, otherFile := args[0], args[1]

	baseEntries, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}
	otherEntries, err := parser.ParseFile(otherFile)
	if err != nil {
		return fmt.Errorf("parsing other file: %w", err)
	}

	format, _ := cmd.Flags().GetString("format")
	ignoreCase, _ := cmd.Flags().GetBool("ignore-case")
	mask, _ := cmd.Flags().GetBool("mask")
	outPath, _ := cmd.Flags().GetString("output")

	opts := parser.DefaultDiffValuesOptions()
	opts.IgnoreCase = ignoreCase
	opts.MaskSensitive = mask

	diffs := parser.DiffValues(baseEntries, otherEntries, opts)

	w := cmd.OutOrStdout()
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return parser.WriteDiffValuesReport(diffs, format, w)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/differ"
	"envoy-cli/internal/parser"
)

var (
	diffFormat  string
	diffMasked  bool
	diffOutput  string
)

var diffCmd = &cobra.Command{
	Use:   "diff <base> <target>",
	Short: "Diff two .env files and report differences",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVarP(&diffFormat, "format", "f", "text", "Output format: text or json")
	diffCmd.Flags().BoolVarP(&diffMasked, "mask", "m", true, "Mask sensitive values in output")
	diffCmd.Flags().StringVarP(&diffOutput, "output", "o", "", "Write output to file instead of stdout")
}

func runDiff(cmd *cobra.Command, args []string) error {
	baseFile, targetFile := args[0], args[1]

	baseEntries, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("parsing base file %q: %w", baseFile, err)
	}

	targetEntries, err := parser.ParseFile(targetFile)
	if err != nil {
		return fmt.Errorf("parsing target file %q: %w", targetFile, err)
	}

	opts := differ.DiffOptions{
		MaskSecrets: diffMasked,
		MaskOptions: parser.DefaultMaskOptions(),
	}

	results := differ.Diff(baseEntries, targetEntries, opts)

	out := cmd.OutOrStdout()
	if diffOutput != "" {
		f, err := os.Create(diffOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	return differ.WriteReport(out, results, diffFormat)
}

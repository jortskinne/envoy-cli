package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var pinCmd = &cobra.Command{
	Use:   "pin <file>",
	Short: "Mark env entries as pinned to prevent accidental overwrites",
	Args:  cobra.ExactArgs(1),
	RunE:  runPin,
}

func init() {
	pinCmd.Flags().StringSliceP("keys", "k", nil, "Keys to pin (default: all)")
	pinCmd.Flags().BoolP("overwrite", "o", false, "Re-pin already-pinned entries")
	pinCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	pinCmd.Flags().StringP("format", "f", "text", "Output format: text|json")
	pinCmd.Flags().StringP("output", "O", "", "Write report to file instead of stdout")
	RootCmd.AddCommand(pinCmd)
}

func runPin(cmd *cobra.Command, args []string) error {
	file := args[0]
	keys, _ := cmd.Flags().GetStringSlice("keys")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	format, _ := cmd.Flags().GetString("format")
	outputFile, _ := cmd.Flags().GetString("output")

	entries, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultPinOptions()
	opts.Keys = keys
	opts.Overwrite = overwrite
	opts.DryRun = dryRun

	updated, results, err := parser.Pin(entries, opts)
	if err != nil {
		return err
	}

	if !dryRun {
		if err := parser.WriteEnvFile(file, updated); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	report := parser.BuildPinReport(results)

	w := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot open output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return parser.WritePinReport(w, report, format)
}

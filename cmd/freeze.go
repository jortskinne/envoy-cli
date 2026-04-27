package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var freezeCmd = &cobra.Command{
	Use:   "freeze <file>",
	Short: "Mark env entries as frozen (read-only)",
	Args:  cobra.ExactArgs(1),
	RunE:  runFreeze,
}

func init() {
	freezeCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to freeze")
	freezeCmd.Flags().Bool("all", false, "Freeze all entries")
	freezeCmd.Flags().Bool("dry-run", false, "Print result without writing")
	freezeCmd.Flags().String("output", "", "Write output to file instead of stdout")
	freezeCmd.Flags().String("format", "text", "Output format: text or json")
	RootCmd.AddCommand(freezeCmd)
}

func runFreeze(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("freeze: parse error: %w", err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	all, _ := cmd.Flags().GetBool("all")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFile, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")

	// Flatten comma-joined values from --keys flag
	var flatKeys []string
	for _, k := range keys {
		for _, part := range strings.Split(k, ",") {
			if t := strings.TrimSpace(part); t != "" {
				flatKeys = append(flatKeys, t)
			}
		}
	}

	opts := parser.DefaultFreezeOptions()
	opts.Keys = flatKeys
	opts.FreezeAll = all
	opts.DryRun = dryRun

	updated, err := parser.Freeze(entries, opts)
	if err != nil {
		return err
	}

	report := parser.BuildFreezeReport(entries, updated, opts.Tag)

	w := os.Stdout
	if outputFile != "" && !dryRun {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("freeze: cannot create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if !dryRun {
		if err := parser.WriteEnvFile(filePath, updated); err != nil {
			return fmt.Errorf("freeze: write error: %w", err)
		}
	}

	return parser.WriteFreezeReport(w, report, format)
}

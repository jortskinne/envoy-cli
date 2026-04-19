package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var maskAllCmd = &cobra.Command{
	Use:   "maskall <file>",
	Short: "Mask all sensitive values in an env file and print a report",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaskAll,
}

func init() {
	maskAllCmd.Flags().StringP("output", "o", "", "Write masked env to file")
	maskAllCmd.Flags().StringP("format", "f", "text", "Report format: text|json")
	maskAllCmd.Flags().StringSliceP("keys", "k", nil, "Explicit keys to mask (comma-separated)")
	maskAllCmd.Flags().Int("reveal", 0, "Reveal last N characters of masked values")
	rootCmd.AddCommand(maskAllCmd)
}

func runMaskAll(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	reveal, _ := cmd.Flags().GetInt("reveal")
	format, _ := cmd.Flags().GetString("format")
	outFile, _ := cmd.Flags().GetString("output")

	// Normalise keys (trim whitespace)
	for i, k := range keys {
		keys[i] = strings.TrimSpace(k)
	}

	opts := parser.DefaultMaskAllOptions()
	opts.Keys = keys
	opts.RevealTrailing = reveal

	masked := parser.MaskAll(entries, opts)
	report := parser.BuildMaskReport(entries, masked)

	// Write masked env file if requested
	if outFile != "" {
		if err := parser.WriteEnvFile(outFile, masked); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	} else {
		for _, e := range masked {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
		}
	}

	// Always print report to stderr
	if err := parser.WriteMaskReport(os.Stderr, report, format); err != nil {
		return err
	}
	return nil
}

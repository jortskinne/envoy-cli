package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

func init() {
	var outputFormat string
	var outputFile string
	var overrideOS bool
	var failMissing bool
	var write bool

	cmd := &cobra.Command{
		Use:   "resolve <file>",
		Short: "Resolve env file values against OS environment variables",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResolve(args[0], outputFormat, outputFile, overrideOS, failMissing, write)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text or json")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write report to file")
	cmd.Flags().BoolVar(&overrideOS, "override", false, "Override file values with OS env vars")
	cmd.Flags().BoolVar(&failMissing, "fail-missing", false, "Exit with error if any key has no value")
	cmd.Flags().BoolVar(&write, "write", false, "Write resolved values back to file")

	rootCmd.AddCommand(cmd)
}

func runResolve(file, format, outputFile string, overrideOS, failMissing, write bool) error {
	entries, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultResolveOptions()
	opts.OverrideWithOS = overrideOS
	opts.FailOnMissing = failMissing

	results, err := parser.Resolve(entries, opts)
	if err != nil {
		return err
	}

	if write {
		resolved := make([]parser.Entry, len(results))
		for i, r := range results {
			resolved[i] = r.Entry
		}
		if err := parser.WriteEnvFile(file, resolved); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		fmt.Println("Resolved values written to", file)
		return nil
	}

	out := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	return parser.WriteResolveReport(out, results, format)
}

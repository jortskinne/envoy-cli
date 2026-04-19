package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
	"envoy-cli/internal/validator"
)

func init() {
	var requiredKeys string
	var outputFormat string
	var outputFile string
	var disallowEmpty bool
	var enforceUpperSnake bool

	validateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate a .env file against rules",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0], requiredKeys, outputFormat, outputFile, disallowEmpty, enforceUpperSnake)
		},
	}

	validateCmd.Flags().StringVarP(&requiredKeys, "required", "r", "", "Comma-separated list of required keys")
	validateCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text or json")
	validateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write report to file instead of stdout")
	validateCmd.Flags().BoolVar(&disallowEmpty, "disallow-empty", false, "Fail if any value is empty")
	validateCmd.Flags().BoolVar(&enforceUpperSnake, "upper-snake", false, "Enforce UPPER_SNAKE_CASE key naming")

	RootCmd.AddCommand(validateCmd)
}

func runValidate(filePath, requiredKeys, outputFormat, outputFile string, disallowEmpty, enforceUpperSnake bool) error {
	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	opts := validator.Options{
		DisallowEmptyValues: disallowEmpty,
		EnforceUpperSnake:   enforceUpperSnake,
	}

	if requiredKeys != "" {
		for _, k := range strings.Split(requiredKeys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.RequiredKeys = append(opts.RequiredKeys, k)
			}
		}
	}

	results := validator.Validate(entries, opts)

	out := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	if err := validator.WriteValidationReport(out, results, outputFormat); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	if hasErrors(results) {
		os.Exit(1)
	}

	return nil
}

// hasErrors returns true if any validation result has a level of "error".
func hasErrors(results []validator.Result) bool {
	for _, r := range results {
		if r.Level == "error" {
			return true
		}
	}
	return false
}

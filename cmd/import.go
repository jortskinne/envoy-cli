package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var importCmd = &cobra.Command{
	Use:   "import <base-file> <src-file>",
	Short: "Import entries from a dotenv or JSON file into a base env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runImport,
}

func init() {
	importCmd.Flags().BoolP("overwrite", "w", false, "Overwrite existing keys with values from src")
	importCmd.Flags().Bool("skip-invalid", false, "Skip invalid lines instead of failing")
	importCmd.Flags().StringP("format", "f", "", "Source format: dotenv or json (default: auto-detect)")
	importCmd.Flags().StringP("output", "o", "", "Write result to file instead of stdout")
	importCmd.Flags().Bool("dry-run", false, "Print result without writing to the base file")
	RootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	basePath := args[0]
	srcPath := args[1]

	overwrite, _ := cmd.Flags().GetBool("overwrite")
	skipInvalid, _ := cmd.Flags().GetBool("skip-invalid")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	base, err := parser.ParseFile(basePath)
	if err != nil {
		return fmt.Errorf("import: cannot parse base file: %w", err)
	}

	opts := parser.DefaultImportOptions()
	opts.Overwrite = overwrite
	opts.SkipInvalid = skipInvalid
	opts.Format = format

	result, err := parser.Import(base, srcPath, opts)
	if err != nil {
		return err
	}

	if dryRun {
		for _, e := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
		}
		return nil
	}

	dest := basePath
	if output != "" {
		dest = output
	}

	if err := parser.WriteEnvFile(dest, result); err != nil {
		return fmt.Errorf("import: failed to write output: %w", err)
	}

	fmt.Fprintf(os.Stderr, "imported %d entries from %q into %q\n", len(result), srcPath, dest)
	return nil
}

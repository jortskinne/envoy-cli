package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var grepCmd = &cobra.Command{
	Use:   "grep <pattern> <file>",
	Short: "Search for entries matching a pattern",
	Args:  cobra.ExactArgs(2),
	RunE:  runGrep,
}

func init() {
	grepCmd.Flags().Bool("keys-only", false, "Search only in keys")
	grepCmd.Flags().Bool("values-only", false, "Search only in values")
	grepCmd.Flags().BoolP("invert", "v", false, "Invert match")
	grepCmd.Flags().Bool("case-sensitive", false, "Case-sensitive matching")
	grepCmd.Flags().StringP("output", "o", "", "Write output to file")
	RootCmd.AddCommand(grepCmd)
}

func runGrep(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	filePath := args[1]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultGrepOptions()
	opts.Pattern = pattern

	keysOnly, _ := cmd.Flags().GetBool("keys-only")
	valuesOnly, _ := cmd.Flags().GetBool("values-only")
	if keysOnly {
		opts.SearchValues = false
	}
	if valuesOnly {
		opts.SearchKeys = false
	}
	opts.Invert, _ = cmd.Flags().GetBool("invert")
	opts.CaseSensitive, _ = cmd.Flags().GetBool("case-sensitive")

	result, err := parser.Grep(entries, opts)
	if err != nil {
		return fmt.Errorf("grep error: %w", err)
	}

	outFile, _ := cmd.Flags().GetString("output")
	w := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	for _, e := range result {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}

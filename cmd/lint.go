package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var lintCmd = &cobra.Command{
	Use:   "lint [file]",
	Short: "Lint an .env file for style and correctness issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	lintCmd.Flags().BoolP("no-empty", "e", true, "Warn on empty values")
	lintCmd.Flags().BoolP("upper-snake", "u", true, "Enforce UPPER_SNAKE_CASE keys")
	lintCmd.Flags().BoolP("no-duplicates", "d", true, "Warn on duplicate keys")
	RootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	noEmpty, _ := cmd.Flags().GetBool("no-empty")
	upperSnake, _ := cmd.Flags().GetBool("upper-snake")
	noDuplicates, _ := cmd.Flags().GetBool("no-duplicates")

	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	opts := parser.LintOptions{
		DisallowEmptyValues: noEmpty,
		EnforceUpperSnake:   upperSnake,
		WarnDuplicateKeys:   noDuplicates,
		DisallowLeadingSpace: true,
	}

	issues := parser.Lint(entries, opts)

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(issues)
	default:
		if len(issues) == 0 {
			fmt.Println("No lint issues found.")
			return nil
		}
		for _, issue := range issues {
			severityLabel := "[WARN]"
			if issue.Severity == parser.LintError {
				severityLabel = "[ERROR]"
			}
			fmt.Fprintf(os.Stdout, "%s %s: %s\n", severityLabel, issue.Key, issue.Message)
		}
		return nil
	}
}

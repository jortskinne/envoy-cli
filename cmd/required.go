package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var requiredCmd = &cobra.Command{
	Use:   "required <file> <KEY1,KEY2,...>",
	Short: "Check that required keys exist and are non-empty in an env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runRequired,
}

func init() {
	requiredCmd.Flags().Bool("allow-empty", false, "Allow keys to be present but empty")
	rootCmd.AddCommand(requiredCmd)
}

func runRequired(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	keys := strings.Split(args[1], ",")

	allowEmpty, _ := cmd.Flags().GetBool("allow-empty")

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	opts := parser.DefaultRequiredOptions()
	opts.AllowEmpty = allowEmpty

	results, err := parser.CheckRequired(entries, keys, opts)

	for _, r := range results {
		status := "OK"
		if !r.Present {
			status = "MISSING"
		} else if r.Empty {
			if allowEmpty {
				status = "EMPTY (allowed)"
			} else {
				status = "EMPTY"
			}
		}
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", r.Key, status)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}

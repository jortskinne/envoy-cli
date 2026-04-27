package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var protectCmd = &cobra.Command{
	Use:   "protect <file>",
	Short: "Mark specific keys as protected in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runProtect,
}

func init() {
	protectCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to protect (required)")
	protectCmd.Flags().Bool("allow-empty", false, "Protect keys even if their value is empty")
	protectCmd.Flags().Bool("dry-run", false, "Print result without writing to disk")
	protectCmd.Flags().String("output", "", "Write output to file instead of stdout")
	rootCmd.AddCommand(protectCmd)
}

func runProtect(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	keys, _ := cmd.Flags().GetStringSlice("keys")
	allowEmpty, _ := cmd.Flags().GetBool("allow-empty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")

	if len(keys) == 0 {
		return fmt.Errorf("--keys is required")
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filePath, err)
	}

	opts := parser.DefaultProtectOptions()
	opts.Keys = keys
	opts.AllowEmpty = allowEmpty
	opts.DryRun = dryRun

	out, res, err := parser.Protect(entries, opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Protected: %s\n", strings.Join(res.Protected, ", "))
	if len(res.Skipped) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Skipped:   %s\n", strings.Join(res.Skipped, ", "))
	}

	if dryRun {
		return parser.WriteEnvFile(os.Stdout, out)
	}

	dest := filePath
	if output != "" {
		dest = output
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	return parser.WriteEnvFile(f, out)
}

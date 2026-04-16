package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <source> <target>",
	Short: "Promote keys from one .env file into another",
	Args:  cobra.ExactArgs(2),
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing keys in target")
	promoteCmd.Flags().StringSliceP("keys", "k", nil, "Comma-separated list of keys to promote")
	promoteCmd.Flags().BoolP("dry-run", "n", false, "Preview changes without writing")
	RootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]
	targetPath := args[1]

	overwrite, _ := cmd.Flags().GetBool("overwrite")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	src, err := parser.ParseFile(sourcePath)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}

	dst, err := parser.ParseFile(targetPath)
	if err != nil {
		return fmt.Errorf("reading target: %w", err)
	}

	opts := parser.DefaultPromoteOptions()
	opts.Overwrite = overwrite
	opts.Keys = keys
	opts.DryRun = dryRun

	res, err := parser.Promote(src, dst, opts)
	if err != nil {
		return err
	}

	if len(res.Added) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "promoted: %s\n", strings.Join(res.Added, ", "))
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "skipped:  %s\n", strings.Join(res.Skipped, ", "))
	}

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: no changes written)")
		return nil
	}

	if err := parser.WriteEnvFile(targetPath, res.Merged); err != nil {
		return fmt.Errorf("writing target: %w", err)
	}

	fmt.Fprintln(os.Stderr, "target updated successfully")
	return nil
}

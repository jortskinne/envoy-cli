package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var renameDryRun bool
var renameIgnoreMissing bool

var renameCmd = &cobra.Command{
	Use:   "rename <file> <old-key> <new-key>",
	Short: "Rename a key in an .env file",
	Args:  cobra.ExactArgs(3),
	RunE:  runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().BoolVar(&renameDryRun, "dry-run", false, "Preview rename without writing changes")
	renameCmd.Flags().BoolVar(&renameIgnoreMissing, "ignore-missing", false, "Do not error if old key is absent")
}

func runRename(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	oldKey := args[1]
	newKey := args[2]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultRenameOptions()
	opts.DryRun = renameDryRun
	opts.ErrorIfMissing = !renameIgnoreMissing

	updated, result, err := parser.Rename(entries, oldKey, newKey, opts)
	if err != nil {
		return err
	}

	if result.Skipped {
		fmt.Fprintf(os.Stderr, "skipped: %s\n", result.Reason)
		return nil
	}

	if renameDryRun {
		fmt.Printf("[dry-run] would rename %q → %q\n", oldKey, newKey)
		return nil
	}

	if err := parser.WriteEnvFile(filePath, updated); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Printf("renamed %q → %q in %s\n", result.OldKey, result.NewKey, filePath)
	return nil
}

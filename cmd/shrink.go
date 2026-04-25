package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var shrinkCmd = &cobra.Command{
	Use:   "shrink <file>",
	Short: "Reduce an env file by removing comments, empty values, and duplicate keys",
	Args:  cobra.ExactArgs(1),
	RunE:  runShrink,
}

func init() {
	shrinkCmd.Flags().Bool("keep-comments", false, "Preserve comment-only lines")
	shrinkCmd.Flags().Bool("remove-empty", false, "Remove entries with empty values")
	shrinkCmd.Flags().Bool("keep-dupes", false, "Do not deduplicate keys")
	shrinkCmd.Flags().Bool("no-trim", false, "Do not trim whitespace from values")
	shrinkCmd.Flags().StringP("output", "o", "", "Write result to file instead of stdout")
	shrinkCmd.Flags().Bool("dry-run", false, "Print removed count without writing output")
	RootCmd.AddCommand(shrinkCmd)
}

func runShrink(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultShrinkOptions()

	if v, _ := cmd.Flags().GetBool("keep-comments"); v {
		opts.RemoveComments = false
	}
	if v, _ := cmd.Flags().GetBool("remove-empty"); v {
		opts.RemoveEmpty = true
	}
	if v, _ := cmd.Flags().GetBool("keep-dupes"); v {
		opts.DedupeKeys = false
	}
	if v, _ := cmd.Flags().GetBool("no-trim"); v {
		opts.TrimValues = false
	}

	result, removed := parser.Shrink(entries, opts)

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "shrink: %d entr(ies) would be removed (%d remaining)\n", removed, len(result))
		return nil
	}

	outPath, _ := cmd.Flags().GetString("output")
	var dest *os.File
	if outPath != "" {
		dest, err = os.Create(outPath)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	if err := parser.WriteEnvFile(dest, result); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "shrink: removed %d entr(ies)\n", removed)
	return nil
}

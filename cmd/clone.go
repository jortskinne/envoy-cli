package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/parser"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <src> <dst>",
	Short: "Clone entries from one .env file into another",
	Args:  cobra.ExactArgs(2),
	RunE:  runClone,
}

func init() {
	cloneCmd.Flags().String("prefix", "", "Only clone keys matching this prefix")
	cloneCmd.Flags().Bool("strip-prefix", false, "Strip the prefix from cloned keys")
	cloneCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in destination")
	cloneCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	RootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	prefix, _ := cmd.Flags().GetString("prefix")
	stripPrefix, _ := cmd.Flags().GetBool("strip-prefix")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	srcEntries, err := parser.ParseFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}

	var dstEntries []parser.Entry
	if _, statErr := os.Stat(dstPath); statErr == nil {
		dstEntries, err = parser.ParseFile(dstPath)
		if err != nil {
			return fmt.Errorf("reading dst: %w", err)
		}
	}

	opts := parser.DefaultCloneOptions()
	opts.Prefix = prefix
	opts.StripPrefix = stripPrefix
	opts.Overwrite = overwrite

	result, count, err := parser.Clone(dstEntries, srcEntries, opts)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Would clone %d entries into %s\n", count, dstPath)
		return nil
	}

	if err := parser.WriteEnvFile(dstPath, result); err != nil {
		return fmt.Errorf("writing dst: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Cloned %d entries into %s\n", count, dstPath)
	return nil
}

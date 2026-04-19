package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/parser"
)

var setkeyOverwrite bool
var setkeyDryRun bool

var setkeyCmd = &cobra.Command{
	Use:   "setkey <file> <KEY> <value>",
	Short: "Add or update a single key in an env file",
	Args:  cobra.ExactArgs(3),
	RunE:  runSetKey,
}

func init() {
	setkeyCmd.Flags().BoolVar(&setkeyOverwrite, "overwrite", false, "Overwrite existing key")
	setkeyCmd.Flags().BoolVar(&setkeyDryRun, "dry-run", false, "Print result without writing")
	rootCmd.AddCommand(setkeyCmd)
}

func runSetKey(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	key := args[1]
	value := args[2]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := parser.DefaultSetKeyOptions()
	opts.Overwrite = setkeyOverwrite
	opts.DryRun = setkeyDryRun

	result, err := parser.SetKey(entries, key, value, opts)
	if err != nil {
		return err
	}

	if setkeyDryRun {
		for _, e := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
		}
		return nil
	}

	if err := parser.WriteEnvFile(filePath, result); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Set %s in %s\n", key, filePath)
	return nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
	"envoy-cli/internal/syncer"
)

var (
	syncOverwrite  bool
	syncFullSync   bool
	syncOutputFile string
	syncDryRun     bool
)

var syncCmd = &cobra.Command{
	Use:   "sync [base] [target]",
	Short: "Sync keys from a base .env file into a target .env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncOverwrite, "overwrite", false, "Overwrite existing keys in target with values from base")
	syncCmd.Flags().BoolVar(&syncFullSync, "full", false, "Remove keys in target that are not present in base")
	syncCmd.Flags().StringVarP(&syncOutputFile, "output", "o", "", "Write synced result to file instead of stdout")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Print what would be written without making changes")
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	baseEntries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse base file %q: %w", args[0], err)
	}

	targetEntries, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("failed to parse target file %q: %w", args[1], err)
	}

	opts := syncer.DefaultSyncOptions()
	opts.Overwrite = syncOverwrite
	opts.FullSync = syncFullSync

	synced := syncer.Sync(baseEntries, targetEntries, opts)

	if syncDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "# Dry run — no files written")
		for _, e := range synced {
			if e.Comment != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", e.Comment)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	}

	dest := args[1]
	if syncOutputFile != "" {
		dest = syncOutputFile
	}

	if err := syncer.WriteEnvFile(dest, synced); err != nil {
		return fmt.Errorf("failed to write synced env file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Synced %d entries to %s\n", len(synced), dest)
	return nil
}

func init() { _ = os.Getenv } // ensure os import used

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var patchCmd = &cobra.Command{
	Use:   "patch <file>",
	Short: "Apply set/delete/rename operations to an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runPatch,
}

func init() {
	patchCmd.Flags().StringArrayP("set", "s", nil, "Set a key: KEY=VALUE")
	patchCmd.Flags().StringArrayP("delete", "d", nil, "Delete a key")
	patchCmd.Flags().StringArrayP("rename", "r", nil, "Rename a key: OLD=NEW")
	patchCmd.Flags().Bool("ignore-missing", false, "Ignore missing keys for delete/rename")
	patchCmd.Flags().Bool("dry-run", false, "Print result without writing")
	RootCmd.AddCommand(patchCmd)
}

func runPatch(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var ops []parser.PatchOperation

	if sets, _ := cmd.Flags().GetStringArray("set"); len(sets) > 0 {
		for _, s := range sets {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("--set value must be KEY=VALUE, got: %s", s)
			}
			ops = append(ops, parser.PatchOperation{Action: "set", Key: parts[0], Value: parts[1]})
		}
	}

	if deletes, _ := cmd.Flags().GetStringArray("delete"); len(deletes) > 0 {
		for _, k := range deletes {
			ops = append(ops, parser.PatchOperation{Action: "delete", Key: k})
		}
	}

	if renames, _ := cmd.Flags().GetStringArray("rename"); len(renames) > 0 {
		for _, r := range renames {
			parts := strings.SplitN(r, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("--rename value must be OLD=NEW, got: %s", r)
			}
			ops = append(ops, parser.PatchOperation{Action: "rename", Key: parts[0], NewKey: parts[1]})
		}
	}

	ignoreMissing, _ := cmd.Flags().GetBool("ignore-missing")
	opts := parser.PatchOptions{IgnoreMissing: ignoreMissing}

	result, err := parser.Patch(entries, ops, opts)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		for _, e := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	}
	return parser.WriteEnvFile(filePath, result)
}

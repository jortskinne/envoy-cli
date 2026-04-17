package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var groupCmd = &cobra.Command{
	Use:   "group <file>",
	Short: "Group env entries by key prefix and write with section headers",
	Args:  cobra.ExactArgs(1),
	RunE:  runGroup,
}

func init() {
	groupCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	groupCmd.Flags().StringP("separator", "s", "_", "Key prefix separator")
	groupCmd.Flags().Bool("no-headers", false, "Omit section comment headers")
	RootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	output, _ := cmd.Flags().GetString("output")
	separator, _ := cmd.Flags().GetString("separator")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	opts := parser.GroupOptions{
		Separator:     separator,
		CommentHeader: !noHeaders,
	}

	grouped := parser.GroupByPrefix(entries, opts)
	flat := parser.FlattenGrouped(grouped, !noHeaders)

	var dest *os.File
	if output != "" {
		dest, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	return parser.WriteEnvFile(dest, flat)
}

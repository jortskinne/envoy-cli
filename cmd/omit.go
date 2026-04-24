package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var omitCmd = &cobra.Command{
	Use:   "omit <file>",
	Short: "Omit specific keys from an env file by key name, prefix, or sensitivity",
	Args:  cobra.ExactArgs(1),
	RunE:  runOmit,
}

func init() {
	omitCmd.Flags().StringSliceP("keys", "k", nil, "Comma-separated list of keys to omit")
	omitCmd.Flags().StringP("prefix", "p", "", "Omit all keys with this prefix")
	omitCmd.Flags().Bool("sensitive", false, "Omit all sensitive keys (e.g. PASSWORD, SECRET, TOKEN)")
	omitCmd.Flags().StringP("output", "o", "", "Write result to file instead of stdout")
	omitCmd.Flags().Bool("dry-run", false, "Print result without writing to disk")
	rootCmd.AddCommand(omitCmd)
}

func runOmit(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	prefix, _ := cmd.Flags().GetString("prefix")
	sensitive, _ := cmd.Flags().GetBool("sensitive")
	output, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	opts := parser.DefaultOmitOptions()
	opts.Keys = keys
	opts.Prefix = prefix
	opts.OmitSensitive = sensitive

	result, err := parser.Omit(entries, opts)
	if err != nil {
		return fmt.Errorf("omit failed: %w", err)
	}

	if dryRun || output == "" {
		var sb strings.Builder
		for _, e := range result {
			if e.Comment != "" {
				sb.WriteString("# " + e.Comment + "\n")
			}
			sb.WriteString(e.Key + "=" + e.Value + "\n")
		}
		fmt.Print(sb.String())
		if dryRun {
			return nil
		}
	}

	if output != "" {
		if err := parser.WriteEnvFile(output, result); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d entries to %s\n", len(result), output)
	}

	return nil
}

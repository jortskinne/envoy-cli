package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var transformCmd = &cobra.Command{
	Use:   "transform <file>",
	Short: "Apply value transformations to .env entries",
	Args:  cobra.ExactArgs(1),
	RunE:  runTransform,
}

func init() {
	transformCmd.Flags().StringSlice("keys", nil, "Comma-separated keys to transform (default: all)")
	transformCmd.Flags().Bool("uppercase", false, "Convert values to uppercase")
	transformCmd.Flags().Bool("lowercase", false, "Convert values to lowercase")
	transformCmd.Flags().Bool("trim", true, "Trim whitespace from values")
	transformCmd.Flags().String("output", "", "Write result to file instead of stdout")
	transformCmd.Flags().Bool("dry-run", false, "Print result without writing")
	RootCmd.AddCommand(transformCmd)
}

func runTransform(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	uppercase, _ := cmd.Flags().GetBool("uppercase")
	lowercase, _ := cmd.Flags().GetBool("lowercase")
	trim, _ := cmd.Flags().GetBool("trim")
	output, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	opts := parser.DefaultTransformOptions()
	opts.Keys = keys
	opts.Uppercase = uppercase
	opts.Lowercase = lowercase
	opts.TrimSpace = trim

	result := parser.Transform(entries, opts)

	if dryRun || output == "" {
		for _, e := range result {
			if e.Comment != "" {
				fmt.Fprintf(os.Stdout, "# %s\n", e.Comment)
			}
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		if dryRun {
			return nil
		}
	}

	if output != "" {
		var sb strings.Builder
		for _, e := range result {
			if e.Comment != "" {
				sb.WriteString("# " + e.Comment + "\n")
			}
			sb.WriteString(e.Key + "=" + e.Value + "\n")
		}
		if err := os.WriteFile(output, []byte(sb.String()), 0644); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return nil
}

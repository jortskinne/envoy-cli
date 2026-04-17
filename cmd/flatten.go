package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/parser"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten <file>",
	Short: "Flatten nested keys using a separator into dot-notation",
	Args:  cobra.ExactArgs(1),
	RunE:  runFlatten,
}

func init() {
	flattenCmd.Flags().String("separator", "__", "Key separator to split on")
	flattenCmd.Flags().String("prefix", "", "Only flatten keys with this prefix")
	flattenCmd.Flags().Bool("lowercase", false, "Convert resulting keys to lowercase")
	flattenCmd.Flags().Bool("dry-run", false, "Print result without writing")
	flattenCmd.Flags().String("output", "", "Write result to file instead of stdout")
	rootCmd.AddCommand(flattenCmd)
}

func runFlatten(cmd *cobra.Command, args []string) error {
	path := args[0]

	entries, err := parser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	sep, _ := cmd.Flags().GetString("separator")
	prefix, _ := cmd.Flags().GetString("prefix")
	lower, _ := cmd.Flags().GetBool("lowercase")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")

	opts := parser.DefaultFlattenOptions()
	opts.Separator = sep
	opts.Prefix = prefix
	opts.Lowercase = lower

	result := parser.Flatten(entries, opts)

	if dryRun || output == "" {
		w := os.Stdout
		if output != "" {
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()
			w = f
		}
		for _, e := range result {
			if e.Comment != "" {
				fmt.Fprintf(w, "# %s\n", e.Comment)
			}
			fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	}

	return parser.WriteEnvFile(output, result)
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var defaultsCmd = &cobra.Command{
	Use:   "defaults <file>",
	Short: "Apply default values to a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runDefaults,
}

func init() {
	defaultsCmd.Flags().StringArrayP("set", "s", nil, "Default key=value pairs to apply (repeatable)")
	defaultsCmd.Flags().BoolP("overwrite", "w", false, "Overwrite existing values with defaults")
	defaultsCmd.Flags().Bool("allow-empty", false, "Allow empty default values to be applied")
	defaultsCmd.Flags().StringP("output", "o", "", "Write result to file instead of stdout")
	defaultsCmd.Flags().Bool("dry-run", false, "Print result without writing")
	rootCmd.AddCommand(defaultsCmd)
}

func runDefaults(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	sets, _ := cmd.Flags().GetStringArray("set")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	allowEmpty, _ := cmd.Flags().GetBool("allow-empty")
	output, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	defaults := make(map[string]string, len(sets))
	for _, s := range sets {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set value (expected key=value)", s)
		}
		defaults[parts[0]] = parts[1]
	}

	opts := parser.DefaultDefaultsOptions()
	opts.Overwrite = overwrite
	opts.SkipEmpty = !allowEmpty

	result := parser.ApplyDefaults(entries, defaults, opts)

	if dryRun || output == "" {
		w := os.Stdout
		if output != "" && !dryRun {
			w, err = os.Create(output)
			if err != nil {
				return err
			}
			defer w.Close()
		}
		return parser.WriteEnvFile(w, result)
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return parser.WriteEnvFile(f, result)
}

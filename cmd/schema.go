package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var schemaCmd = &cobra.Command{
	Use:   "schema [env-file] [schema-file]",
	Short: "Validate an env file against a JSON schema and apply defaults",
	Args:  cobra.ExactArgs(2),
	RunE:  runSchema,
}

var applyDefaults bool

func init() {
	schemaCmd.Flags().BoolVar(&applyDefaults, "apply-defaults", false, "Write missing keys with default values back to the env file")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	envPath := args[0]
	schemaPath := args[1]

	entries, err := parser.ParseFile(envPath)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	schema, err := parser.LoadSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	required := schema.RequiredKeys()
	sort.Strings(required)

	existing := make(map[string]string, len(entries))
	for _, e := range entries {
		existing[e.Key] = e.Value
	}

	if missing := missingKeys(required, existing); len(missing) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Missing required keys:\n")
		for _, k := range missing {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", k)
		}
		os.Exit(1)
	}

	if applyDefaults {
		updated := schema.ApplyDefaults(entries)
		if err := parser.WriteEnvFileFromEntries(envPath, updated); err != nil {
			return fmt.Errorf("writing env file: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Defaults applied and file updated.")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Schema validation passed.")
	return nil
}

// missingKeys returns the keys from required that are not present in existing.
func missingKeys(required []string, existing map[string]string) []string {
	var missing []string
	for _, k := range required {
		if _, ok := existing[k]; !ok {
			missing = append(missing, k)
		}
	}
	return missing
}

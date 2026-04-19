package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var envmapCmd = &cobra.Command{
	Use:   "envmap <file>",
	Short: "Convert an .env file to a key=value map and display or export it",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvmap,
}

func init() {
	envmapCmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	envmapCmd.Flags().StringP("filter-prefix", "p", "", "Only include keys with this prefix")
	RootCmd.AddCommand(envmapCmd)
}

func runEnvmap(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	outFmt, _ := cmd.Flags().GetString("output")
	prefix, _ := cmd.Flags().GetString("filter-prefix")

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	m := parser.ToMap(entries)

	if prefix != "" {
		m = parser.FilterMap(m, func(k, _ string) bool {
			return len(k) >= len(prefix) && k[:len(prefix)] == prefix
		})
	}

	switch outFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m)
	default:
		keys := parser.FromMap(m)
		for _, e := range keys {
			fmt.Printf("%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var castCmd = &cobra.Command{
	Use:   "cast <file>",
	Short: "Infer and normalize value types in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runCast,
}

func init() {
	castCmd.Flags().StringSliceP("keys", "k", nil, "Comma-separated list of keys to cast (default: all)")
	castCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	castCmd.Flags().BoolP("strict", "s", false, "Fail on cast errors")
	RootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	file := args[0]
	keys, _ := cmd.Flags().GetStringSlice("keys")
	format, _ := cmd.Flags().GetString("format")
	strict, _ := cmd.Flags().GetBool("strict")

	entries, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	opts := parser.DefaultCastOptions()
	opts.Keys = keys
	opts.StrictMode = strict

	results, err := parser.Cast(entries, opts)
	if err != nil {
		return err
	}

	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	default:
		if len(results) == 0 {
			fmt.Println("No entries to cast.")
			return nil
		}
		fmt.Printf("%-30s %-10s %s\n", "KEY", "TYPE", "NORMALIZED")
		fmt.Println(strings.Repeat("-", 60))
		for _, r := range results {
			fmt.Printf("%-30s %-10s %s\n", r.Key, r.InferredType, r.Normalized)
		}
	}
	return nil
}

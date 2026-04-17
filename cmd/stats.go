package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var statsCmd = &cobra.Command{
	Use:   "stats <file>",
	Short: "Show statistics about a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runStats,
}

func init() {
	statsCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	statsCmd.Flags().StringP("output", "o", "", "Write output to file")
	RootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	outFile, _ := cmd.Flags().GetString("output")

	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return err
	}

	opts := parser.DefaultStatsOptions()
	stats := parser.ComputeStats(entries, opts)

	w := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	return parser.WriteStatsReport(w, stats, format)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/internal/audit"
	"github.com/envoy-cli/internal/parser"
	"github.com/spf13/cobra"
)

var auditFormat string
var auditOutput string

var auditCmd = &cobra.Command{
	Use:   "audit <base> <target>",
	Short: "Generate an audit log of changes between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().StringVarP(&auditFormat, "format", "f", "text", "Output format: text or json")
	auditCmd.Flags().StringVarP(&auditOutput, "output", "o", "", "Write output to file instead of stdout")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	baseEntries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}

	targetEntries, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading target file: %w", err)
	}

	baseMap := make(map[string]parser.Entry, len(baseEntries))
	for _, e := range baseEntries {
		baseMap[e.Key] = e
	}
	targetMap := make(map[string]parser.Entry, len(targetEntries))
	for _, e := range targetEntries {
		targetMap[e.Key] = e
	}

	opts := audit.DefaultAuditOptions()
	log := audit.Build(baseMap, targetMap, opts)

	w := cmd.OutOrStdout()
	if auditOutput != "" {
		f, err := os.Create(auditOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return audit.WriteAuditReport(log, auditFormat, w)
}

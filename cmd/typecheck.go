package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var typecheckCmd = &cobra.Command{
	Use:   "typecheck [file]",
	Short: "Validate value types in a .env file against declared rules",
	Args:  cobra.ExactArgs(1),
	RunE:  runTypecheck,
}

func init() {
	typecheckCmd.Flags().StringSlice("rule", nil, "Type rules in KEY:TYPE or KEY:regex:PATTERN format")
	typecheckCmd.Flags().String("format", "text", "Output format: text or json")
	typecheckCmd.Flags().String("output", "", "Write output to file")
	rootCmd.AddCommand(typecheckCmd)
}

func runTypecheck(cmd *cobra.Command, args []string) error {
	rawRules, _ := cmd.Flags().GetStringSlice("rule")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return err
	}

	var rules []parser.TypeRule
	for _, raw := range rawRules {
		parts := strings.SplitN(raw, ":", 3)
		if len(parts) < 2 {
			continue
		}
		rule := parser.TypeRule{Key: parts[0], Type: parts[1]}
		if len(parts) == 3 {
			rule.Pattern = parts[2]
		}
		rules = append(rules, rule)
	}

	issues := parser.TypeCheck(entries, parser.TypeCheckOptions{Rules: rules})

	w := cmd.OutOrStdout()
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	return parser.WriteTypeCheckReport(w, issues, format)
}

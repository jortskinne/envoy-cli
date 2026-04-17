package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var templateCmd = &cobra.Command{
	Use:   "template [env-file] [template-file]",
	Short: "Render a template file using values from an .env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runTemplate,
}

func init() {
	templateCmd.Flags().StringP("output", "o", "", "Write rendered output to file instead of stdout")
	templateCmd.Flags().Bool("strict", true, "Fail on unresolved placeholders")
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	envFile := args[0]
	tmplFile := args[1]

	output, _ := cmd.Flags().GetString("output")
	strict, _ := cmd.Flags().GetBool("strict")

	entries, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	opts := parser.DefaultTemplateOptions()
	opts.StrictMode = strict

	rendered, err := parser.RenderTemplateFile(tmplFile, entries, opts)
	if err != nil {
		return err
	}

	if output != "" {
		if err := os.WriteFile(output, []byte(rendered+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rendered template written to %s\n", output)
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), rendered)
	return nil
}

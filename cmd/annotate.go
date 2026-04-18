package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate <file>",
	Short: "Add or view inline comments on env keys",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnnotate,
}

func init() {
	annotateCmd.Flags().StringArrayP("set", "s", nil, "KEY=comment pairs to annotate")
	annotateCmd.Flags().BoolP("overwrite", "w", false, "Overwrite existing comments")
	annotateCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	annotateCmd.Flags().StringP("output", "o", "", "Write output to file")
	annotateCmd.Flags().BoolP("extract", "e", false, "Extract and display existing annotations")
	RootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	file := args[0]
	entries, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	format, _ := cmd.Flags().GetString("format")
	outPath, _ := cmd.Flags().GetString("output")
	extract, _ := cmd.Flags().GetBool("extract")

	w := os.Stdout
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	if extract {
		annotations := parser.ExtractAnnotations(entries)
		return parser.WriteAnnotateReport(w, annotations, format)
	}

	sets, _ := cmd.Flags().GetStringArray("set")
	annotations := make(map[string]string, len(sets))
	for _, s := range sets {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set value %q, expected KEY=comment", s)
		}
		annotations[parts[0]] = parts[1]
	}

	overwrite, _ := cmd.Flags().GetBool("overwrite")
	opts := parser.DefaultAnnotateOptions()
	opts.Overwrite = overwrite

	result := parser.Annotate(entries, annotations, opts)
	return parser.WriteEnvFile(file, result)
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var tagCmd = &cobra.Command{
	Use:   "tag <file>",
	Short: "Apply or extract labels on env entries via inline comments",
	Args:  cobra.ExactArgs(1),
	RunE:  runTag,
}

func init() {
	tagCmd.Flags().StringSliceP("set", "s", nil, "key=label pairs to apply (e.g. APP_NAME=core)")
	tagCmd.Flags().BoolP("extract", "e", false, "extract existing tags instead of applying")
	tagCmd.Flags().BoolP("overwrite", "w", false, "overwrite existing inline comments")
	tagCmd.Flags().StringP("format", "f", "text", "output format: text|json")
	tagCmd.Flags().StringP("output", "o", "", "write result to file instead of stdout")
	tagCmd.Flags().StringP("separator", "", "#tag:", "inline comment prefix used for tags")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	file := args[0]
	entries, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	extract, _ := cmd.Flags().GetBool("extract")
	format, _ := cmd.Flags().GetString("format")
	outFile, _ := cmd.Flags().GetString("output")
	sep, _ := cmd.Flags().GetString("separator")

	w := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("output: %w", err)
		}
		defer f.Close()
		w = f
	}

	if extract {
		tags := parser.ExtractTags(entries, sep)
		report := parser.BuildTagReport(tags, len(entries))
		return parser.WriteTagReport(w, report, format)
	}

	setPairs, _ := cmd.Flags().GetStringSlice("set")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	tags := make(map[string]string)
	for _, pair := range setPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set pair %q (expected key=label)", pair)
		}
		tags[parts[0]] = parts[1]
	}

	opts := parser.DefaultTagOptions()
	opts.Tags = tags
	opts.Overwrite = overwrite
	opts.Separator = sep

	result, err := parser.Tag(entries, opts)
	if err != nil {
		return err
	}
	return parser.WriteEnvFile(w, result)
}

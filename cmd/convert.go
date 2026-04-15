package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var (
	convertFormat  string
	convertMask    bool
	convertOutput  string
)

var convertCmd = &cobra.Command{
	Use:   "convert <file>",
	Short: "Convert a .env file to another format (dotenv, export, json, yaml)",
	Args:  cobra.ExactArgs(1),
	RunE:  runConvert,
}

func init() {
	convertCmd.Flags().StringVarP(&convertFormat, "format", "f", "dotenv", "Output format: dotenv, export, json, yaml")
	convertCmd.Flags().BoolVar(&convertMask, "mask", false, "Mask sensitive values in output")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "", "Write output to file instead of stdout")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	opts := parser.DefaultConvertOptions()
	opts.Format = parser.OutputFormat(convertFormat)
	opts.MaskSecrets = convertMask

	result, err := parser.Convert(entries, opts)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	if convertOutput != "" {
		if err := os.WriteFile(convertOutput, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Written to %s\n", convertOutput)
		return nil
	}

	fmt.Fprint(cmd.OutOrStdout(), result)
	return nil
}

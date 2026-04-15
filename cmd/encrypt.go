package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <file>",
	Short: "Encrypt sensitive values in an .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runEncrypt,
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt <file>",
	Short: "Decrypt sensitive values in an .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runDecrypt,
}

var encryptOutput string
var encryptPassphrase string

func init() {
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "write result to file instead of stdout")
	encryptCmd.Flags().StringVar(&encryptPassphrase, "passphrase", "", "passphrase for encryption (or set ENVOY_SECRET_KEY)")
	decryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "write result to file instead of stdout")
	decryptCmd.Flags().StringVar(&encryptPassphrase, "passphrase", "", "passphrase for decryption (or set ENVOY_SECRET_KEY)")
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}

func resolvePassphrase() (string, error) {
	if encryptPassphrase != "" {
		return encryptPassphrase, nil
	}
	v := os.Getenv("ENVOY_SECRET_KEY")
	if v == "" {
		return "", fmt.Errorf("passphrase required: use --passphrase or set ENVOY_SECRET_KEY")
	}
	return v, nil
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	pass, err := resolvePassphrase()
	if err != nil {
		return err
	}
	result, err := parser.EncryptEntries(entries, pass, parser.DefaultEncryptOptions())
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	dest := encryptOutput
	if dest == "" {
		dest = args[0]
	}
	return parser.WriteEnvFile(dest, result)
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	pass, err := resolvePassphrase()
	if err != nil {
		return err
	}
	result, err := parser.DecryptEntries(entries, pass)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	dest := encryptOutput
	if dest == "" {
		dest = args[0]
	}
	return parser.WriteEnvFile(dest, result)
}

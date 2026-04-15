package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runConvertCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"convert"}, args...))
	err := rootCmd.Execute()
	// reset for next test
	rootCmd.SetArgs([]string{})
	_ = cobra.CheckErr
	return buf.String(), err
}

func TestConvertCmd_DotEnvOutput(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=envoy\nPORT=8080\n")
	out, err := runConvertCmd(t, f, "--format", "dotenv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME in output, got: %s", out)
	}
}

func TestConvertCmd_JSONOutput(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=envoy\nPORT=8080\n")
	out, err := runConvertCmd(t, f, "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestConvertCmd_ExportOutput(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=envoy\n")
	out, err := runConvertCmd(t, f, "--format", "export")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestConvertCmd_MaskFlag(t *testing.T) {
	f := writeTempEnv(t, "DB_PASSWORD=supersecret\nAPP_NAME=envoy\n")
	out, err := runConvertCmd(t, f, "--format", "dotenv", "--mask")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected password to be masked, got: %s", out)
	}
}

func TestConvertCmd_FileOutput(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=envoy\nPORT=9000\n")
	out := filepath.Join(t.TempDir(), "out.env")
	_, err := runConvertCmd(t, f, "--format", "dotenv", "--output", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME in file output, got: %s", string(data))
	}
}

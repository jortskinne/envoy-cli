package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runValidateCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)

	RootCmd.SetArgs(append([]string{"validate"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestValidateCmd_NoErrors(t *testing.T) {
	f := writeTempEnvValidate(t, "APP_NAME=myapp\nAPP_ENV=production\n")
	_, err := runValidateCmd(t, f)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateCmd_MissingRequiredKey(t *testing.T) {
	f := writeTempEnvValidate(t, "APP_NAME=myapp\n")
	_, err := runValidateCmd(t, "--required", "APP_ENV", f)
	if err == nil {
		t.Fatal("expected error for missing required key, got nil")
	}
}

func TestValidateCmd_JSONOutput(t *testing.T) {
	f := writeTempEnvValidate(t, "APP_NAME=myapp\nAPP_ENV=staging\n")
	out, err := runValidateCmd(t, "--format", "json", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[") {
		t.Errorf("expected JSON array in output, got: %s", out)
	}
}

func TestValidateCmd_FileOutput(t *testing.T) {
	f := writeTempEnvValidate(t, "APP_NAME=myapp\n")
	outFile := filepath.Join(t.TempDir(), "report.txt")

	_, err := runValidateCmd(t, "--output", outFile, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestValidateCmd_UpperSnakeViolation(t *testing.T) {
	f := writeTempEnvValidate(t, "appName=myapp\n")
	_, err := runValidateCmd(t, "--upper-snake", f)
	if err == nil {
		t.Fatal("expected error for non-upper-snake key, got nil")
	}
}

func writeTempEnvValidate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

var _ = func() *cobra.Command { return RootCmd }()

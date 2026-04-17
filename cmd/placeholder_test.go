package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writePlaceholderTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func runPlaceholderCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	RootCmd.SetOut(&buf)
	RootCmd.SetErr(&buf)
	RootCmd.SetArgs(append([]string{"placeholder"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestPlaceholderCmd_DetectsPlaceholders(t *testing.T) {
	env := writePlaceholderTempEnv(t, "API_KEY=<your-api-key>\nAPP_NAME=myapp\n")
	out, err := runPlaceholderCmd(t, env)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestPlaceholderCmd_JSONOutput(t *testing.T) {
	env := writePlaceholderTempEnv(t, "TOKEN=CHANGE_ME\nHOST=localhost\n")
	out, err := runPlaceholderCmd(t, "--format", "json", env)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "TOKEN") {
		t.Errorf("expected TOKEN in JSON output, got: %s", out)
	}
}

func TestPlaceholderCmd_FileOutput(t *testing.T) {
	env := writePlaceholderTempEnv(t, "SECRET=TODO\n")
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "out.txt")
	_, err := runPlaceholderCmd(t, "--output", outFile, env)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "SECRET") {
		t.Errorf("expected SECRET in file output, got: %s", string(data))
	}
}

func TestPlaceholderCmd_NoPlaceholders(t *testing.T) {
	env := writePlaceholderTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	out, err := runPlaceholderCmd(t, env)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No placeholder") {
		t.Errorf("expected no-placeholder message, got: %s", out)
	}
}

var _ = cobra.Command{}

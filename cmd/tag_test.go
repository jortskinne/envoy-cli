package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

func writeTagTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runTagCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"tag"}, args...))
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

func TestTagCmd_ExtractText(t *testing.T) {
	path := writeTagTempEnv(t, "APP_NAME=myapp #tag:core\nDB_PASS=secret\n")
	out, err := runTagCmd(t, "--extract", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "core") {
		t.Errorf("expected 'core' in output, got: %s", out)
	}
}

func TestTagCmd_ExtractJSON(t *testing.T) {
	path := writeTagTempEnv(t, "APP_NAME=myapp #tag:core\n")
	out, err := runTagCmd(t, "--extract", "--format", "json", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"tagged\"") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestTagCmd_SetTag(t *testing.T) {
	path := writeTagTempEnv(t, "APP_NAME=myapp\nDB_PASS=secret\n")
	out, err := runTagCmd(t, "--set", "APP_NAME=core", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "#tag:core") {
		t.Errorf("expected tagged output, got: %s", out)
	}
}

func TestTagCmd_FileOutput(t *testing.T) {
	path := writeTagTempEnv(t, "APP_NAME=myapp\n")
	out := filepath.Join(t.TempDir(), "out.env")
	_, err := runTagCmd(t, "--set", "APP_NAME=core", "--output", out, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "#tag:core") {
		t.Errorf("expected tagged content in file, got: %s", string(data))
	}
}

func TestTagCmd_MissingArgs(t *testing.T) {
	_ = parser.EnvEntry{} // ensure import used
	_ = cobra.Command{}
	_, err := runTagCmd(t)
	if err == nil {
		t.Error("expected error for missing args")
	}
}

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDiffCmd_TextOutput(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nSECRET_KEY=abc123\n")
	target := writeTempEnv(t, "APP_NAME=envoy\nSECRET_KEY=xyz789\nNEW_VAR=hello\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"diff", "--format", "text", "--mask=false", base, target})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "NEW_VAR") {
		t.Errorf("expected NEW_VAR in output, got:\n%s", out)
	}
}

func TestDiffCmd_JSONOutput(t *testing.T) {
	base := writeTempEnv(t, "DB_HOST=localhost\n")
	target := writeTempEnv(t, "DB_HOST=prod.db\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"diff", "--format", "json", "--mask=false", base, target})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\"key\"") {
		t.Errorf("expected JSON output with 'key' field, got:\n%s", out)
	}
}

func TestDiffCmd_FileOutput(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")
	target := writeTempEnv(t, "FOO=baz\n")
	outFile := filepath.Join(t.TempDir(), "result.txt")

	rootCmd.SetArgs([]string{"diff", "--format", "text", "--output", outFile, base, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "FOO") {
		t.Errorf("expected FOO in file output, got:\n%s", string(data))
	}
}

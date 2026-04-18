package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTypecheckTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runTypecheckCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"typecheck"}, args...))
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func TestTypecheckCmd_ValidInt(t *testing.T) {
	p := writeTypecheckTempEnv(t, "PORT=8080\n")
	out, err := runTypecheckCmd([]string{p, "--rule", "PORT:int"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "passed") {
		t.Errorf("expected pass message, got: %s", out)
	}
}

func TestTypecheckCmd_InvalidInt(t *testing.T) {
	p := writeTypecheckTempEnv(t, "PORT=notanumber\n")
	out, err := runTypecheckCmd([]string{p, "--rule", "PORT:int"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "PORT") || !strings.Contains(out, "FAIL") {
		t.Errorf("expected failure for PORT, got: %s", out)
	}
}

func TestTypecheckCmd_JSONOutput(t *testing.T) {
	p := writeTypecheckTempEnv(t, "ENABLED=yes\n")
	out, err := runTypecheckCmd([]string{p, "--rule", "ENABLED:bool", "--format", "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "\"issues\"") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestTypecheckCmd_FileOutput(t *testing.T) {
	p := writeTypecheckTempEnv(t, "PORT=9000\n")
	outFile := filepath.Join(t.TempDir(), "out.txt")
	_, err := runTypecheckCmd([]string{p, "--rule", "PORT:int", "--output", outFile})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "passed") {
		t.Errorf("expected output in file, got: %s", string(data))
	}
}

func TestTypecheckCmd_MissingArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	_, err := runTypecheckCmd([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

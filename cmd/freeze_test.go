package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

func writeFreezeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeFreezeTempEnv: %v", err)
	}
	return p
}

func runFreezeCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"freeze"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestFreezeCmd_FreezesNamedKey(t *testing.T) {
	p := writeFreezeTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\n")
	_, err := runFreezeCmd(t, p, "--keys", "DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := parser.ParseFile(p)
	var found bool
	for _, e := range entries {
		if e.Key == "DB_HOST" && parser.IsFrozen(e, "@frozen") {
			found = true
		}
	}
	if !found {
		t.Error("expected DB_HOST to be frozen in written file")
	}
}

func TestFreezeCmd_DryRun(t *testing.T) {
	p := writeFreezeTempEnv(t, "APP_ENV=production\n")
	out, err := runFreezeCmd(t, p, "--all", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := parser.ParseFile(p)
	for _, e := range entries {
		if parser.IsFrozen(e, "@frozen") {
			t.Error("dry-run should not modify file")
		}
	}
	_ = out
}

func TestFreezeCmd_JSONOutput(t *testing.T) {
	p := writeFreezeTempEnv(t, "KEY_A=val1\nKEY_B=val2\n")
	out, err := runFreezeCmd(t, p, "--all", "--dry-run", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report parser.FreezeReport
	if jsonErr := json.Unmarshal([]byte(out), &report); jsonErr != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", jsonErr, out)
	}
	if report.Total != 2 {
		t.Errorf("expected total=2, got %d", report.Total)
	}
}

func TestFreezeCmd_FileOutput(t *testing.T) {
	p := writeFreezeTempEnv(t, "LOG_LEVEL=debug\n")
	outFile := filepath.Join(t.TempDir(), "report.txt")

	_, err := runFreezeCmd(t, p, "--keys", "LOG_LEVEL", "--output", outFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(outFile)
	if readErr != nil {
		t.Fatalf("could not read output file: %v", readErr)
	}
	if !strings.Contains(string(data), "LOG_LEVEL") {
		t.Errorf("expected output file to mention LOG_LEVEL, got: %s", data)
	}
}

func TestFreezeCmd_MissingArgs(t *testing.T) {
	RootCmd.SetArgs([]string{"freeze"})
	err := RootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file argument")
	}
	_ = &cobra.Command{} // suppress unused import
}

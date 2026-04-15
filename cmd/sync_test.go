package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runSyncCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"sync"}, args...))
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

func TestSyncCmd_AddsMissingKeys(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=true\nNEW_KEY=hello\n")
	target := writeTempEnv(t, "APP_NAME=old\nDEBUG=false\n")

	out := filepath.Join(t.TempDir(), "synced.env")
	_, err := runSyncCmd(t, []string{base, target, "--output", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	contents := string(data)
	if !strings.Contains(contents, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY to be synced, got:\n%s", contents)
	}
	// Without --overwrite, existing keys should retain target values
	if !strings.Contains(contents, "APP_NAME=old") {
		t.Errorf("expected APP_NAME to retain target value, got:\n%s", contents)
	}
}

func TestSyncCmd_OverwriteFlag(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=true\n")
	target := writeTempEnv(t, "APP_NAME=old\nDEBUG=false\n")

	out := filepath.Join(t.TempDir(), "synced.env")
	_, err := runSyncCmd(t, []string{base, target, "--overwrite", "--output", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	contents := string(data)
	if !strings.Contains(contents, "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME to be overwritten, got:\n%s", contents)
	}
}

func TestSyncCmd_DryRun(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nNEW_KEY=hello\n")
	target := writeTempEnv(t, "APP_NAME=old\n")

	out, err := runSyncCmd(t, []string{base, target, "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Dry run") {
		t.Errorf("expected dry run header in output, got: %s", out)
	}
	if !strings.Contains(out, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY in dry run output, got: %s", out)
	}
}

func TestSyncCmd_MissingArgs(t *testing.T) {
	// Reset cobra state
	rootCmd.SetArgs([]string{"sync"})
	defer func() { rootCmd.SetArgs(nil) }()
	_ = &cobra.Command{} // suppress unused import
	_, err := runSyncCmd(t, []string{})
	if err == nil {
		t.Error("expected error for missing arguments")
	}
}

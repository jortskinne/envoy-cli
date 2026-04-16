package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

func runRenameCmd(args ...string) (string, error) {
	buf := new(strings.Builder)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"rename"}, args...))
	_ = cobra.EnableCommandSorting
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestRenameCmd_Success(t *testing.T) {
	f, _ := os.CreateTemp("", "*.env")
	f.WriteString("DB_HOST=localhost\nDB_PORT=5432\n")
	f.Close()
	defer os.Remove(f.Name())

	_, err := runRenameCmd(f.Name(), "DB_HOST", "DATABASE_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := parser.ParseFile(f.Name())
	found := false
	for _, e := range entries {
		if e.Key == "DATABASE_HOST" {
			found = true
		}
	}
	if !found {
		t.Error("expected DATABASE_HOST in renamed file")
	}
}

func TestRenameCmd_DryRun(t *testing.T) {
	f, _ := os.CreateTemp("", "*.env")
	f.WriteString("DB_HOST=localhost\n")
	f.Close()
	defer os.Remove(f.Name())

	out, err := runRenameCmd("--dry-run", f.Name(), "DB_HOST", "DATABASE_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run output, got: %s", out)
	}

	entries, _ := parser.ParseFile(f.Name())
	if entries[0].Key != "DB_HOST" {
		t.Error("dry run should not modify file")
	}
}

func TestRenameCmd_MissingArgs(t *testing.T) {
	_, err := runRenameCmd("only-one-arg")
	if err == nil {
		t.Error("expected error for missing args")
	}
}

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

func runStatsCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	statsCmd.SetOut(buf)
	statsCmd.SetErr(buf)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	statsCmd.ResetFlags()
	statsCmd.Flags().StringP("format", "f", "text", "")
	statsCmd.Flags().StringP("output", "o", "", "")
	statsCmd.RunE = runStats
	_ = cobra.EnableCommandSorting
	RootCmd.SetArgs(append([]string{"stats"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func writeStatsTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestStatsCmd_TextOutput(t *testing.T) {
	f := writeStatsTempEnv(t, "DB_HOST=localhost\nDB_PASSWORD=secret\nAPP_NAME=envoy\nPORT=8080\n")
	out, err := runStatsCmd([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Total keys") {
		t.Errorf("expected 'Total keys' in output, got: %s", out)
	}
}

func TestStatsCmd_JSONOutput(t *testing.T) {
	f := writeStatsTempEnv(t, "API_KEY=abc\nAPI_SECRET=xyz\n")
	out, err := runStatsCmd([]string{"--format", "json", f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Total") {
		t.Errorf("expected JSON with Total field, got: %s", out)
	}
}

func TestStatsCmd_FileOutput(t *testing.T) {
	f := writeStatsTempEnv(t, "DB_HOST=localhost\nDB_PASS=s3cr3t\n")
	outPath := filepath.Join(t.TempDir(), "stats.txt")
	_, err := runStatsCmd([]string{"--output", outPath, f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(outPath)
	if !strings.Contains(string(data), "Total") {
		t.Errorf("expected output file to contain stats, got: %s", string(data))
	}
}

func TestStatsCmd_MissingArgs(t *testing.T) {
	_, err := runStatsCmd([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestComputeStats_EmptyInput(t *testing.T) {
	s := parser.ComputeStats(nil, parser.DefaultStatsOptions())
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

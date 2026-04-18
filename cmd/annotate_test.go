package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

func writeAnnotateTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "annotate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runAnnotateCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"annotate"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestAnnotateCmd_ExtractText(t *testing.T) {
	file := writeAnnotateTempEnv(t, "APP_NAME=envoy # the app name\nSECRET=abc\n")
	_, err := runAnnotateCmd(file, "--extract")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnnotateCmd_ExtractJSON(t *testing.T) {
	file := writeAnnotateTempEnv(t, "APP_NAME=envoy # the app name\n")
	out, err := runAnnotateCmd(file, "--extract", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME") && !strings.Contains(out, "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestAnnotateCmd_SetAnnotation(t *testing.T) {
	file := writeAnnotateTempEnv(t, "APP_ENV=staging\n")
	_, err := runAnnotateCmd(file, "--set", "APP_ENV=deployment environment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := parser.ParseFile(file)
	for _, e := range entries {
		if e.Key == "APP_ENV" && e.Comment != "deployment environment" {
			t.Errorf("expected comment set, got %q", e.Comment)
		}
	}
}

func TestAnnotateCmd_MissingArgs(t *testing.T) {
	cmd := &cobra.Command{}
	err := runAnnotate(cmd, []string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

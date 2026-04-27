package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writePinTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func runPinCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "pin"}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	pinCmd.ResetFlags()
	init_pin()
	pinCmd.SetOut(buf)
	pinCmd.SetErr(buf)
	pinCmd.SetArgs(args)
	err := pinCmd.Execute()
	return buf.String(), err
}

func TestPinCmd_PinsNamedKey(t *testing.T) {
	path := writePinTempEnv(t, "API_KEY=secret123\nDATABASE_URL=postgres://localhost/db\n")
	out := filepath.Join(t.TempDir(), "out.env")

	_, err := runPinCmd(t, path, "--keys", "API_KEY", "--output", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if !strings.Contains(string(data), "pinned") {
		t.Errorf("expected pinned comment in output, got:\n%s", string(data))
	}
}

func TestPinCmd_DryRun(t *testing.T) {
	path := writePinTempEnv(t, "SECRET=abc\nHOST=localhost\n")

	out, err := runPinCmd(t, path, "--keys", "SECRET", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected dry-run output to contain SECRET, got:\n%s", out)
	}
}

func TestPinCmd_JSONOutput(t *testing.T) {
	path := writePinTempEnv(t, "TOKEN=xyz\nPORT=8080\n")
	out := filepath.Join(t.TempDir(), "report.json")

	_, err := runPinCmd(t, path, "--keys", "TOKEN", "--format", "json", "--report", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read report: %v", err)
	}
	if !strings.Contains(string(data), "pinned") {
		t.Errorf("expected JSON report to contain pinned field, got:\n%s", string(data))
	}
}

func TestPinCmd_MissingArgs(t *testing.T) {
	_, err := runPinCmd(t)
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestPinCmd_AlreadyPinned(t *testing.T) {
	// Key already has a # pinned comment — should not duplicate
	path := writePinTempEnv(t, "# pinned\nAPI_KEY=secret\n")
	out := filepath.Join(t.TempDir(), "out.env")

	_, err := runPinCmd(t, path, "--keys", "API_KEY", "--output", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	count := strings.Count(string(data), "# pinned")
	if count != 1 {
		t.Errorf("expected exactly one '# pinned' comment, got %d in:\n%s", count, string(data))
	}
}

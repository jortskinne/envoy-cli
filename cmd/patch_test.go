package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writePatchTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func runPatchCmd(args ...string) (string, error) {
	buf := new(strings.Builder)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"patch"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestPatchCmd_SetKey(t *testing.T) {
	file := writePatchTempEnv(t, "APP_ENV=staging\nDB_HOST=localhost\n")
	_, err := runPatchCmd(file, "--set", "APP_ENV=production", "--dry-run")
	require.NoError(t, err)
}

func TestPatchCmd_DryRun_PrintsResult(t *testing.T) {
	file := writePatchTempEnv(t, "FOO=bar\nBAZ=qux\n")
	out, err := runPatchCmd(file, "--set", "FOO=updated", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "FOO=updated")
	assert.Contains(t, out, "BAZ=qux")
}

func TestPatchCmd_DeleteKey(t *testing.T) {
	file := writePatchTempEnv(t, "FOO=bar\nBAZ=qux\n")
	out, err := runPatchCmd(file, "--delete", "BAZ", "--dry-run")
	require.NoError(t, err)
	assert.NotContains(t, out, "BAZ")
}

func TestPatchCmd_RenameKey(t *testing.T) {
	file := writePatchTempEnv(t, "OLD_KEY=value\n")
	out, err := runPatchCmd(file, "--rename", "OLD_KEY=NEW_KEY", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "NEW_KEY=value")
	assert.NotContains(t, out, "OLD_KEY")
}

func TestPatchCmd_MissingArgs(t *testing.T) {
	_, err := runPatchCmd()
	assert.Error(t, err)
}

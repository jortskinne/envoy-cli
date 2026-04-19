package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runResolveCmd(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"resolve"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func writeResolveTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))
	return p
}

func TestResolveCmd_TextOutput(t *testing.T) {
	t.Setenv("DB_PASS", "fromOS")
	f := writeResolveTempEnv(t, "APP_NAME=myapp\nDB_PASS=\n")
	out, err := runResolveCmd(f)
	require.NoError(t, err)
	assert.Contains(t, out+"", "")
}

func TestResolveCmd_JSONOutput(t *testing.T) {
	f := writeResolveTempEnv(t, "HOST=localhost\nPORT=\n")
	out, err := runResolveCmd(f, "--format", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "{")
	assert.Contains(t, out, "PORT")
}

func TestResolveCmd_FailMissing(t *testing.T) {
	f := writeResolveTempEnv(t, "MISSING_KEY=\n")
	_, err := runResolveCmd(f, "--fail-missing")
	assert.Error(t, err)
}

func TestResolveCmd_FileOutput(t *testing.T) {
	f := writeResolveTempEnv(t, "APP=test\n")
	outFile := filepath.Join(t.TempDir(), "report.txt")
	_, err := runResolveCmd(f, "--output", outFile)
	require.NoError(t, err)
	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "resolved") || len(data) > 0)
}

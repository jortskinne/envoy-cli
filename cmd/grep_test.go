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

func writeGrepTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func runGrepCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"grep"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestGrepCmd_MatchesPattern(t *testing.T) {
	path := writeGrepTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\nAPP_NAME=myapp\n")
	out, err := runGrepCmd("^DB_", path, "--keys-only")
	require.NoError(t, err)
	assert.Contains(t, out, "DB_HOST")
	assert.Contains(t, out, "DB_PASS")
	assert.NotContains(t, out, "APP_NAME")
}

func TestGrepCmd_InvertFlag(t *testing.T) {
	path := writeGrepTempEnv(t, "DB_HOST=localhost\nAPP_NAME=myapp\n")
	out, err := runGrepCmd("^DB_", path, "--keys-only", "--invert")
	require.NoError(t, err)
	assert.NotContains(t, out, "DB_HOST")
	assert.Contains(t, out, "APP_NAME")
}

func TestGrepCmd_FileOutput(t *testing.T) {
	path := writeGrepTempEnv(t, "API_KEY=abc\nLOG_LEVEL=debug\n")
	outPath := filepath.Join(t.TempDir(), "out.env")
	_, err := runGrepCmd("API", path, "--keys-only", "--output", outPath)
	require.NoError(t, err)
	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "API_KEY"))
}

func TestGrepCmd_MissingArgs(t *testing.T) {
	_, err := runGrepCmd("pattern")
	assert.Error(t, err)
}

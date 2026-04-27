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

func writeProtectTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func runProtectCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"protect"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestProtectCmd_MarksKey(t *testing.T) {
	env := writeProtectTempEnv(t, "APP_NAME=envoy\nDB_PASS=secret\n")
	out, err := runProtectCmd(env, "--keys", "APP_NAME")
	require.NoError(t, err)
	assert.Contains(t, out, "Protected: APP_NAME")

	data, _ := os.ReadFile(env)
	assert.Contains(t, string(data), "PROTECTED")
}

func TestProtectCmd_DryRun(t *testing.T) {
	env := writeProtectTempEnv(t, "TOKEN=abc123\n")
	out, err := runProtectCmd(env, "--keys", "TOKEN", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "TOKEN")
}

func TestProtectCmd_FileOutput(t *testing.T) {
	env := writeProtectTempEnv(t, "API_KEY=xyz\n")
	dest := filepath.Join(t.TempDir(), "out.env")
	_, err := runProtectCmd(env, "--keys", "API_KEY", "--output", dest)
	require.NoError(t, err)
	data, _ := os.ReadFile(dest)
	assert.Contains(t, string(data), "API_KEY")
}

func TestProtectCmd_MissingKeys(t *testing.T) {
	env := writeProtectTempEnv(t, "FOO=bar\n")
	_, err := runProtectCmd(env)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "--keys") || err != nil)
}

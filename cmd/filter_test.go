package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runFilterCmd(args ...string) (string, error) {
	buf := new(strings.Builder)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"filter"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func writeFilterTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func TestFilterCmd_ByPrefix(t *testing.T) {
	env := writeFilterTempEnv(t, "APP_NAME=myapp\nAPP_ENV=prod\nDB_HOST=localhost\n")
	out, err := runFilterCmd(env, "--prefix", "APP_")
	require.NoError(t, err)
	assert.Contains(t, out, "APP_NAME")
	assert.Contains(t, out, "APP_ENV")
	assert.NotContains(t, out, "DB_HOST")
}

func TestFilterCmd_Exclude(t *testing.T) {
	env := writeFilterTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nDB_HOST=localhost\n")
	out, err := runFilterCmd(env, "--exclude", "DB_PASSWORD")
	require.NoError(t, err)
	assert.NotContains(t, out, "DB_PASSWORD")
	assert.Contains(t, out, "APP_NAME")
}

func TestFilterCmd_SensitiveOnly(t *testing.T) {
	env := writeFilterTempEnv(t, "APP_NAME=myapp\nAPI_KEY=abc123\nDB_PASSWORD=secret\n")
	out, err := runFilterCmd(env, "--sensitive")
	require.NoError(t, err)
	assert.NotContains(t, out, "APP_NAME")
	assert.Contains(t, out, "API_KEY")
}

func TestFilterCmd_FileOutput(t *testing.T) {
	env := writeFilterTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	outFile := filepath.Join(t.TempDir(), "out.env")
	_, err := runFilterCmd(env, "--keys", "APP_NAME", "--out", outFile)
	require.NoError(t, err)
	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "APP_NAME=myapp")
	assert.NotContains(t, string(data), "DB_HOST")
}

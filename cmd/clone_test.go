package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runCloneCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"clone"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestCloneCmd_ClonesEntries(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dst := filepath.Join(dir, "dst.env")

	require.NoError(t, os.WriteFile(src, []byte("DB_HOST=localhost\nDB_PORT=5432\n"), 0644))
	require.NoError(t, os.WriteFile(dst, []byte("APP_NAME=myapp\n"), 0644))

	out, err := runCloneCmd(src, dst)
	require.NoError(t, err)
	assert.Contains(t, out, "Cloned 2 entries")

	content, _ := os.ReadFile(dst)
	assert.Contains(t, string(content), "DB_HOST=localhost")
	assert.Contains(t, string(content), "APP_NAME=myapp")
}

func TestCloneCmd_DryRun(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dst := filepath.Join(dir, "dst.env")

	require.NoError(t, os.WriteFile(src, []byte("KEY=value\n"), 0644))
	require.NoError(t, os.WriteFile(dst, []byte(""), 0644))

	out, err := runCloneCmd("--dry-run", src, dst)
	require.NoError(t, err)
	assert.Contains(t, out, "Would clone")

	content, _ := os.ReadFile(dst)
	assert.NotContains(t, string(content), "KEY=value")
}

func TestCloneCmd_StripPrefix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dst := filepath.Join(dir, "dst.env")

	require.NoError(t, os.WriteFile(src, []byte("STAGE_DB_HOST=localhost\nSTAGE_DB_PORT=5432\n"), 0644))
	require.NoError(t, os.WriteFile(dst, []byte(""), 0644))

	_, err := runCloneCmd("--prefix=STAGE_", "--strip-prefix", src, dst)
	require.NoError(t, err)

	content, _ := os.ReadFile(dst)
	assert.Contains(t, string(content), "DB_HOST=localhost")
	assert.NotContains(t, string(content), "STAGE_")
}

func TestCloneCmd_MissingArgs(t *testing.T) {
	_, err := runCloneCmd("only-one-arg")
	assert.Error(t, err)
}

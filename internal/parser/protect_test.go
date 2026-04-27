package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeProtectEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "EMPTY_KEY", Value: ""},
		{Key: "ALREADY", Value: "yes", Comment: "PROTECTED"},
	}
}

func TestProtect_MarksKeys(t *testing.T) {
	entries := makeProtectEntries()
	opts := DefaultProtectOptions()
	opts.Keys = []string{"APP_NAME", "DB_PASSWORD"}

	out, res, err := Protect(entries, opts)
	require.NoError(t, err)
	assert.Equal(t, []string{"APP_NAME", "DB_PASSWORD"}, res.Protected)
	assert.Empty(t, res.Skipped)

	for _, e := range out {
		if e.Key == "APP_NAME" || e.Key == "DB_PASSWORD" {
			assert.Equal(t, "PROTECTED", e.Comment)
		}
	}
}

func TestProtect_SkipsEmptyValue(t *testing.T) {
	entries := makeProtectEntries()
	opts := DefaultProtectOptions()
	opts.Keys = []string{"EMPTY_KEY"}

	_, res, err := Protect(entries, opts)
	require.NoError(t, err)
	assert.Contains(t, res.Skipped, "EMPTY_KEY")
	assert.Empty(t, res.Protected)
}

func TestProtect_AllowEmpty(t *testing.T) {
	entries := makeProtectEntries()
	opts := DefaultProtectOptions()
	opts.Keys = []string{"EMPTY_KEY"}
	opts.AllowEmpty = true

	_, res, err := Protect(entries, opts)
	require.NoError(t, err)
	assert.Contains(t, res.Protected, "EMPTY_KEY")
}

func TestProtect_SkipsAlreadyProtected(t *testing.T) {
	entries := makeProtectEntries()
	opts := DefaultProtectOptions()
	opts.Keys = []string{"ALREADY"}

	_, res, err := Protect(entries, opts)
	require.NoError(t, err)
	assert.Contains(t, res.Skipped, "ALREADY")
}

func TestProtect_NoKeysReturnsError(t *testing.T) {
	entries := makeProtectEntries()
	opts := DefaultProtectOptions()

	_, _, err := Protect(entries, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no keys specified")
}

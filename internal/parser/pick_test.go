package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePickEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestPick_ReturnsRequestedKeys(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{"APP_NAME", "DB_PORT"}

	result, err := Pick(entries, opts)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "APP_NAME", result[0].Key)
	assert.Equal(t, "DB_PORT", result[1].Key)
}

func TestPick_PreservesKeyOrder(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{"SECRET_KEY", "APP_ENV", "DB_HOST"}

	result, err := Pick(entries, opts)
	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "SECRET_KEY", result[0].Key)
	assert.Equal(t, "APP_ENV", result[1].Key)
	assert.Equal(t, "DB_HOST", result[2].Key)
}

func TestPick_IgnoresMissingKeysInLenientMode(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{"APP_NAME", "DOES_NOT_EXIST"}

	result, err := Pick(entries, opts)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "APP_NAME", result[0].Key)
}

func TestPick_StrictMode_ErrorOnMissingKey(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{"APP_NAME", "MISSING_KEY"}
	opts.StrictMode = true

	_, err := Pick(entries, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MISSING_KEY")
}

func TestPick_EmptyKeys_ReturnsEmpty(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{}

	result, err := Pick(entries, opts)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestPick_SkipsBlankKeyStrings(t *testing.T) {
	entries := makePickEntries()
	opts := DefaultPickOptions()
	opts.Keys = []string{"", "APP_ENV", ""}

	result, err := Pick(entries, opts)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "APP_ENV", result[0].Key)
}

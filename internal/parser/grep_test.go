package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeGrepEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "API_KEY", Value: "abcdef"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "LOG_LEVEL", Value: "debug"},
	}
}

func TestGrep_MatchesKey(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.SearchValues = false
	opts.Pattern = "^DB_"
	result, err := Grep(entries, opts)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "DB_HOST", result[0].Key)
	assert.Equal(t, "DB_PASSWORD", result[1].Key)
}

func TestGrep_MatchesValue(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.SearchKeys = false
	opts.Pattern = "localhost"
	result, err := Grep(entries, opts)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "DB_HOST", result[0].Key)
}

func TestGrep_CaseInsensitive(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.Pattern = "debug"
	opts.CaseSensitive = false
	result, err := Grep(entries, opts)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "LOG_LEVEL", result[0].Key)
}

func TestGrep_Invert(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.Pattern = "^DB_"
	opts.SearchValues = false
	opts.Invert = true
	result, err := Grep(entries, opts)
	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestGrep_InvalidPattern(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.Pattern = "["
	_, err := Grep(entries, opts)
	assert.Error(t, err)
}

func TestGrep_EmptyResult(t *testing.T) {
	entries := makeGrepEntries()
	opts := DefaultGrepOptions()
	opts.Pattern = "NONEXISTENT"
	result, err := Grep(entries, opts)
	require.NoError(t, err)
	assert.Empty(t, result)
}

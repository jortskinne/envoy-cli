package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeExtractEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestExtract_ByKeys(t *testing.T) {
	opts := DefaultExtractOptions()
	opts.Keys = []string{"DB_HOST", "SECRET_KEY"}
	out, err := Extract(makeExtractEntries(), opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "DB_HOST", out[0].Key)
	assert.Equal(t, "SECRET_KEY", out[1].Key)
}

func TestExtract_ByPrefix(t *testing.T) {
	opts := DefaultExtractOptions()
	opts.Prefix = "APP_"
	out, err := Extract(makeExtractEntries(), opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	for _, e := range out {
		assert.True(t, len(e.Key) > 0)
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	opts := DefaultExtractOptions()
	opts.Prefix = "APP_"
	opts.StripPrefix = true
	out, err := Extract(makeExtractEntries(), opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	keys := []string{out[0].Key, out[1].Key}
	assert.Contains(t, keys, "NAME")
	assert.Contains(t, keys, "ENV")
}

func TestExtract_KeysAndPrefix(t *testing.T) {
	opts := DefaultExtractOptions()
	opts.Keys = []string{"SECRET_KEY"}
	opts.Prefix = "DB_"
	out, err := Extract(makeExtractEntries(), opts)
	require.NoError(t, err)
	assert.Len(t, out, 3)
}

func TestExtract_NoMatch(t *testing.T) {
	opts := DefaultExtractOptions()
	opts.Keys = []string{"NONEXISTENT"}
	out, err := Extract(makeExtractEntries(), opts)
	require.NoError(t, err)
	assert.Empty(t, out)
}

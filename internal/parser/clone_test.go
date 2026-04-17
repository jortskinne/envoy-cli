package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeCloneEntries() (dst, src []Entry) {
	dst = []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
	}
	src = []Entry{
		{Key: "STAGE_DB_HOST", Value: "localhost"},
		{Key: "STAGE_DB_PORT", Value: "5432"},
		{Key: "STAGE_SECRET", Value: "abc123"},
		{Key: "OTHER_KEY", Value: "other"},
	}
	return
}

func TestClone_AppendsNewKeys(t *testing.T) {
	dst, src := makeCloneEntries()
	opts := DefaultCloneOptions()
	result, count, err := Clone(dst, src, opts)
	require.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Len(t, result, 6)
}

func TestClone_FiltersByPrefix(t *testing.T) {
	dst, src := makeCloneEntries()
	opts := DefaultCloneOptions()
	opts.Prefix = "STAGE_"
	result, count, err := Clone(dst, src, opts)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 5)
}

func TestClone_StripPrefix(t *testing.T) {
	dst, src := makeCloneEntries()
	opts := DefaultCloneOptions()
	opts.Prefix = "STAGE_"
	opts.StripPrefix = true
	result, count, err := Clone(dst, src, opts)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
	keys := make([]string, len(result))
	for i, e := range result {
		keys[i] = e.Key
	}
	assert.Contains(t, keys, "DB_HOST")
	assert.Contains(t, keys, "DB_PORT")
	assert.Contains(t, keys, "SECRET")
}

func TestClone_DoesNotOverwriteByDefault(t *testing.T) {
	dst := []Entry{{Key: "DB_HOST", Value: "prod-host"}}
	src := []Entry{{Key: "DB_HOST", Value: "stage-host"}}
	opts := DefaultCloneOptions()
	result, _, err := Clone(dst, src, opts)
	require.NoError(t, err)
	assert.Equal(t, "prod-host", result[0].Value)
}

func TestClone_OverwriteFlag(t *testing.T) {
	dst := []Entry{{Key: "DB_HOST", Value: "prod-host"}}
	src := []Entry{{Key: "DB_HOST", Value: "stage-host"}}
	opts := DefaultCloneOptions()
	opts.Overwrite = true
	result, _, err := Clone(dst, src, opts)
	require.NoError(t, err)
	assert.Equal(t, "stage-host", result[0].Value)
}

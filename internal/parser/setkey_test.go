package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSetKeyEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestSetKey_AddsNewKey(t *testing.T) {
	entries := makeSetKeyEntries()
	opts := DefaultSetKeyOptions()
	out, err := SetKey(entries, "APP_PORT", "8080", opts)
	require.NoError(t, err)
	assert.Len(t, out, 3)
	assert.Equal(t, "APP_PORT", out[2].Key)
	assert.Equal(t, "8080", out[2].Value)
}

func TestSetKey_ExistingKey_NoOverwrite_Error(t *testing.T) {
	entries := makeSetKeyEntries()
	opts := DefaultSetKeyOptions()
	_, err := SetKey(entries, "APP_NAME", "newname", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestSetKey_ExistingKey_Overwrite(t *testing.T) {
	entries := makeSetKeyEntries()
	opts := DefaultSetKeyOptions()
	opts.Overwrite = true
	out, err := SetKey(entries, "APP_NAME", "newname", opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "newname", out[0].Value)
}

func TestSetKey_EmptyKey_Error(t *testing.T) {
	entries := makeSetKeyEntries()
	opts := DefaultSetKeyOptions()
	_, err := SetKey(entries, "", "value", opts)
	require.Error(t, err)
}

func TestSetKey_DoesNotMutateOriginal(t *testing.T) {
	entries := makeSetKeyEntries()
	opts := DefaultSetKeyOptions()
	opts.Overwrite = true
	_, err := SetKey(entries, "APP_NAME", "changed", opts)
	require.NoError(t, err)
	assert.Equal(t, "myapp", entries[0].Value)
}

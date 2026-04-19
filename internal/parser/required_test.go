package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRequiredEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "EMPTY_KEY", Value: ""},
	}
}

func TestCheckRequired_AllPresent(t *testing.T) {
	entries := makeRequiredEntries()
	results, err := CheckRequired(entries, []string{"HOST", "PORT"}, DefaultRequiredOptions())
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.True(t, results[0].Present)
	assert.True(t, results[1].Present)
}

func TestCheckRequired_MissingKey(t *testing.T) {
	entries := makeRequiredEntries()
	_, err := CheckRequired(entries, []string{"HOST", "SECRET"}, DefaultRequiredOptions())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SECRET")
}

func TestCheckRequired_EmptyValueFails(t *testing.T) {
	entries := makeRequiredEntries()
	_, err := CheckRequired(entries, []string{"EMPTY_KEY"}, DefaultRequiredOptions())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EMPTY_KEY")
}

func TestCheckRequired_EmptyValueAllowed(t *testing.T) {
	entries := makeRequiredEntries()
	opts := DefaultRequiredOptions()
	opts.AllowEmpty = true
	results, err := CheckRequired(entries, []string{"EMPTY_KEY"}, opts)
	require.NoError(t, err)
	assert.True(t, results[0].Empty)
}

func TestCheckRequired_EmptyRequired(t *testing.T) {
	entries := makeRequiredEntries()
	results, err := CheckRequired(entries, []string{}, DefaultRequiredOptions())
	require.NoError(t, err)
	assert.Empty(t, results)
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeResolveEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: ""},
		{Key: "PORT", Value: ""},
	}
}

func TestResolve_FillMissingFromOS(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret123")
	entries := makeResolveEntries()
	opts := DefaultResolveOptions()
	results, err := Resolve(entries, opts)
	require.NoError(t, err)
	assert.Equal(t, "myapp", results[0].Entry.Value)
	assert.Equal(t, "file", results[0].Source)
	assert.Equal(t, "secret123", results[1].Entry.Value)
	assert.Equal(t, "os", results[1].Source)
	assert.True(t, results[2].Missing)
	assert.Equal(t, "empty", results[2].Source)
}

func TestResolve_OverrideWithOS(t *testing.T) {
	t.Setenv("APP_NAME", "overridden")
	entries := makeResolveEntries()
	opts := DefaultResolveOptions()
	opts.OverrideWithOS = true
	results, err := Resolve(entries, opts)
	require.NoError(t, err)
	assert.Equal(t, "overridden", results[0].Entry.Value)
	assert.Equal(t, "os", results[0].Source)
}

func TestResolve_FailOnMissing(t *testing.T) {
	entries := makeResolveEntries()
	opts := DefaultResolveOptions()
	opts.FailOnMissing = true
	_, err := Resolve(entries, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
	assert.Contains(t, err.Error(), "PORT")
}

func TestResolve_NoMissingKeys(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
	opts := DefaultResolveOptions()
	opts.FailOnMissing = true
	results, err := Resolve(entries, opts)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.False(t, r.Missing)
	}
}

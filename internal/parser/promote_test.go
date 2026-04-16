package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePromoteEntries(kvs ...string) []EnvEntry {
	var out []EnvEntry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, EnvEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := makePromoteEntries("NEW_KEY", "hello")
	dst := makePromoteEntries("EXISTING", "world")
	res, err := Promote(src, dst, DefaultPromoteOptions())
	require.NoError(t, err)
	assert.Len(t, res.Merged, 2)
	assert.Contains(t, res.Added, "NEW_KEY")
}

func TestPromote_SkipsExistingByDefault(t *testing.T) {
	src := makePromoteEntries("KEY", "new")
	dst := makePromoteEntries("KEY", "old")
	res, err := Promote(src, dst, DefaultPromoteOptions())
	require.NoError(t, err)
	assert.Contains(t, res.Skipped, "KEY")
	assert.Equal(t, "old", res.Merged[0].Value)
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := makePromoteEntries("KEY", "new")
	dst := makePromoteEntries("KEY", "old")
	opts := DefaultPromoteOptions()
	opts.Overwrite = true
	res, err := Promote(src, dst, opts)
	require.NoError(t, err)
	assert.Equal(t, "new", res.Merged[0].Value)
	assert.Contains(t, res.Added, "KEY")
}

func TestPromote_FilterByKeys(t *testing.T) {
	src := makePromoteEntries("A", "1", "B", "2", "C", "3")
	dst := makePromoteEntries("X", "9")
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"A", "C"}
	res, err := Promote(src, dst, opts)
	require.NoError(t, err)
	assert.Len(t, res.Merged, 3)
	keys := make([]string, len(res.Merged))
	for i, e := range res.Merged {
		keys[i] = e.Key
	}
	assert.Contains(t, keys, "A")
	assert.Contains(t, keys, "C")
	assert.NotContains(t, keys, "B")
}

func TestPromote_MissingKeyInSource_Error(t *testing.T) {
	src := makePromoteEntries("A", "1")
	dst := makePromoteEntries("X", "9")
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"MISSING"}
	_, err := Promote(src, dst, opts)
	assert.Error(t, err)
}

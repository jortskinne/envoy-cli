package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSliceEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "GAMMA", Value: "3"},
		{Key: "DELTA", Value: "4"},
		{Key: "EPSILON", Value: "5"},
	}
}

func TestSlice_FullRange(t *testing.T) {
	entries := makeSliceEntries()
	opts := DefaultSliceOptions()
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 5)
	assert.Equal(t, "ALPHA", out[0].Key)
	assert.Equal(t, "EPSILON", out[4].Key)
}

func TestSlice_StartOffset(t *testing.T) {
	entries := makeSliceEntries()
	opts := DefaultSliceOptions()
	opts.Start = 2
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 3)
	assert.Equal(t, "GAMMA", out[0].Key)
}

func TestSlice_StartAndEnd(t *testing.T) {
	entries := makeSliceEntries()
	opts := SliceOptions{Start: 1, End: 3}
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "BETA", out[0].Key)
	assert.Equal(t, "GAMMA", out[1].Key)
}

func TestSlice_EndBeyondLength(t *testing.T) {
	entries := makeSliceEntries()
	opts := SliceOptions{Start: 3, End: 100}
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "DELTA", out[0].Key)
	assert.Equal(t, "EPSILON", out[1].Key)
}

func TestSlice_FilterByKeys(t *testing.T) {
	entries := makeSliceEntries()
	opts := SliceOptions{Keys: []string{"ALPHA", "GAMMA", "EPSILON"}, Start: 0, End: -1}
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 3)
	assert.Equal(t, "ALPHA", out[0].Key)
	assert.Equal(t, "GAMMA", out[1].Key)
	assert.Equal(t, "EPSILON", out[2].Key)
}

func TestSlice_FilterByKeysWithRange(t *testing.T) {
	entries := makeSliceEntries()
	opts := SliceOptions{Keys: []string{"ALPHA", "GAMMA", "EPSILON"}, Start: 1, End: 2}
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, "GAMMA", out[0].Key)
}

func TestSlice_DoesNotMutateOriginal(t *testing.T) {
	entries := makeSliceEntries()
	opts := SliceOptions{Start: 0, End: 2}
	out, err := Slice(entries, opts)
	require.NoError(t, err)
	out[0].Key = "MUTATED"
	assert.Equal(t, "ALPHA", entries[0].Key)
}

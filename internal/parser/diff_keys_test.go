package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeDiffKeyEntries(keys ...string) []EnvEntry {
	out := make([]EnvEntry, len(keys))
	for i, k := range keys {
		out[i] = EnvEntry{Key: k, Value: "val"}
	}
	return out
}

func TestDiffKeys_OnlyInBase(t *testing.T) {
	base := makeDiffKeyEntries("A", "B", "C")
	other := makeDiffKeyEntries("B", "C")
	result := DiffKeys(base, other, DefaultDiffKeysOptions())
	assert.Equal(t, []string{"A"}, result.OnlyInBase)
	assert.Empty(t, result.OnlyInOther)
	assert.ElementsMatch(t, []string{"B", "C"}, result.InBoth)
}

func TestDiffKeys_OnlyInOther(t *testing.T) {
	base := makeDiffKeyEntries("X")
	other := makeDiffKeyEntries("X", "Y", "Z")
	result := DiffKeys(base, other, DefaultDiffKeysOptions())
	assert.Empty(t, result.OnlyInBase)
	assert.ElementsMatch(t, []string{"Y", "Z"}, result.OnlyInOther)
	assert.Equal(t, []string{"X"}, result.InBoth)
}

func TestDiffKeys_InBoth(t *testing.T) {
	base := makeDiffKeyEntries("FOO", "BAR")
	other := makeDiffKeyEntries("FOO", "BAR")
	result := DiffKeys(base, other, DefaultDiffKeysOptions())
	assert.Empty(t, result.OnlyInBase)
	assert.Empty(t, result.OnlyInOther)
	assert.ElementsMatch(t, []string{"BAR", "FOO"}, result.InBoth)
}

func TestDiffKeys_IgnoreCase(t *testing.T) {
	base := makeDiffKeyEntries("foo", "BAR")
	other := makeDiffKeyEntries("FOO", "bar")
	opts := DefaultDiffKeysOptions()
	opts.IgnoreCase = true
	result := DiffKeys(base, other, opts)
	assert.Empty(t, result.OnlyInBase)
	assert.Empty(t, result.OnlyInOther)
	assert.Len(t, result.InBoth, 2)
}

func TestDiffKeys_EmptyInputs(t *testing.T) {
	result := DiffKeys(nil, nil, DefaultDiffKeysOptions())
	assert.Empty(t, result.OnlyInBase)
	assert.Empty(t, result.OnlyInOther)
	assert.Empty(t, result.InBoth)
}

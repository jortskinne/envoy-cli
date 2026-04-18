package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeReorderEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "C", Value: "3"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "D", Value: "4"},
	}
}

func TestReorder_MatchesKeyOrder(t *testing.T) {
	entries := makeReorderEntries()
	result := Reorder(entries, []string{"A", "B", "C"}, DefaultReorderOptions())
	assert.Equal(t, "A", result[0].Key)
	assert.Equal(t, "B", result[1].Key)
	assert.Equal(t, "C", result[2].Key)
	// D appended at end
	assert.Equal(t, "D", result[3].Key)
}

func TestReorder_SkipsMissingKeysInReference(t *testing.T) {
	entries := makeReorderEntries()
	// "Z" is not in entries, should be silently skipped
	result := Reorder(entries, []string{"Z", "A", "C"}, DefaultReorderOptions())
	assert.Equal(t, "A", result[0].Key)
	assert.Equal(t, "C", result[1].Key)
}

func TestReorder_DropsExtraWhenNotAppend(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.AppendExtra = false
	result := Reorder(entries, []string{"A", "C"}, opts)
	assert.Len(t, result, 2)
	assert.Equal(t, "A", result[0].Key)
	assert.Equal(t, "C", result[1].Key)
}

func TestReorder_EmptyKeys(t *testing.T) {
	entries := makeReorderEntries()
	result := Reorder(entries, []string{}, DefaultReorderOptions())
	// All entries appended as extras
	assert.Len(t, result, 4)
}

func TestReorder_EmptyEntries(t *testing.T) {
	result := Reorder([]EnvEntry{}, []string{"A", "B"}, DefaultReorderOptions())
	assert.Empty(t, result)
}

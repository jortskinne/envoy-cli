package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeShrinkEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "  # comment", Value: ""},
		{Key: "DB_HOST", Value: "  localhost  "},
		{Key: "EMPTY_KEY", Value: ""},
		{Key: "APP_NAME", Value: "overridden"},
		{Key: "SECRET", Value: "abc123"},
	}
}

func TestShrink_RemovesComments(t *testing.T) {
	entries := makeShrinkEntries()
	opts := DefaultShrinkOptions()
	opts.RemoveEmpty = false
	opts.DedupeKeys = false
	opts.TrimValues = false

	result, removed := Shrink(entries, opts)

	assert.Equal(t, 1, removed)
	for _, e := range result {
		assert.NotContains(t, e.Key, "#")
	}
}

func TestShrink_TrimsValues(t *testing.T) {
	entries := makeShrinkEntries()
	opts := DefaultShrinkOptions()
	opts.RemoveComments = false
	opts.RemoveEmpty = false
	opts.DedupeKeys = false

	result, _ := Shrink(entries, opts)

	for _, e := range result {
		if e.Key == "DB_HOST" {
			assert.Equal(t, "localhost", e.Value)
		}
	}
}

func TestShrink_RemovesEmpty(t *testing.T) {
	entries := makeShrinkEntries()
	opts := DefaultShrinkOptions()
	opts.RemoveComments = false
	opts.DedupeKeys = false

	result, _ := Shrink(entries, opts)

	for _, e := range result {
		assert.NotEmpty(t, e.Value, "entry %q should not have empty value", e.Key)
	}
}

func TestShrink_DeduplicatesKeys(t *testing.T) {
	entries := makeShrinkEntries()
	opts := DefaultShrinkOptions()
	opts.RemoveComments = false
	opts.RemoveEmpty = false
	opts.TrimValues = false

	result, removed := Shrink(entries, opts)

	keys := make(map[string]int)
	for _, e := range result {
		keys[e.Key]++
	}
	assert.Equal(t, 1, keys["APP_NAME"], "APP_NAME should appear exactly once")
	assert.Equal(t, 1, removed)

	// Last value wins.
	for _, e := range result {
		if e.Key == "APP_NAME" {
			assert.Equal(t, "overridden", e.Value)
		}
	}
}

func TestShrink_DefaultOptions(t *testing.T) {
	entries := makeShrinkEntries()
	opts := DefaultShrinkOptions()

	result, removed := Shrink(entries, opts)

	// comment removed (1) + APP_NAME deduped (1) = 2
	assert.Equal(t, 2, removed)
	assert.Len(t, result, 4)
}

func TestShrink_EmptyInput(t *testing.T) {
	result, removed := Shrink([]EnvEntry{}, DefaultShrinkOptions())
	assert.Empty(t, result)
	assert.Equal(t, 0, removed)
}

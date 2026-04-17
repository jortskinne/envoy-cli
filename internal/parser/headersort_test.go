package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeGroupEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc"},
		{Key: "NOPREFIX", Value: "standalone"},
	}
}

func TestGroupByPrefix_BasicGrouping(t *testing.T) {
	entries := makeGroupEntries()
	opts := DefaultGroupOptions()
	result := GroupByPrefix(entries, opts)

	assert.Equal(t, 4, len(result.Sections))
	assert.Equal(t, "DB", result.Sections[0].Header)
	assert.Equal(t, 2, len(result.Sections[0].Entries))
	assert.Equal(t, "APP", result.Sections[1].Header)
	assert.Equal(t, 2, len(result.Sections[1].Entries))
}

func TestGroupByPrefix_NoSeparatorGoesToOther(t *testing.T) {
	entries := makeGroupEntries()
	opts := DefaultGroupOptions()
	result := GroupByPrefix(entries, opts)

	last := result.Sections[len(result.Sections)-1]
	assert.Equal(t, "OTHER", last.Header)
	assert.Equal(t, 1, len(last.Entries))
	assert.Equal(t, "NOPREFIX", last.Entries[0].Key)
}

func TestFlattenGrouped_WithHeaders(t *testing.T) {
	entries := makeGroupEntries()
	opts := DefaultGroupOptions()
	grouped := GroupByPrefix(entries, opts)
	flat := FlattenGrouped(grouped, true)

	// 4 sections => 4 comment headers + 6 entries
	assert.Equal(t, 10, len(flat))
	assert.True(t, flat[0].IsComment)
	assert.Equal(t, "# --- DB ---", flat[0].Comment)
}

func TestFlattenGrouped_WithoutHeaders(t *testing.T) {
	entries := makeGroupEntries()
	opts := DefaultGroupOptions()
	grouped := GroupByPrefix(entries, opts)
	flat := FlattenGrouped(grouped, false)

	assert.Equal(t, 6, len(flat))
	for _, e := range flat {
		assert.False(t, e.IsComment)
	}
}

func TestGroupByPrefix_EmptyInput(t *testing.T) {
	result := GroupByPrefix([]EnvEntry{}, DefaultGroupOptions())
	assert.Equal(t, 0, len(result.Sections))
}

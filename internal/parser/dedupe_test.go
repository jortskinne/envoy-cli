package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeDedupeEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "first"},
		{Key: "DEBUG", Value: "true"},
		{Key: "APP_NAME", Value: "second"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "false"},
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	entries := makeDedupeEntries()
	result := Dedupe(entries, DefaultDedupeOptions())

	assert.Len(t, result.Entries, 3)
	assert.Equal(t, "first", result.Entries[0].Value)  // APP_NAME kept first
	assert.Equal(t, "true", result.Entries[1].Value)   // DEBUG kept first
	assert.Len(t, result.Removed, 2)
}

func TestDedupe_KeepLast(t *testing.T) {
	entries := makeDedupeEntries()
	opts := DedupeOptions{KeepFirst: false}
	result := Dedupe(entries, opts)

	assert.Len(t, result.Entries, 3)
	assert.Equal(t, "second", result.Entries[0].Value) // APP_NAME replaced
	assert.Equal(t, "false", result.Entries[1].Value)  // DEBUG replaced
	assert.Len(t, result.Removed, 2)
}

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := []EnvEntry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	result := Dedupe(entries, DefaultDedupeOptions())

	assert.Len(t, result.Entries, 2)
	assert.Empty(t, result.Removed)
}

func TestDedupe_EmptyInput(t *testing.T) {
	result := Dedupe([]EnvEntry{}, DefaultDedupeOptions())
	assert.Empty(t, result.Entries)
	assert.Empty(t, result.Removed)
}

func TestDedupe_DoesNotMutateOriginal(t *testing.T) {
	entries := makeDedupeEntries()
	copy := append([]EnvEntry{}, entries...)
	Dedupe(entries, DefaultDedupeOptions())
	assert.Equal(t, copy, entries)
}

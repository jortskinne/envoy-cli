package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makePrefixEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "APP_NAME", Value: "envoy"},
	}
}

func TestAddPrefix_Basic(t *testing.T) {
	entries := makePrefixEntries()
	result := AddPrefix(entries, "APP_", DefaultPrefixOptions())
	assert.Equal(t, "APP_HOST", result[0].Key)
	assert.Equal(t, "APP_PORT", result[1].Key)
	// already has prefix — should not double-prefix
	assert.Equal(t, "APP_NAME", result[2].Key)
}

func TestAddPrefix_DryRun(t *testing.T) {
	entries := makePrefixEntries()
	opts := DefaultPrefixOptions()
	opts.DryRun = true
	result := AddPrefix(entries, "PRE_", opts)
	// keys must remain unchanged in dry-run
	for i, e := range result {
		assert.Equal(t, entries[i].Key, e.Key)
	}
}

func TestRemovePrefix_Basic(t *testing.T) {
	entries := []EnvEntry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "NAME", Value: "envoy"},
	}
	result := RemovePrefix(entries, "APP_", DefaultPrefixOptions())
	assert.Len(t, result, 3)
	assert.Equal(t, "HOST", result[0].Key)
	assert.Equal(t, "PORT", result[1].Key)
	assert.Equal(t, "NAME", result[2].Key)
}

func TestRemovePrefix_SkipsEmptyKey(t *testing.T) {
	entries := []EnvEntry{
		{Key: "APP_", Value: "bad"},
		{Key: "APP_HOST", Value: "localhost"},
	}
	result := RemovePrefix(entries, "APP_", DefaultPrefixOptions())
	assert.Len(t, result, 1)
	assert.Equal(t, "HOST", result[0].Key)
}

func TestRemovePrefix_DryRun(t *testing.T) {
	entries := []EnvEntry{
		{Key: "APP_HOST", Value: "localhost"},
	}
	opts := DefaultPrefixOptions()
	opts.DryRun = true
	result := RemovePrefix(entries, "APP_", opts)
	assert.Equal(t, "APP_HOST", result[0].Key)
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeOmitEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "DEBUG", Value: "true"},
		{Key: "SECRET_TOKEN", Value: "tok"},
	}
}

func TestOmit_ByExplicitKeys(t *testing.T) {
	entries := makeOmitEntries()
	opts := DefaultOmitOptions()
	opts.Keys = []string{"DEBUG", "APP_NAME"}

	result, err := Omit(entries, opts)
	require.NoError(t, err)

	keys := make([]string, 0, len(result))
	for _, e := range result {
		keys = append(keys, e.Key)
	}
	assert.NotContains(t, keys, "DEBUG")
	assert.NotContains(t, keys, "APP_NAME")
	assert.Contains(t, keys, "DB_PASSWORD")
	assert.Contains(t, keys, "API_KEY")
}

func TestOmit_ByPrefix(t *testing.T) {
	entries := makeOmitEntries()
	opts := DefaultOmitOptions()
	opts.Prefix = "DB_"

	result, err := Omit(entries, opts)
	require.NoError(t, err)

	for _, e := range result {
		assert.False(t, len(e.Key) >= 3 && e.Key[:3] == "DB_", "key %q should be omitted", e.Key)
	}
	assert.Len(t, result, 4)
}

func TestOmit_SensitiveKeys(t *testing.T) {
	entries := makeOmitEntries()
	opts := DefaultOmitOptions()
	opts.OmitSensitive = true

	result, err := Omit(entries, opts)
	require.NoError(t, err)

	for _, e := range result {
		assert.False(t, IsSensitive(e.Key), "sensitive key %q should be omitted", e.Key)
	}
}

func TestOmit_NoOptions_ReturnsAll(t *testing.T) {
	entries := makeOmitEntries()
	opts := DefaultOmitOptions()

	result, err := Omit(entries, opts)
	require.NoError(t, err)
	assert.Len(t, result, len(entries))
}

func TestOmit_CombinedKeysAndPrefix(t *testing.T) {
	entries := makeOmitEntries()
	opts := DefaultOmitOptions()
	opts.Keys = []string{"DEBUG"}
	opts.Prefix = "API_"

	result, err := Omit(entries, opts)
	require.NoError(t, err)

	keys := make([]string, 0)
	for _, e := range result {
		keys = append(keys, e.Key)
	}
	assert.NotContains(t, keys, "DEBUG")
	assert.NotContains(t, keys, "API_KEY")
	assert.Contains(t, keys, "APP_NAME")
	assert.Contains(t, keys, "DB_HOST")
}

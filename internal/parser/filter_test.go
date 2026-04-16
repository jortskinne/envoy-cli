package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeFilterEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "abc123"},
	}
}

func TestFilter_ByKeys(t *testing.T) {
	entries := makeFilterEntries()
	opts := DefaultFilterOptions()
	opts.Keys = []string{"APP_NAME", "DB_HOST"}
	out := Filter(entries, opts)
	assert.Len(t, out, 2)
	assert.Equal(t, "APP_NAME", out[0].Key)
	assert.Equal(t, "DB_HOST", out[1].Key)
}

func TestFilter_ByPrefix(t *testing.T) {
	entries := makeFilterEntries()
	opts := DefaultFilterOptions()
	opts.Prefix = "APP_"
	out := Filter(entries, opts)
	assert.Len(t, out, 2)
	for _, e := range out {
		assert.True(t, len(e.Key) > 4 && e.Key[:4] == "APP_")
	}
}

func TestFilter_Exclude(t *testing.T) {
	entries := makeFilterEntries()
	opts := DefaultFilterOptions()
	opts.Exclude = []string{"DB_PASSWORD", "API_KEY"}
	out := Filter(entries, opts)
	for _, e := range out {
		assert.NotEqual(t, "DB_PASSWORD", e.Key)
		assert.NotEqual(t, "API_KEY", e.Key)
	}
	assert.Len(t, out, 3)
}

func TestFilter_SensitiveOnly(t *testing.T) {
	entries := makeFilterEntries()
	opts := DefaultFilterOptions()
	opts.Sensitive = true
	out := Filter(entries, opts)
	for _, e := range out {
		assert.True(t, IsSensitive(e.Key, DefaultMaskOptions()))
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	entries := makeFilterEntries()
	out := Filter(entries, DefaultFilterOptions())
	assert.Len(t, out, len(entries))
}

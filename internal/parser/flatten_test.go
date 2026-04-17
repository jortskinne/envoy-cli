package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeFlattenEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "DB__HOST", Value: "localhost"},
		{Key: "DB__PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "CACHE__TTL", Value: "300"},
	}
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	entries := makeFlattenEntries()
	opts := DefaultFlattenOptions()
	out := Flatten(entries, opts)

	assert.Equal(t, "DB.HOST", out[0].Key)
	assert.Equal(t, "DB.PORT", out[1].Key)
	assert.Equal(t, "APP_NAME", out[2].Key) // unchanged
	assert.Equal(t, "CACHE.TTL", out[3].Key)
}

func TestFlatten_PreservesValues(t *testing.T) {
	entries := makeFlattenEntries()
	out := Flatten(entries, DefaultFlattenOptions())
	assert.Equal(t, "localhost", out[0].Value)
	assert.Equal(t, "5432", out[1].Value)
}

func TestFlatten_Lowercase(t *testing.T) {
	entries := []EnvEntry{
		{Key: "DB__HOST", Value: "localhost"},
	}
	opts := DefaultFlattenOptions()
	opts.Lowercase = true
	out := Flatten(entries, opts)
	assert.Equal(t, "db.host", out[0].Key)
}

func TestFlatten_PrefixFilter(t *testing.T) {
	entries := makeFlattenEntries()
	opts := DefaultFlattenOptions()
	opts.Prefix = "DB"
	out := Flatten(entries, opts)

	// DB__ keys flattened, others unchanged
	assert.Equal(t, "DB.HOST", out[0].Key)
	assert.Equal(t, "DB.PORT", out[1].Key)
	assert.Equal(t, "APP_NAME", out[2].Key)
	assert.Equal(t, "CACHE__TTL", out[3].Key) // not touched
}

func TestFlatten_NoSeparatorPassthrough(t *testing.T) {
	entries := []EnvEntry{
		{Key: "SIMPLE", Value: "val"},
	}
	out := Flatten(entries, DefaultFlattenOptions())
	assert.Equal(t, "SIMPLE", out[0].Key)
}

func TestFlatten_EmptyInput(t *testing.T) {
	out := Flatten([]EnvEntry{}, DefaultFlattenOptions())
	assert.Empty(t, out)
}

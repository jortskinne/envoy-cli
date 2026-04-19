package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeMapEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DEBUG", Value: "true"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestToMap_Basic(t *testing.T) {
	entries := makeMapEntries()
	m := ToMap(entries)
	assert.Equal(t, "envoy", m["APP_NAME"])
	assert.Equal(t, "true", m["DEBUG"])
	assert.Equal(t, "abc123", m["SECRET_KEY"])
	assert.Len(t, m, 3)
}

func TestFromMap_SortedOrder(t *testing.T) {
	m := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	entries := FromMap(m)
	assert.Equal(t, "A_KEY", entries[0].Key)
	assert.Equal(t, "M_KEY", entries[1].Key)
	assert.Equal(t, "Z_KEY", entries[2].Key)
}

func TestMergeMap_NoOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original", "B": "keep"}
	src := map[string]string{"A": "new", "C": "added"}
	out := MergeMap(dst, src, false)
	assert.Equal(t, "original", out["A"])
	assert.Equal(t, "keep", out["B"])
	assert.Equal(t, "added", out["C"])
}

func TestMergeMap_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "replaced"}
	out := MergeMap(dst, src, true)
	assert.Equal(t, "replaced", out["A"])
}

func TestFilterMap_BySensitiveKey(t *testing.T) {
	m := map[string]string{
		"SECRET_KEY": "s3cr3t",
		"APP_NAME":   "envoy",
		"DB_PASSWORD": "pass",
	}
	out := FilterMap(m, func(k, _ string) bool {
		return IsSensitive(k)
	})
	assert.Contains(t, out, "SECRET_KEY")
	assert.Contains(t, out, "DB_PASSWORD")
	assert.NotContains(t, out, "APP_NAME")
}

func TestToMap_EmptyEntries(t *testing.T) {
	m := ToMap([]Entry{})
	assert.Empty(t, m)
}

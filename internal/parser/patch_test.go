package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePatchEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
}

func TestPatch_SetExisting(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "set", Key: "DB_HOST", Value: "db.prod.internal"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "db.prod.internal", findValue(out, "DB_HOST"))
}

func TestPatch_SetNew(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "set", Key: "REDIS_URL", Value: "redis://localhost"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "redis://localhost", findValue(out, "REDIS_URL"))
	assert.Len(t, out, 4)
}

func TestPatch_Delete(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "delete", Key: "DB_PORT"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "", findValue(out, "DB_PORT"))
}

func TestPatch_DeleteMissing_Error(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "delete", Key: "NONEXISTENT"}}
	_, err := Patch(entries, ops, DefaultPatchOptions())
	assert.Error(t, err)
}

func TestPatch_DeleteMissing_Ignore(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "delete", Key: "NONEXISTENT"}}
	out, err := Patch(entries, ops, PatchOptions{IgnoreMissing: true})
	require.NoError(t, err)
	assert.Len(t, out, 3)
}

func TestPatch_Rename(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "rename", Key: "APP_ENV", NewKey: "APP_ENVIRONMENT"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "production", findValue(out, "APP_ENVIRONMENT"))
	assert.Equal(t, "", findValue(out, "APP_ENV"))
}

func TestPatch_UnknownAction(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOperation{{Action: "upsert", Key: "X"}}
	_, err := Patch(entries, ops, DefaultPatchOptions())
	assert.Error(t, err)
}

func findValue(entries []EnvEntry, key string) string {
	for _, e := range entries {
		if e.Key == key {
			return e.Value
		}
	}
	return ""
}

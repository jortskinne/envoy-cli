package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCrossCheckBase() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func makeCrossCheckOther() []EnvEntry {
	return []EnvEntry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "prod-db"},
		{Key: "EXTRA_KEY", Value: "extra"},
	}
}

func TestCrossCheck_DetectsMissingInOther(t *testing.T) {
	base := makeCrossCheckBase()
	other := makeCrossCheckOther()
	opts := DefaultCrossCheckOptions()

	results := CrossCheck(base, other, opts)

	statuses := map[string]string{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	assert.Equal(t, "ok", statuses["APP_NAME"])
	assert.Equal(t, "ok", statuses["DB_HOST"])
	assert.Equal(t, "missing_in_other", statuses["SECRET_KEY"])
}

func TestCrossCheck_DetectsMissingInBase(t *testing.T) {
	base := makeCrossCheckBase()
	other := makeCrossCheckOther()
	opts := DefaultCrossCheckOptions()
	opts.RequireAllOther = true

	results := CrossCheck(base, other, opts)

	statuses := map[string]string{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	assert.Equal(t, "missing_in_base", statuses["EXTRA_KEY"])
}

func TestCrossCheck_IgnoreCase(t *testing.T) {
	base := []EnvEntry{{Key: "app_name", Value: "myapp"}}
	other := []EnvEntry{{Key: "APP_NAME", Value: "myapp"}}
	opts := DefaultCrossCheckOptions()
	opts.IgnoreCase = true

	results := CrossCheck(base, other, opts)

	assert.Len(t, results, 1)
	assert.Equal(t, "ok", results[0].Status)
}

func TestCrossCheck_OKEntryHasBothValues(t *testing.T) {
	base := []EnvEntry{{Key: "DB_HOST", Value: "localhost"}}
	other := []EnvEntry{{Key: "DB_HOST", Value: "prod-db"}}
	opts := DefaultCrossCheckOptions()

	results := CrossCheck(base, other, opts)

	assert.Len(t, results, 1)
	assert.Equal(t, "ok", results[0].Status)
	assert.Equal(t, "localhost", results[0].BaseValue)
	assert.Equal(t, "prod-db", results[0].OtherValue)
}

func TestCrossCheck_EmptyBase(t *testing.T) {
	results := CrossCheck([]EnvEntry{}, makeCrossCheckOther(), DefaultCrossCheckOptions())
	assert.Empty(t, results)
}

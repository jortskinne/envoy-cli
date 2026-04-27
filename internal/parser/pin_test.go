package parser

import (
	"testing"
)

func makePinEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "SECRET", Value: "abc", Comment: "# sensitive"},
		{Key: "DB_URL", Value: "postgres://", Comment: "# @pinned"},
	}
}

func TestPin_ExplicitKeys(t *testing.T) {
	entries := makePinEntries()
	opts := DefaultPinOptions()
	opts.Keys = []string{"HOST", "PORT"}

	out, results, err := Pin(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Pinned {
			t.Errorf("expected %s to be pinned", r.Key)
		}
	}
	for _, e := range out {
		if (e.Key == "HOST" || e.Key == "PORT") && !IsPinned(e) {
			t.Errorf("expected %s comment to contain @pinned", e.Key)
		}
	}
}

func TestPin_SkipsAlreadyPinned(t *testing.T) {
	entries := makePinEntries()
	opts := DefaultPinOptions()

	_, results, _ := Pin(entries, opts)

	var dbResult *PinResult
	for i := range results {
		if results[i].Key == "DB_URL" {
			dbResult = &results[i]
		}
	}
	if dbResult == nil {
		t.Fatal("expected DB_URL in results")
	}
	if !dbResult.Skipped {
		t.Errorf("expected DB_URL to be skipped (already pinned)")
	}
}

func TestPin_OverwriteFlag(t *testing.T) {
	entries := makePinEntries()
	opts := DefaultPinOptions()
	opts.Overwrite = true
	opts.Keys = []string{"DB_URL"}

	_, results, _ := Pin(entries, opts)

	var dbResult *PinResult
	for i := range results {
		if results[i].Key == "DB_URL" {
			dbResult = &results[i]
		}
	}
	if dbResult == nil || !dbResult.Pinned {
		t.Errorf("expected DB_URL to be re-pinned with Overwrite=true")
	}
}

func TestPin_DryRun(t *testing.T) {
	entries := makePinEntries()
	opts := DefaultPinOptions()
	opts.DryRun = true
	opts.Keys = []string{"HOST"}

	out, results, _ := Pin(entries, opts)

	if len(results) == 0 || !results[0].Pinned {
		t.Errorf("expected result to show HOST as pinned")
	}
	for _, e := range out {
		if e.Key == "HOST" && IsPinned(e) {
			t.Errorf("dry-run should not modify entry comment")
		}
	}
}

func TestIsPinned_DetectsTag(t *testing.T) {
	e := EnvEntry{Key: "X", Value: "1", Comment: "# @pinned"}
	if !IsPinned(e) {
		t.Error("expected IsPinned to return true")
	}
}

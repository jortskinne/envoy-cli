package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func sampleWatchEvent() WatchEvent {
	return WatchEvent{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		File:      ".env.production",
		Count:     3,
		Entries: []EnvEntry{
			{Key: "FOO", Value: "bar"},
		},
	}
}

func TestWriteWatchEvent_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	e := sampleWatchEvent()
	if err := WriteWatchEvent(&buf, e, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, ".env.production") {
		t.Errorf("expected file name in output, got: %s", out)
	}
	if !strings.Contains(out, "3 entries") {
		t.Errorf("expected entry count in output, got: %s", out)
	}
	if !strings.Contains(out, "12:00:00") {
		t.Errorf("expected timestamp in output, got: %s", out)
	}
}

func TestWriteWatchEvent_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	e := sampleWatchEvent()
	if err := WriteWatchEvent(&buf, e, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded WatchEvent
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if decoded.File != ".env.production" {
		t.Errorf("expected file .env.production, got %s", decoded.File)
	}
	if decoded.Count != 3 {
		t.Errorf("expected count 3, got %d", decoded.Count)
	}
}

func TestWriteWatchEvent_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	e := sampleWatchEvent()
	if err := WriteWatchEvent(&buf, e, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for default format")
	}
}

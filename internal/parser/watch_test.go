package parser

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func writeWatchTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeWatchTempEnv(t, "FOO=bar\n")

	var mu sync.Mutex
	var received []EnvEntry

	opts := DefaultWatchOptions()
	opts.Interval = 50 * time.Millisecond
	opts.OnChange = func(entries []EnvEntry) {
		mu.Lock()
		defer mu.Unlock()
		received = entries
	}

	done := make(chan struct{})
	go func() {
		_ = Watch(path, opts, done)
	}()

	time.Sleep(80 * time.Millisecond)
	_ = os.WriteFile(path, []byte("FOO=changed\nBAR=new\n"), 0644)
	time.Sleep(150 * time.Millisecond)
	close(done)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 entries after change, got %d", len(received))
	}
}

func TestWatch_NoCallbackWhenUnchanged(t *testing.T) {
	path := writeWatchTempEnv(t, "FOO=bar\n")

	callCount := 0
	opts := DefaultWatchOptions()
	opts.Interval = 40 * time.Millisecond
	opts.OnChange = func(_ []EnvEntry) { callCount++ }

	done := make(chan struct{})
	go func() { _ = Watch(path, opts, done) }()
	time.Sleep(180 * time.Millisecond)
	close(done)

	if callCount != 0 {
		t.Errorf("expected 0 callbacks for unchanged file, got %d", callCount)
	}
}

func TestWatch_MissingFile_ReturnsError(t *testing.T) {
	done := make(chan struct{})
	close(done)
	err := Watch("/nonexistent/.env", DefaultWatchOptions(), done)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWatch_ErrorCallbackOnParseFailure(t *testing.T) {
	path := writeWatchTempEnv(t, "FOO=bar\n")

	var gotErr error
	opts := DefaultWatchOptions()
	opts.Interval = 40 * time.Millisecond
	opts.OnError = func(err error) { gotErr = err }

	done := make(chan struct{})
	go func() { _ = Watch(path, opts, done) }()
	time.Sleep(60 * time.Millisecond)
	_ = os.Remove(path)
	time.Sleep(100 * time.Millisecond)
	close(done)

	if gotErr == nil {
		t.Error("expected OnError to be called after file removal")
	}
}

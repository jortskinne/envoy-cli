package parser

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the file watcher behaviour.
type WatchOptions struct {
	Interval  time.Duration
	OnChange  func(entries []EnvEntry)
	OnError   func(err error)
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval: 2 * time.Second,
		OnChange: func(_ []EnvEntry) {},
		OnError:  func(_ error) {},
	}
}

// Watch monitors a .env file for changes and invokes callbacks.
// It blocks until the done channel is closed.
func Watch(path string, opts WatchOptions, done <-chan struct{}) error {
	lastHash, err := fileHash(path)
	if err != nil {
		return fmt.Errorf("watch: initial read failed: %w", err)
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			h, err := fileHash(path)
			if err != nil {
				opts.OnError(err)
				continue
			}
			if h != lastHash {
				lastHash = h
				entries, err := ParseFile(path)
				if err != nil {
					opts.OnError(err)
					continue
				}
				opts.OnChange(entries)
			}
		}
	}
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

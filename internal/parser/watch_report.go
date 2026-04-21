package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WatchEvent captures a single detected change event.
type WatchEvent struct {
	Timestamp time.Time  `json:"timestamp"`
	File      string     `json:"file"`
	Count     int        `json:"entry_count"`
	Entries   []EnvEntry `json:"entries,omitempty"`
}

// WriteWatchEvent writes a human-readable or JSON event line to w.
func WriteWatchEvent(w io.Writer, event WatchEvent, format string) error {
	switch format {
	case "json":
		return writeWatchEventJSON(w, event)
	default:
		return writeWatchEventText(w, event)
	}
}

func writeWatchEventText(w io.Writer, e WatchEvent) error {
	_, err := fmt.Fprintf(w,
		"[%s] %s changed — %d entries loaded\n",
		e.Timestamp.Format("15:04:05"),
		e.File,
		e.Count,
	)
	return err
}

func writeWatchEventJSON(w io.Writer, e WatchEvent) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

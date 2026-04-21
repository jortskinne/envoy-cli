package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/parser"
)

var watchCmd = &cobra.Command{
	Use:   "watch <file>",
	Short: "Watch a .env file and print changes as they occur",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	watchCmd.Flags().DurationP("interval", "i", 2*time.Second, "Poll interval (e.g. 500ms, 2s)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := args[0]
	format, _ := cmd.Flags().GetString("format")
	interval, _ := cmd.Flags().GetDuration("interval")

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file not found: %s", path)
	}

	fmt.Fprintf(os.Stderr, "Watching %s (interval: %s) — press Ctrl+C to stop\n", path, interval)

	opts := parser.DefaultWatchOptions()
	opts.Interval = interval
	opts.OnChange = func(entries []parser.EnvEntry) {
		event := parser.WatchEvent{
			Timestamp: time.Now(),
			File:      path,
			Count:     len(entries),
			Entries:   entries,
		}
		_ = parser.WriteWatchEvent(os.Stdout, event, format)
	}
	opts.OnError = func(err error) {
		fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
	}

	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(done)
	}()

	return parser.Watch(path, opts, done)
}

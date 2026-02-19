package configs

import (
	"log/slog"
	"os"
)

// NewLogger() returns a new logger instance for use
func NewLogger(isDev bool) *slog.Logger {
	var handler slog.Handler
	if isDev {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	return slog.New(handler)
}

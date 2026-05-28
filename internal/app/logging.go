package app

import (
	"log/slog"
	"os"
)

// SetupLogging initializes the application's structured JSON logger
// and sets it as the default slog logger. This should be called at
// the very start of main() before any other initialization occurs.
func SetupLogging() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

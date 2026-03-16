package logs

import (
	"log/slog"
)

func NoopHandler() slog.Handler {
	return slog.DiscardHandler
}

func NoopLogger() *slog.Logger {
	return slog.New(NoopHandler())
}

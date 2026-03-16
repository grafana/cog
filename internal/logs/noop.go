package logs

import (
	"io"
	"log/slog"
)

func NoopHandler() slog.Handler {
	return slog.NewTextHandler(io.Discard(), nil)
}

func NoopLogger() *slog.Logger {
	return slog.New(NoopHandler())
}

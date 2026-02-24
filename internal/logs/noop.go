package logs

import (
	"log/slog"
	"os"
)

func NoopHandler() slog.Handler {
	return slog.NewTextHandler(os.NewFile(0, os.DevNull), nil)
}

func NoopLogger() *slog.Logger {
	return slog.New(NoopHandler())
}

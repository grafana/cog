package logs

import (
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type slogWrapper struct {
	inner *slog.Logger
}

func HCLLoggerFromSlog(logger *slog.Logger) hclog.Logger {
	return &slogWrapper{
		inner: logger,
	}
}

func hclLevelToSlog(level hclog.Level) slog.Level {
	switch level {
	case hclog.Trace:
		return slog.LevelDebug - 1
	case hclog.Debug:
		return slog.LevelDebug
	case hclog.Info:
		return slog.LevelInfo
	case hclog.Warn:
		return slog.LevelWarn
	case hclog.Error:
		return slog.LevelError
	case hclog.Off:
		return slog.LevelError + 4
	default:
		return slog.LevelInfo
	}
}

func (s slogWrapper) Log(level hclog.Level, msg string, args ...any) {
	levelPrefixes := []string{"DEBUG ", "INFO ", "WARN ", "ERROR "}
	for _, prefix := range levelPrefixes {
		msg = strings.TrimPrefix(msg, prefix)
	}

	s.inner.Log(nil, hclLevelToSlog(level), msg, args...)
}

func (s slogWrapper) Trace(msg string, args ...any) {
	s.Log(hclog.Trace, msg, args...)
}

func (s slogWrapper) Debug(msg string, args ...any) {
	s.Log(hclog.Debug, msg, args...)
}

func (s slogWrapper) Info(msg string, args ...any) {
	s.Log(hclog.Info, msg, args...)
}

func (s slogWrapper) Warn(msg string, args ...any) {
	s.Log(hclog.Warn, msg, args...)
}

func (s slogWrapper) Error(msg string, args ...any) {
	s.Log(hclog.Error, msg, args...)
}

func (s slogWrapper) IsTrace() bool {
	return s.inner.Enabled(nil, slog.LevelDebug-1)
}

func (s slogWrapper) IsDebug() bool {
	return s.inner.Enabled(nil, slog.LevelDebug)
}

func (s slogWrapper) IsInfo() bool {
	return s.inner.Enabled(nil, slog.LevelInfo)
}

func (s slogWrapper) IsWarn() bool {
	return s.inner.Enabled(nil, slog.LevelWarn)
}

func (s slogWrapper) IsError() bool {
	return s.inner.Enabled(nil, slog.LevelError)
}

func (s slogWrapper) ImpliedArgs() []any {
	panic("slogWrapper.ImpliedArgs(): not implemented")
}

func (s slogWrapper) With(args ...any) hclog.Logger {
	panic("slogWrapper.With(): not implemented")
}

func (s slogWrapper) Name() string {
	return ""
}

func (s slogWrapper) Named(_ string) hclog.Logger {
	return s
}

func (s slogWrapper) ResetNamed(_ string) hclog.Logger {
	return s
}

func (s slogWrapper) SetLevel(level hclog.Level) {
}

func (s slogWrapper) GetLevel() hclog.Level {
	return hclog.Trace
}

func (s slogWrapper) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	panic("slogWrapper.StandardLogger(): not implemented")
}

func (s slogWrapper) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	panic("slogWrapper.StandardWriter(): not implemented")
}

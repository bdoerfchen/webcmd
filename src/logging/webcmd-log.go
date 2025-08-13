package logging

import (
	"context"
	"log/slog"
	"os"
)

type loggerKey struct{}

var contextLoggerKey = loggerKey{}

func New(level slog.Level, color bool) *slog.Logger {
	return slog.New(newLogHandler(os.Stdout, &slog.HandlerOptions{Level: level}, color))
}

func FromContext(ctx context.Context) *slog.Logger {
	// Get logger as value from context by specific key
	logger, ok := ctx.Value(contextLoggerKey).(*slog.Logger)
	if ok {
		return logger
	}

	// If no logger in context, return debug logger
	return New(slog.LevelDebug, false).With(slog.String("logger", "tmp"))
}

func AddToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger)
}

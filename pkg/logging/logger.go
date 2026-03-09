package logging

import (
	"context"
	"log/slog"
	"os"
)

// NewLogger creates a structured JSON logger.
func NewLogger(serviceName string, level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler).With(
		slog.String("service", serviceName),
	)
}

type loggerKey struct{}

// WithLogger stores logger in context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext retrieves logger from context, returns default if not found.
func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

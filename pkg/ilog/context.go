package ilog

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}

func ContextWithAttrs(ctx context.Context, attrs ...any) context.Context {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		log = log.With(attrs...)
		return context.WithValue(ctx, ctxKey{}, log)
	}

	log := slog.Default().With(attrs...)
	return context.WithValue(ctx, ctxKey{}, log)
}

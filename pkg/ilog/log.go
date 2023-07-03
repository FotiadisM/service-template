package ilog

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

type handler struct {
	slog.Handler

	ctxKeys []any
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	for _, k := range h.ctxKeys {
		if v := ctx.Value(k); v != nil {
			// TODO(FotiadisM): improve functionality
			r.Add(k, v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

type Option func(*handler)

func WithContextKeys(keys ...any) Option {
	return func(h *handler) {
		if len(keys) > 0 {
			h.ctxKeys = keys
		}
	}
}

func NewHandler(opts ...Option) slog.Handler {
	h := &handler{
		Handler: slog.NewJSONHandler(os.Stdout, nil),
		ctxKeys: []any{},
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

type ctxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

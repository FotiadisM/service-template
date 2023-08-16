package ilog

import (
	"context"
	"log/slog"
	"os"
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

func NewHandler(ho *slog.HandlerOptions, opts ...Option) slog.Handler {
	h := &handler{
		Handler: slog.NewJSONHandler(os.Stdout, ho),
		ctxKeys: []any{},
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

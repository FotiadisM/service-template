package ilog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
)

func try(callback func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unexpected error: %+v", r)
			}
		}
	}()

	err = callback()

	return
}

type FanoutHandler struct {
	handlers []slog.Handler
}

var _ slog.Handler = (*FanoutHandler)(nil)

func NewFanoutHandler(handlers ...slog.Handler) slog.Handler {
	return &FanoutHandler{
		handlers: handlers,
	}
}

func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// If errs is empty, or contains only nil errors, this returns nil
	return errors.Join(errs...)
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		handlers[i] = h.WithAttrs(slices.Clone(attrs))
	}

	return NewFanoutHandler(handlers...)
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	handlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		handlers[i] = h.WithGroup(name)
	}

	return NewFanoutHandler(handlers...)
}

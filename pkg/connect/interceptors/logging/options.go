package logging

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

type FilterFunc func(ctx context.Context, spec connect.Spec) bool

type CodeToLevelFunc func(code connect.Code) slog.Level

type ErrorDetailsAttrFunc func(details []*connect.ErrorDetail) []slog.Attr

type options struct {
	filterFunc           FilterFunc
	codeToLevelFunc      CodeToLevelFunc
	errorDetailsAttrFunc ErrorDetailsAttrFunc
}

func defaultOptions() *options {
	return &options{
		filterFunc:           func(_ context.Context, _ connect.Spec) bool { return true },
		codeToLevelFunc:      DefaultCodeToLevelFunc,
		errorDetailsAttrFunc: DefaultErrorDetailsAttrFunc,
	}
}

type Option func(*options)

func WithFilterFunc(f FilterFunc) Option {
	return func(o *options) {
		o.filterFunc = f
	}
}

func WithCodeToLevelFunc(f CodeToLevelFunc) Option {
	return func(o *options) {
		o.codeToLevelFunc = f
	}
}

func WithErrorDetailsAttrFunc(f ErrorDetailsAttrFunc) Option {
	return func(o *options) {
		o.errorDetailsAttrFunc = f
	}
}

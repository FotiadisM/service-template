package logging

import (
	"log/slog"

	"google.golang.org/grpc/codes"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
)

type CodeToLevelFunc func(code codes.Code) slog.Level

type GRPCStatusDetailsToLogAttrsFunc func(details []any) []slog.Attr

type options struct {
	filter          filter.Filter
	codeToLevelFunc CodeToLevelFunc

	// requestHeaders enables logging of all request headers, however sensitive
	// headers like authorization, cookie and set-cookie are hidden.
	requestHeaders bool

	// hideRequestHeaders are requests headers which are redacted from the logs
	hideRequestHeaders []string

	grpcDetailsToLogAttrsFunc GRPCStatusDetailsToLogAttrsFunc
}

func defaultOptions() *options {
	return &options{
		codeToLevelFunc:           DefaultServerCodeToLevel,
		filter:                    filter.Any(),
		grpcDetailsToLogAttrsFunc: DefaultGRPCStatusDetailsToLogAttrs,
	}
}

type Option func(*options)

func WithCodeToLevelFunc(fn CodeToLevelFunc) Option {
	return func(c *options) {
		c.codeToLevelFunc = fn
	}
}

func WithFilter(filter filter.Filter) Option {
	return func(c *options) {
		c.filter = filter
	}
}

func WithRequestHeaders(enabled bool) Option {
	return func(o *options) {
		o.requestHeaders = enabled
	}
}

func WithHideRequestHeaders(headers []string) Option {
	return func(o *options) {
		o.hideRequestHeaders = headers
	}
}

func WithGRPCDetailsToLogAttrsFunc(fn GRPCStatusDetailsToLogAttrsFunc) Option {
	return func(o *options) {
		o.grpcDetailsToLogAttrsFunc = fn
	}
}

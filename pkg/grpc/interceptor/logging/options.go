package logging

import (
	"log/slog"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
	"google.golang.org/grpc/codes"
)

type CodeToLevel func(code codes.Code) slog.Level

type options struct {
	levelFunc CodeToLevel
	filter    filter.Filter
}

func defaultOptions() *options {
	return &options{
		levelFunc: DefaultServerCodeToLevel,
		filter:    filter.Any(),
	}
}

type Option func(*options)

func WithCodeToLevelFunc(fn CodeToLevel) Option {
	return func(c *options) {
		c.levelFunc = fn
	}
}

func WithFilter(filter filter.Filter) Option {
	return func(c *options) {
		c.filter = filter
	}
}

package errsanitizer

import (
	"context"

	"connectrpc.com/connect"
)

type Interceptor struct {
	snFunc SanitizerFunc
}

var _ connect.Interceptor = &Interceptor{}

func NewInterceptor(opts ...Option) *Interceptor {
	options := defaultOptions()
	for _, fn := range opts {
		fn(options)
	}
	return &Interceptor{snFunc: options.snFunc}
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
		res, err = next(ctx, req)

		return res, i.snFunc(err)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) (err error) {
		err = next(ctx, conn)

		return i.snFunc(err)
	}
}

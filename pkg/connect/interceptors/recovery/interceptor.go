package recovery

import (
	"context"

	"connectrpc.com/connect"
)

type Interceptor struct {
	recoveryFn RecoveryFunc
}

var _ connect.Interceptor = &Interceptor{}

func NewInterceptor(opts ...Option) *Interceptor {
	options := defaultOptions()
	for _, fn := range opts {
		fn(options)
	}
	return &Interceptor{recoveryFn: options.recoveryFn}
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = i.recoveryFn(ctx, p)
			}
		}()

		res, err = next(ctx, req)

		return res, err
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = i.recoveryFn(ctx, p)
			}
		}()

		err = next(ctx, conn)

		return err
	}
}

package recovery

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.recoveryFn(ctx, p)
			}
		}()

		res, err = handler(ctx, req)

		return res, err
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.recoveryFn(stream.Context(), p)
			}
		}()

		err = handler(srv, stream)

		return err
	}
}

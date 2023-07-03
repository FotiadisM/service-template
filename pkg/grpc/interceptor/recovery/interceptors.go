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
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.recoveryFn(ctx, p)
			}
		}()

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.recoveryFn(stream.Context(), p)
			}
		}()

		return handler(srv, stream)
	}
}

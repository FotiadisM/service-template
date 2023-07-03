package filter

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(i grpc.UnaryServerInterceptor, filter Filter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if filter(info.FullMethod) {
			return handler(ctx, req)
		}
		return i(ctx, req, info, handler)
	}
}

func StreamServerInterceptor(i grpc.StreamServerInterceptor, filter Filter) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if filter(info.FullMethod) {
			return handler(srv, ss)
		}
		return i(srv, ss, info, handler)
	}
}

func UnaryClientInterceptor(i grpc.UnaryClientInterceptor, filter Filter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if filter(method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		return i(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func StreamClientInterceptor(i grpc.StreamClientInterceptor, filter Filter) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if filter(method) {
			return streamer(ctx, desc, cc, method, opts...)
		}
		return i(ctx, desc, cc, method, streamer, opts...)
	}
}

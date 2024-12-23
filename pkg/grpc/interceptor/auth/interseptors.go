package auth

import (
	"context"

	"google.golang.org/grpc"
)

type AuthFunc func(ctx context.Context) (context.Context, error) //nolint

func UnaryServerInterceptor(fn AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		newCtx, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// wrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type wrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context //nolint:containedctx
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context().
func (w *wrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// wrapServerStream returns a ServerStream that has the ability to overwrite context.
func wrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}

func StreamServerInterceptor(fn AuthFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		newCtx, err := fn(stream.Context())
		if err != nil {
			return err
		}

		ws := wrapServerStream(stream)
		ws.WrappedContext = newCtx

		return handler(srv, ws)
	}
}

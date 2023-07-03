package logging

import (
	"context"
	"strings"

	"github.com/FotiadisM/mock-microservice/pkg/ilog"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(logger *slog.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if options.filter(info.FullMethod) {
			return handler(ctx, req)
		}

		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if traceID.IsValid() {
			logger = logger.With("trace.id", traceID)
		}

		fullName := strings.TrimLeft(info.FullMethod, "/")
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			ctx = ilog.ContextWithLogger(ctx, logger)
			return handler(ctx, req)
		}

		logger = logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		ctx = ilog.ContextWithLogger(ctx, logger)
		res, err := handler(ctx, req)

		st := status.Convert(err)
		if err != nil {
			logger = logger.With("error", st.Message())
		}

		lvl := options.levelFunc(st.Code())

		logger.LogAttrs(ctx, lvl, "finished call", slog.Int("rpc.grpc.status_code", int(st.Code())))

		return res, err
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

func StreamServerInterceptor(logger *slog.Logger, opts ...Option) grpc.StreamServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if options.filter(info.FullMethod) {
			return handler(srv, ss)
		}

		ws := wrapServerStream(ss)
		ctx := ws.Context()

		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if traceID.IsValid() {
			logger = logger.With("trace.id", traceID)
		}

		fullName := strings.TrimLeft(info.FullMethod, "/")
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			ws.WrappedContext = ilog.ContextWithLogger(ctx, logger)
			return handler(srv, ss)
		}

		logger = logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		ws.WrappedContext = ilog.ContextWithLogger(ctx, logger)
		err := handler(srv, ss)

		st := status.Convert(err)
		if err != nil {
			logger = logger.With("error", st.Message())
		}

		lvl := options.levelFunc(st.Code())

		logger.LogAttrs(ctx, lvl, "finished call", slog.Int("rpc.grpc.status_code", int(st.Code())))

		return err
	}
}

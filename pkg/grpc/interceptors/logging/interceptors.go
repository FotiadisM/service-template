package logging

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/FotiadisM/mock-microservice/pkg/ilog"
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

		fullName := strings.TrimLeft(info.FullMethod, "/")
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			ctx = ilog.ContextWithLogger(ctx, logger)
			return handler(ctx, req)
		}

		ctxLogger := logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if traceID.IsValid() {
			ctxLogger = ctxLogger.With("trace.id", traceID)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if ok && options.requestHeaders {
			for key, value := range md {
				if slices.Index(options.hideRequestHeaders, key) != -1 {
					continue
				}
				ctxLogger = ctxLogger.With("rpc.grpc.request.metadata."+key, value)
			}
		}

		ctx = ilog.ContextWithLogger(ctx, ctxLogger)
		start := time.Now()
		res, err := handler(ctx, req)
		duration := time.Since(start)

		logAttrs := []slog.Attr{slog.Int64("rpc.server.duration", duration.Milliseconds())}
		st := status.Convert(err)
		if err != nil {
			logAttrs = append(logAttrs, ilog.Err(errors.New(st.Message()))) //nolint:err113
		}

		detailsAttrs := options.grpcDetailsToLogAttrsFunc(st.Details())
		if len(detailsAttrs) > 0 {
			logAttrs = append(logAttrs, detailsAttrs...)
		}

		level := options.codeToLevelFunc(st.Code())
		logAttrs = append(logAttrs, slog.Int("rpc.grpc.status_code", int(st.Code())))

		ctxLogger.LogAttrs(ctx, level, "request_end", logAttrs...)

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

		fullName := strings.TrimLeft(info.FullMethod, "/")
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			ws.WrappedContext = ilog.ContextWithLogger(ctx, logger)
			return handler(srv, ss)
		}

		ctxLogger := logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if traceID.IsValid() {
			ctxLogger = ctxLogger.With("trace.id", traceID)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if ok && options.requestHeaders {
			for key, value := range md {
				if slices.Index(options.hideRequestHeaders, key) != -1 {
					continue
				}
				ctxLogger = ctxLogger.With("rpc.grpc.request.metadata."+key, value)
			}
		}

		ws.WrappedContext = ilog.ContextWithLogger(ctx, ctxLogger)
		err := handler(srv, ws)

		logAttrs := []slog.Attr{}
		st := status.Convert(err)
		if err != nil {
			logAttrs = append(logAttrs, ilog.Err(errors.New(st.Message()))) //nolint:err113
		}

		detailsAttrs := options.grpcDetailsToLogAttrsFunc(st.Details())
		if len(detailsAttrs) > 0 {
			logAttrs = append(logAttrs, detailsAttrs...)
		}

		level := options.codeToLevelFunc(st.Code())
		logAttrs = append(logAttrs, slog.Int("rpc.grpc.status_code", int(st.Code())))

		ctxLogger.LogAttrs(ctx, level, "request_end", logAttrs...)

		return err
	}
}

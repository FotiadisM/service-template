package logging

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/FotiadisM/mock-microservice/pkg/ilog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				attrs := []slog.Attr{}
				for i, fv := range t.FieldViolations {
					attrs = append(attrs, slog.Group(strconv.Itoa(i),
						slog.String("field", fv.Field),
						slog.String("description", fv.Description),
					))
				}
				logger = logger.With(attrs)
			case *errdetails.DebugInfo:
			case *errdetails.ErrorInfo:
				md := []slog.Attr{}
				for k, v := range t.Metadata {
					md = append(md, slog.String(k, v))
				}
				logger = logger.With(slog.Group(
					"error_info",
					slog.String("reason", t.Reason),
					slog.String("domain", t.Domain),
					slog.Group("metadata", md),
				))
			case *errdetails.PreconditionFailure:
			case *errdetails.RequestInfo:
			}
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

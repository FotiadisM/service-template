package logger

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	slog.Handler
}

func NewHandler() slog.Handler {
	h := slog.NewJSONHandler(os.Stdout, nil)
	return handler{h}
}

type ctxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

// DefaultServerCodeToLevel is the helper mapper that maps gRPC return codes to log levels for server side.
func DefaultServerCodeToLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK, codes.NotFound, codes.Canceled, codes.AlreadyExists, codes.InvalidArgument, codes.Unauthenticated:
		return slog.LevelInfo

	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted,
		codes.OutOfRange, codes.Unavailable:
		return slog.LevelWarn

	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
		return slog.LevelError

	default:
		return slog.LevelError
	}
}

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
			ctx = ContextWithLogger(ctx, logger)
			return handler(ctx, req)
		}

		logger = logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		ctx = ContextWithLogger(ctx, logger)
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
			ws.WrappedContext = ContextWithLogger(ctx, logger)
			return handler(srv, ss)
		}

		logger = logger.With(
			"rpc.service", parts[0],
			"rpc.method", parts[1],
		)

		ws.WrappedContext = ContextWithLogger(ctx, logger)
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

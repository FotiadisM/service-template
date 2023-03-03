package logger

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var defaultLogger *zap.Logger

type ctxKey struct{}

func New(debug bool) *zap.Logger {
	config := zap.NewProductionConfig()
	if debug {
		config = zap.NewDevelopmentConfig()
	}

	defaultLogger = zap.Must(config.Build())

	return defaultLogger
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return logger
	}

	return defaultLogger
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fields := logging.ExtractFields(ctx)
		vals := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			vals = append(vals, zap.String(fields[i], fields[i+1]))
		}
		logger = logger.With(vals...)
		ctx = context.WithValue(ctx, ctxKey{}, logger)

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newStream := middleware.WrapServerStream(ss)
		ctx := newStream.Context()

		fields := logging.ExtractFields(ctx)
		vals := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			vals = append(vals, zap.String(fields[i], fields[i+1]))
		}
		logger = logger.With(vals...)
		newStream.WrappedContext = context.WithValue(ctx, ctxKey{}, logger)

		return handler(srv, newStream)
	}
}

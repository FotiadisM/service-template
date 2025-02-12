package logging

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"slices"
	"strings"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"

	"github.com/FotiadisM/service-template/pkg/ilog"
)

const (
	systemKey         = "rpc.system"
	methodKey         = "rpc.method"
	serviceKey        = "rpc.service"
	peerNameKey       = "net.peer.name"
	peerPortKey       = "net.peer.port"
	serverDurationKey = "rpc.server.duration"
	metadataPrefixKey = "rpc.connect_rpc.request.metadata."

	grpcProtocol    = "grpc"
	grpcwebString   = "grpcweb"
	grpcwebProtocol = "grpc_web"
	connectString   = "connect"
	connectProtocol = "connect_rpc"
)

func protocolAttribute(protocol string) slog.Attr {
	switch protocol {
	case grpcwebString:
		return slog.String(systemKey, grpcwebProtocol)
	case grpcProtocol:
		return slog.String(systemKey, grpcProtocol)
	case connectString:
		return slog.String(systemKey, connectProtocol)
	default:
		return slog.String(systemKey, protocol)
	}
}

func addressAttributes(address string) []slog.Attr {
	if host, port, err := net.SplitHostPort(address); err == nil {
		return []slog.Attr{
			slog.String(peerNameKey, host),
			slog.String(peerPortKey, port),
		}
	}
	return []slog.Attr{slog.String(peerNameKey, address)}
}

func procedureAttributes(procedure string) []any {
	name := strings.TrimLeft(procedure, "/")
	parts := strings.SplitN(name, "/", 2)
	var attrs []any

	switch len(parts) {
	case 0:
		return attrs // invalid
	case 1:
		// fall back to treating the whole string as the method
		if method := parts[0]; method != "" {
			attrs = append(attrs, slog.String(methodKey, method))
		}
	default:
		if svc := parts[0]; svc != "" {
			attrs = append(attrs, slog.String(serviceKey, svc))
		}
		if method := parts[1]; method != "" {
			attrs = append(attrs, slog.String(methodKey, method))
		}
	}

	return attrs
}

type Interceptor struct {
	opts   *options
	logger *slog.Logger
}

var _ connect.Interceptor = &Interceptor{}

func NewInterceptor(logger *slog.Logger, opts ...Option) *Interceptor {
	options := defaultOptions()
	for _, fn := range opts {
		fn(options)
	}
	return &Interceptor{opts: options, logger: logger}
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient || !i.opts.filterFunc(ctx, req.Spec()) {
			return next(ctx, req)
		}

		ctxLogger := i.logger.With(procedureAttributes(req.Spec().Procedure)...)

		span := trace.SpanContextFromContext(ctx)
		if span.IsValid() {
			ctxLogger = ctxLogger.With("trace.id", span.TraceID(), "span.id", span.SpanID())
		}

		ctx = ilog.ContextWithLogger(ctx, ctxLogger)
		start := time.Now()
		res, err := next(ctx, req)
		duration := time.Since(start)

		level := slog.LevelInfo
		logAttrs := []slog.Attr{
			protocolAttribute(req.Peer().Protocol),
			slog.Int64(serverDurationKey, duration.Milliseconds()),
		}
		if i.opts.withPeer {
			logAttrs = append(logAttrs, addressAttributes(req.Peer().Addr)...)
		}
		if err != nil {
			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				level = i.opts.codeToLevelFunc(connectErr.Code())
				logAttrs = append(logAttrs,
					ilog.Err(errors.New(connectErr.Message())), //nolint:err113
					slog.String("rpc.connect_rpc.error_code", connectErr.Code().String()),
				)
				logAttrs = append(logAttrs, i.opts.errorDetailsAttrFunc(connectErr.Details())...)
			} else {
				level = slog.LevelError
				logAttrs = append(logAttrs,
					ilog.Err(err),
					slog.String("rpc.connect_rpc.error_code", connect.CodeInternal.String()),
				)
			}
		}

		if i.opts.withRequestsHeaders {
			for k, v := range req.Header() {
				if slices.Index(i.opts.hiddenRequestHeaders, k) != -1 {
					continue
				}
				logAttrs = append(logAttrs, slog.Any(metadataPrefixKey+strings.ToLower(k), v))
			}
		}

		ctxLogger.LogAttrs(ctx, level, "request_end", logAttrs...)

		return res, err
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if !i.opts.filterFunc(ctx, conn.Spec()) {
			return next(ctx, conn)
		}

		ctxLogger := i.logger.With(procedureAttributes(conn.Spec().Procedure)...)

		span := trace.SpanContextFromContext(ctx)
		if span.IsValid() {
			ctxLogger = ctxLogger.With("trace.id", span.TraceID(), "span.id", span.SpanID())
		}

		ctx = ilog.ContextWithLogger(ctx, ctxLogger)
		start := time.Now()
		err := next(ctx, conn)
		duration := time.Since(start)

		level := slog.LevelInfo
		logAttrs := []slog.Attr{
			protocolAttribute(conn.Peer().Protocol),
			slog.Int64(serverDurationKey, duration.Milliseconds()),
		}
		if i.opts.withPeer {
			logAttrs = append(logAttrs, addressAttributes(conn.Peer().Addr)...)
		}
		if err != nil {
			if connectErr := new(connect.Error); errors.As(err, &connectErr) {
				level = i.opts.codeToLevelFunc(connectErr.Code())
				logAttrs = append(logAttrs, ilog.Err(errors.New(connectErr.Message()))) //nolint:err113
				logAttrs = append(logAttrs, i.opts.errorDetailsAttrFunc(connectErr.Details())...)
			} else {
				level = slog.LevelError
				logAttrs = append(logAttrs, ilog.Err(err))
			}
		}

		if i.opts.withRequestsHeaders {
			for k, v := range conn.RequestHeader() {
				if slices.Index(i.opts.hiddenRequestHeaders, k) != -1 {
					continue
				}
				logAttrs = append(logAttrs, slog.Any(metadataPrefixKey+strings.ToLower(k), v))
			}
		}

		ctxLogger.LogAttrs(ctx, level, "request_end", logAttrs...)

		return err
	}
}

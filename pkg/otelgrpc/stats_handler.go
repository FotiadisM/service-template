package otelgrpc

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

const instrumentationName = "github.com/FotiadisM/otelgrpc"

type serverConfig struct {
	propagator     propagation.TextMapPropagator
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
	errorHandler   otel.ErrorHandler
}

func newServerConfig() *serverConfig {
	return &serverConfig{
		propagator:     otel.GetTextMapPropagator(),
		tracerProvider: otel.GetTracerProvider(),
		meterProvider:  otel.GetMeterProvider(),
		errorHandler:   otel.GetErrorHandler(),
	}
}

type ServerOption interface {
	apply(c *serverConfig)
}

type serverOptFunc func(c *serverConfig)

func (f serverOptFunc) apply(c *serverConfig) {
	f(c)
}

func WithTextMapPropagator(mp propagation.TextMapPropagator) ServerOption {
	return serverOptFunc(func(c *serverConfig) {
		c.propagator = mp
	})
}

func WithTracerProvider(tp trace.TracerProvider) ServerOption {
	return serverOptFunc(func(c *serverConfig) {
		c.tracerProvider = tp
	})
}

func WithMeterProvider(mp metric.MeterProvider) ServerOption {
	return serverOptFunc(func(c *serverConfig) {
		c.meterProvider = mp
	})
}

func WithErrorHandler(eh otel.ErrorHandler) ServerOption {
	return serverOptFunc(func(c *serverConfig) {
		c.errorHandler = eh
	})
}

func ServerStatsHandler(options ...ServerOption) stats.Handler {
	config := newServerConfig()
	for _, o := range options {
		o.apply(config)
	}

	tracer := config.tracerProvider.Tracer(
		instrumentationName,
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	meter := config.meterProvider.Meter(
		instrumentationName,
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	handler := &serverStatsHandler{
		propagator:   config.propagator,
		tracer:       tracer,
		meter:        meter,
		errorHandler: config.errorHandler,
	}

	var err error
	if handler.duration, err = meter.Int64Histogram(
		"rpc.server.duration",
		instrument.WithUnit("ms"),
		instrument.WithDescription("measures duration of inbound RPC"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requestSize, err = meter.Int64Histogram(
		"rpc.server.request.size",
		instrument.WithUnit("By"),
		instrument.WithDescription("measures size of RPC request messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responseSize, err = meter.Int64Histogram(
		"rpc.server.response.size",
		instrument.WithUnit("By"),
		instrument.WithDescription("measures size of RPC response messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requests, err = meter.Int64Histogram(
		"rpc.server.requests_per_rpc",
		instrument.WithUnit("{count}"),
		instrument.WithDescription("measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responses, err = meter.Int64Histogram(
		"rpc.server.responses_per_rpc",
		instrument.WithUnit("{count}"),
		instrument.WithDescription("measures the number of messages sent per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}

	return handler
}

type serverObserverCtxKey struct{}

type serverObserver struct {
	msgSentCount    int
	msgReceiveCount int

	// isStreaming is used to avoid measuring duration in streaming RPCs
	isStreaming bool

	attrs []attribute.KeyValue
}

type serverStatsHandler struct {
	propagator   propagation.TextMapPropagator
	tracer       trace.Tracer
	meter        metric.Meter
	errorHandler otel.ErrorHandler

	duration     instrument.Int64Histogram
	requestSize  instrument.Int64Histogram
	responseSize instrument.Int64Histogram
	requests     instrument.Int64Histogram
	responses    instrument.Int64Histogram
}

// assert that serverStatsHandler implements the stats.Handler interface.
var _ stats.Handler = &serverStatsHandler{}

func (h *serverStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	name, attrs := spanInfo(info.FullMethodName)

	ctx, _ = h.tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)

	ovserver := &serverObserver{
		attrs: attrs,
	}

	return context.WithValue(ctx, serverObserverCtxKey{}, ovserver)
}

func (h *serverStatsHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	span := trace.SpanFromContext(ctx)

	observer, ok := ctx.Value(serverObserverCtxKey{}).(*serverObserver)
	if !ok {
		observer = &serverObserver{}
	}

	switch rs := rpcStats.(type) {
	case *stats.InHeader:
		host, portStr, err := net.SplitHostPort(rs.LocalAddr.String())
		if err != nil {
			break
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return
		}

		attrs := []attribute.KeyValue{
			semconv.NetPeerName(host),
			semconv.NetPeerPort(port),
		}
		observer.attrs = append(observer.attrs, attrs...)
		span.SetAttributes(attrs...)

	case *stats.Begin:
		observer.isStreaming = rs.IsClientStream || rs.IsServerStream

	case *stats.InPayload:
		h.requestSize.Record(ctx, int64(rs.Length), observer.attrs...)

		observer.msgReceiveCount++
		span.AddEvent("message", trace.WithAttributes(
			semconv.MessageTypeReceived,
			semconv.MessageID(observer.msgReceiveCount),
			semconv.MessageUncompressedSize(rs.Length),
			semconv.MessageCompressedSize(rs.WireLength),
		))

	case *stats.OutPayload:
		h.responseSize.Record(ctx, int64(rs.Length), observer.attrs...)

		observer.msgSentCount++
		span.AddEvent("message", trace.WithAttributes(
			semconv.MessageTypeSent,
			semconv.MessageID(observer.msgSentCount),
			semconv.MessageUncompressedSize(rs.Length),
			semconv.MessageCompressedSize(rs.WireLength),
		))

	case *stats.End:
		rpcCode := grpcCodes.OK
		rpcMesg := ""
		if rs.Error != nil {
			st, ok := status.FromError(rs.Error)
			if ok {
				rpcCode = st.Code()
				rpcMesg = st.Message()
			} else {
				rpcCode = grpcCodes.Internal
				rpcMesg = rs.Error.Error()
			}
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(rpcCode)))
		observer.attrs = append(observer.attrs, semconv.RPCGRPCStatusCodeKey.Int(int(rpcCode)))
		if rpcCode != grpcCodes.OK {
			span.SetStatus(codes.Error, rpcMesg)
		}

		if !observer.isStreaming {
			duration := rs.EndTime.Sub(rs.BeginTime) / time.Millisecond
			h.duration.Record(ctx, int64(duration), observer.attrs...)
		}
		h.requests.Record(ctx, int64(observer.msgReceiveCount), observer.attrs...)
		h.responses.Record(ctx, int64(observer.msgSentCount), observer.attrs...)

		span.End()
	}
}

func (h *serverStatsHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *serverStatsHandler) HandleConn(_ context.Context, _ stats.ConnStats) {}

// spanInfo returns a span name and all appropriate attributes from the gRPC method.
func spanInfo(fullMethod string) (string, []attribute.KeyValue) {
	name, attrs := parseFullMethod(fullMethod)
	attrs = append(attrs, semconv.RPCSystemGRPC)
	return name, attrs
}

// parseFullMethod returns a span name following the OpenTelemetry semantic
// conventions as well as all applicable span attribute.KeyValue attributes based
// on a gRPC's FullMethod.
func parseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCService(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethod(method))
	}
	return name, attrs
}

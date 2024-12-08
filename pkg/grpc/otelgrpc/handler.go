package otelgrpc

import (
	"context"
	"net"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
)

func TraceAttributes(ctx context.Context) (traceID string, attrs []attribute.KeyValue, ok bool) {
	sc := trace.SpanContextFromContext(ctx)
	t := sc.TraceID()
	if !t.IsValid() {
		return "", nil, false
	}

	gctx, _ := ctx.Value(grpcContextKey{}).(*grpcContext)
	if gctx.attrs == nil || gctx.skip {
		return "", nil, false
	}

	return t.String(), gctx.attrs, ok
}

type grpcContextKey struct{}

type grpcContext struct {
	skip bool

	msgSentCount    int
	msgReceiveCount int

	// isStreaming is used to avoid measuring duration in streaming RPCs
	isStreaming bool

	attrs []attribute.KeyValue
}

type statsHandler struct {
	filter   filter.Filter
	spanKind trace.SpanKind

	propagator   propagation.TextMapPropagator
	tracer       trace.Tracer
	meter        metric.Meter
	errorHandler otel.ErrorHandler

	requestMetadata  bool
	responseMetadata bool

	duration     metric.Int64Histogram
	requestSize  metric.Int64Histogram
	responseSize metric.Int64Histogram
	requests     metric.Int64Histogram
	responses    metric.Int64Histogram
}

// assert that statsHandler implements the stats.Handler interface.
var _ stats.Handler = &statsHandler{}

func (h *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if h.filter != nil && h.filter(info.FullMethodName) {
		gctx := &grpcContext{skip: true}
		return context.WithValue(ctx, grpcContextKey{}, gctx)
	}

	ctx = extract(ctx, h.propagator)

	name, attrs := spanInfo(info.FullMethodName)

	ctx, _ = h.tracer.Start(ctx, name,
		trace.WithSpanKind(h.spanKind),
		trace.WithAttributes(attrs...),
	)

	gctx := &grpcContext{
		attrs: attrs,
	}

	return context.WithValue(ctx, grpcContextKey{}, gctx)
}

func (h *statsHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	gctx, _ := ctx.Value(grpcContextKey{}).(*grpcContext)
	if gctx.skip {
		return
	}

	span := trace.SpanFromContext(ctx)

	switch rs := rpcStats.(type) {
	case *stats.InHeader:
		attrs := []attribute.KeyValue{}
		for k, values := range rs.Header {
			attrValue := "["
			for i, v := range values {
				attrValue += v
				if i != len(values)-1 {
					attrValue += ","
				}
			}
			attrs = append(attrs, attribute.String("rpc.grpc.metadata."+k, attrValue))
		}
		span.SetAttributes(attrs...)

		// TODO(FotiadisM): validate server and client
		if addr, ok := rs.RemoteAddr.(*net.TCPAddr); ok {
			gctx.attrs = append(gctx.attrs, semconv.ServerAddress(addr.IP.String()), semconv.ServerPort(addr.Port))
			span.SetAttributes(semconv.ServerAddress(addr.IP.String()), semconv.ServerPort(addr.Port))
		}
		if addr, ok := rs.LocalAddr.(*net.TCPAddr); ok {
			gctx.attrs = append(gctx.attrs, semconv.ClientAddress(addr.IP.String()), semconv.ClientPort(addr.Port))
			span.SetAttributes(semconv.ClientAddress(addr.IP.String()), semconv.ClientPort(addr.Port))
		}

		// TODO(FotiadisM): add request metadata

	case *stats.Begin:
		gctx.isStreaming = rs.IsClientStream || rs.IsServerStream

	case *stats.InPayload:
		h.requestSize.Record(ctx, int64(rs.Length), metric.WithAttributes(gctx.attrs...))

		if gctx.isStreaming {
			gctx.msgReceiveCount++
			span.AddEvent("message", trace.WithAttributes(
				semconv.MessageTypeReceived,
				semconv.MessageID(gctx.msgReceiveCount),
				semconv.MessageUncompressedSize(rs.Length),
				semconv.MessageCompressedSize(rs.WireLength),
			))
		}

	case *stats.OutPayload:
		h.responseSize.Record(ctx, int64(rs.Length), metric.WithAttributes(gctx.attrs...))

		if gctx.isStreaming {
			gctx.msgSentCount++
			span.AddEvent("message", trace.WithAttributes(
				semconv.MessageTypeSent,
				semconv.MessageID(gctx.msgSentCount),
				semconv.MessageUncompressedSize(rs.Length),
				semconv.MessageCompressedSize(rs.WireLength),
			))
		}

	case *stats.End:
		rpcCode := grpcCodes.OK
		rpcMsg := ""
		if rs.Error != nil {
			rpcCode = grpcCodes.Internal
			rpcMsg = rs.Error.Error()
			if st, ok := status.FromError(rs.Error); ok {
				rpcCode = st.Code()
				rpcMsg = st.Message()
			}
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(rpcCode)))
		gctx.attrs = append(gctx.attrs, semconv.RPCGRPCStatusCodeKey.Int(int(rpcCode)))
		if rpcCode != grpcCodes.OK { // set error code propely
			span.SetStatus(codes.Error, rpcMsg)
		}

		if gctx.isStreaming {
			h.requests.Record(ctx, int64(gctx.msgReceiveCount), metric.WithAttributes(gctx.attrs...))
			h.responses.Record(ctx, int64(gctx.msgSentCount), metric.WithAttributes(gctx.attrs...))
		} else {
			duration := rs.EndTime.Sub(rs.BeginTime) / time.Millisecond
			h.duration.Record(ctx, int64(duration), metric.WithAttributes(gctx.attrs...))
		}

		span.End()
	}
}

func (h *statsHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *statsHandler) HandleConn(_ context.Context, _ stats.ConnStats) {}

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

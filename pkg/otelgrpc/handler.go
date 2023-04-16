package otelgrpc

import (
	"context"
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

type statsHandler struct {
	filter   Filter
	spanKind trace.SpanKind

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

// assert that statsHandler implements the stats.Handler interface.
var _ stats.Handler = &statsHandler{}

func (h *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if h.filter != nil && h.filter(info.FullMethodName) {
		observer := &observer{skip: true}
		return ctxWithObserver(ctx, observer)
	}

	name, attrs := spanInfo(info.FullMethodName)

	ctx, _ = h.tracer.Start(ctx, name,
		trace.WithSpanKind(h.spanKind),
		trace.WithAttributes(attrs...),
	)

	observer := &observer{
		attrs: attrs,
	}

	return ctxWithObserver(extract(ctx, h.propagator), observer)
}

func (h *statsHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	observer := observerFromCtx(ctx)
	if observer.skip {
		return
	}

	span := trace.SpanFromContext(ctx)

	switch rs := rpcStats.(type) {
	case *stats.InHeader:
		// TODO(FotiadisM): fix `net.peer.name`

		// host, portStr, err := net.SplitHostPort(rs.LocalAddr.String())
		// if err != nil {
		// 	break
		// }

		// port, err := strconv.Atoi(portStr)
		// if err != nil {
		// 	return
		// }

		// attrs := []attribute.KeyValue{
		// 	semconv.NetPeerName(host),
		// 	semconv.NetPeerPort(port),
		// }
		// observer.attrs = append(observer.attrs, attrs...)
		// span.SetAttributes(attrs...)

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

		if observer.isStreaming {
			h.requests.Record(ctx, int64(observer.msgReceiveCount), observer.attrs...)
			h.responses.Record(ctx, int64(observer.msgSentCount), observer.attrs...)
		} else {
			duration := rs.EndTime.Sub(rs.BeginTime) / time.Millisecond
			h.duration.Record(ctx, int64(duration), observer.attrs...)
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

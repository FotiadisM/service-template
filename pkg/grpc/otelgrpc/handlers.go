package otelgrpc

import (
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/stats"
)

const instrumentationName = "github.com/FotiadisM/otelgrpc"

func newStatsHandler(opts []Option, spanKind trace.SpanKind, role string) stats.Handler {
	options := defaultOptions()
	for _, fn := range opts {
		fn(options)
	}

	tracer := options.tracerProvider.Tracer(
		instrumentationName,
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	meter := options.meterProvider.Meter(
		instrumentationName,
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	handler := &statsHandler{
		filter:           options.filter,
		spanKind:         spanKind,
		propagator:       options.propagator,
		tracer:           tracer,
		meter:            meter,
		errorHandler:     options.errorHandler,
		requestMetadata:  options.requestMetadata,
		responseMetadata: options.responseMetadata,
	}

	var err error
	rpcrole := "rpc." + role + "."
	if handler.duration, err = meter.Int64Histogram(
		rpcrole+"duration",
		metric.WithUnit("ms"),
		metric.WithDescription("measures duration of inbound RPC"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requestSize, err = meter.Int64Histogram(
		rpcrole+"request.size",
		metric.WithUnit("By"),
		metric.WithDescription("measures size of RPC request messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responseSize, err = meter.Int64Histogram(
		rpcrole+"response.size",
		metric.WithUnit("By"),
		metric.WithDescription("measures size of RPC response messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requests, err = meter.Int64Histogram(
		rpcrole+"requests_per_rpc",
		metric.WithUnit("{count}"),
		metric.WithDescription("measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responses, err = meter.Int64Histogram(
		rpcrole+"responses_per_rpc",
		metric.WithUnit("{count}"),
		metric.WithDescription("measures the number of messages sent per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}

	return handler
}

func ServerStatsHandler(opts ...Option) stats.Handler {
	return newStatsHandler(opts, trace.SpanKindServer, "server")
}

func ClientStatsHandler(opts ...Option) stats.Handler {
	return newStatsHandler(opts, trace.SpanKindClient, "client")
}

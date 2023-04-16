package otelgrpc

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/stats"
)

const instrumentationName = "github.com/FotiadisM/otelgrpc"

func ServerStatsHandler(options ...Option) stats.Handler {
	config := newConfig()
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

	handler := &statsHandler{
		filter:       config.filter,
		spanKind:     trace.SpanKindServer,
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

func ClientStatsHandler(options ...Option) stats.Handler {
	config := newConfig()
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

	handler := &statsHandler{
		filter:       config.filter,
		spanKind:     trace.SpanKindClient,
		propagator:   config.propagator,
		tracer:       tracer,
		meter:        meter,
		errorHandler: config.errorHandler,
	}

	var err error
	if handler.duration, err = meter.Int64Histogram(
		"rpc.client.duration",
		instrument.WithUnit("ms"),
		instrument.WithDescription("measures duration of inbound RPC"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requestSize, err = meter.Int64Histogram(
		"rpc.client.request.size",
		instrument.WithUnit("By"),
		instrument.WithDescription("measures size of RPC request messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responseSize, err = meter.Int64Histogram(
		"rpc.client.response.size",
		instrument.WithUnit("By"),
		instrument.WithDescription("measures size of RPC response messages (uncompressed)"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.requests, err = meter.Int64Histogram(
		"rpc.client.requests_per_rpc",
		instrument.WithUnit("{count}"),
		instrument.WithDescription("measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}
	if handler.responses, err = meter.Int64Histogram(
		"rpc.client.responses_per_rpc",
		instrument.WithUnit("{count}"),
		instrument.WithDescription("measures the number of messages sent per RPC. Should be 1 for all non-streaming RPCs"),
	); err != nil {
		handler.errorHandler.Handle(err)
	}

	return handler
}

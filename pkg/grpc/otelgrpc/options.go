package otelgrpc

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
)

type options struct {
	filter         filter.Filter
	propagator     propagation.TextMapPropagator
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
	errorHandler   otel.ErrorHandler

	requestMetadata  bool
	responseMetadata bool
}

func defaultOptions() *options {
	return &options{
		propagator:       otel.GetTextMapPropagator(),
		tracerProvider:   otel.GetTracerProvider(),
		meterProvider:    otel.GetMeterProvider(),
		errorHandler:     otel.GetErrorHandler(),
		requestMetadata:  false,
		responseMetadata: false,
	}
}

type Option func(c *options)

func WithTextMapPropagator(mp propagation.TextMapPropagator) Option {
	return func(c *options) {
		c.propagator = mp
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(c *options) {
		c.tracerProvider = tp
	}
}

func WithMeterProvider(mp metric.MeterProvider) Option {
	return func(c *options) {
		c.meterProvider = mp
	}
}

func WithErrorHandler(eh otel.ErrorHandler) Option {
	return func(c *options) {
		c.errorHandler = eh
	}
}

func WithFilter(filter filter.Filter) Option {
	return func(c *options) {
		c.filter = filter
	}
}

func WithRequestMetadata() Option {
	return func(c *options) {
		c.requestMetadata = true
	}
}

func WithResponseMetadata() Option {
	return func(c *options) {
		c.responseMetadata = true
	}
}

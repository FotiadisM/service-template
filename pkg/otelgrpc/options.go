package otelgrpc

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	filter         Filter
	propagator     propagation.TextMapPropagator
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
	errorHandler   otel.ErrorHandler
}

func defaultConfig() *config {
	return &config{
		propagator:     otel.GetTextMapPropagator(),
		tracerProvider: otel.GetTracerProvider(),
		meterProvider:  otel.GetMeterProvider(),
		errorHandler:   otel.GetErrorHandler(),
	}
}

type Option func(c *config)

func WithTextMapPropagator(mp propagation.TextMapPropagator) Option {
	return func(c *config) {
		c.propagator = mp
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(c *config) {
		c.tracerProvider = tp
	}
}

func WithMeterProvider(mp metric.MeterProvider) Option {
	return func(c *config) {
		c.meterProvider = mp
	}
}

func WithErrorHandler(eh otel.ErrorHandler) Option {
	return func(c *config) {
		c.errorHandler = eh
	}
}

type Filter func(fullMethodName string) bool

func WithFilter(filter Filter) Option {
	return func(c *config) {
		c.filter = filter
	}
}

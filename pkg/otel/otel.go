package otel

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	errorHandler otel.ErrorHandler
}

type Option func(c *config)

func WithErrorHandler(eh otel.ErrorHandler) Option {
	return func(c *config) {
		c.errorHandler = eh
	}
}

func WithErrorHandlerFunc(eh otel.ErrorHandlerFunc) Option {
	return func(c *config) {
		c.errorHandler = eh
	}
}

type Otel struct {
	config *config

	traceProvider *tracesdk.TracerProvider
	meterProvider *metricsdk.MeterProvider
}

func New(ctx context.Context, opts ...Option) (*Otel, error) {
	cfg := &config{}
	for _, o := range opts {
		o(cfg)
	}

	if cfg.errorHandler != nil {
		otel.SetErrorHandler(cfg.errorHandler)
	}

	rsrc := resource.Default()

	te, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp trace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(te),
		tracesdk.WithResource(rsrc),
	)
	otel.SetTracerProvider(tp)

	me, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp metric exporter: %w", err)
	}

	reader := metricsdk.NewPeriodicReader(me)

	mp := metricsdk.NewMeterProvider(
		metricsdk.WithReader(reader),
		metricsdk.WithResource(rsrc),
	)
	otel.SetMeterProvider(mp)

	o := &Otel{
		config:        cfg,
		traceProvider: tp,
		meterProvider: mp,
	}

	return o, nil
}

func (o *Otel) TracerProvider() trace.TracerProvider {
	return o.traceProvider
}

func (o *Otel) MeterProvider() metric.MeterProvider {
	return o.meterProvider
}

func (o *Otel) Shutdown(ctx context.Context) error {
	errs := make(chan error, 2)
	go func() {
		errs <- o.traceProvider.Shutdown(ctx)
	}()
	go func() {
		errs <- o.meterProvider.Shutdown(ctx)
	}()

	return errors.Join(<-errs, <-errs)
}

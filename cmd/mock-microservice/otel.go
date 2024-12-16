package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"

	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type otelShutDownFunc func(ctx context.Context) error

func otelgrpcFilter(ri *stats.RPCTagInfo) bool {
	fullName := strings.TrimLeft(ri.FullMethodName, "/")
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return true
	}
	service := parts[0]

	switch service {
	case "grpc.reflection.v1.ServerReflection":
		return false
	case "grpc.reflection.v1alpha.ServerReflection":
		return false
	case "grpc.health.v1.Health":
		return false
	}

	return true
}

func initializeOTEL(ctx context.Context, log *slog.Logger, exporterAddr string) (otelShutDownFunc, error) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Error("open-telemtry", ilog.Err(err.Error()))
	}))

	conn, err := grpc.NewClient(exporterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel resource %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	fn := func(ctx context.Context) error {
		err := tracerProvider.Shutdown(ctx)
		return errors.Join(err, meterProvider.Shutdown(ctx))
	}

	return fn, nil
}

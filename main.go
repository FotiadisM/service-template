package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/server"
	"github.com/FotiadisM/mock-microservice/internal/service/health"
	servicev1 "github.com/FotiadisM/mock-microservice/internal/service/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	"github.com/FotiadisM/mock-microservice/pkg/logger"
)

type Config struct {
	Store  store.Config
	Server server.Config
}

func main() {
	ctx := context.Background()

	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to process config; %v", err)
		os.Exit(1)
	}

	log := logger.New(config.Server.Debug)

	rs, err := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceName("mock-microservice"),
		semconv.ServiceVersion("0.0.1"),
		semconv.DeploymentEnvironment("dev"),
	))
	if err != nil {
		log.Fatal("failed to create resource", zap.Error(err))
	}

	te, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatal("failed to create trace exporter", zap.Error(err))
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(rs),
		trace.WithBatcher(te),
	)
	otel.SetTracerProvider(tp)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	me, err := stdoutmetric.New(
		stdoutmetric.WithEncoder(enc),
	)
	if err != nil {
		log.Fatal("failed to create metrics exporter", zap.Error(err))
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(rs),
		metric.WithReader(metric.NewPeriodicReader(me)),
	)
	otel.SetMeterProvider(mp)

	store, err := store.New(ctx, config.Store)
	if err != nil {
		log.Fatal("failed to create store", zap.Error(err))
	}

	svc := servicev1.NewService(store)
	healthSvc := health.NewService(store)

	server := server.New(config.Server, log)
	server.Configure()
	server.RegisterService(func(s *grpc.Server, m *runtime.ServeMux) {
		authv1.RegisterAuthServiceServer(s, svc)
		healthv1.RegisterHealthServer(s, healthSvc)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, m, config.Server.GRPCAddr, opts); err != nil {
			log.Fatal("failed to register server", zap.Error(err))
		}
	})
	server.Start(ctx)

	tp.Shutdown(ctx) // nolint:errcheck
	mp.Shutdown(ctx) // nolint:errcheck
}

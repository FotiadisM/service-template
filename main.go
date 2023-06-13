package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/server"
	servicev1 "github.com/FotiadisM/mock-microservice/internal/service/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	"github.com/FotiadisM/mock-microservice/pkg/health"
	"github.com/FotiadisM/mock-microservice/pkg/logger"
	"github.com/FotiadisM/mock-microservice/pkg/otel"

	"github.com/FotiadisM/mock-microservice/pkg/version"
)

type Config struct {
	Store  store.Config
	Server server.Config
}

func main() {
	version.AddFlag(nil)
	flag.Parse()

	ctx := context.Background()

	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to process config; %v", err)
		os.Exit(1)
	}

	log := logger.New(config.Server.Debug)

	otel, err := otel.New(ctx,
		otel.WithErrorHandlerFunc(func(err error) {
			log.Error("otel error occurred", zap.Error(err))
		}),
	)
	if err != nil {
		log.Fatal("failed to initiliaze otel", zap.Error(err))
	}

	store, err := store.New(ctx, config.Store)
	if err != nil {
		log.Fatal("failed to create store", zap.Error(err))
	}

	svc := servicev1.NewService(store)
	healthSvc := health.NewService(nil, nil, nil)

	server := server.New(config.Server, log)
	server.Configure()
	server.RegisterService(func(s *grpc.Server, m *runtime.ServeMux) {
		authv1.RegisterAuthServiceServer(s, svc)
		healthv1.RegisterHealthServer(s, healthSvc)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err = authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, m, config.Server.GRPCAddr, opts); err != nil {
			log.Fatal("failed to register server", zap.Error(err))
		}
	})
	server.Start()
	if err = server.AwaitShutdown(ctx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	}

	if err = otel.Shutdown(ctx); err != nil {
		log.Error("otel shutdown failed", zap.Error(err))
	}
}

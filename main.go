package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	apiauthv1 "github.com/FotiadisM/mock-microservice/api/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/server"
	svcauthv1 "github.com/FotiadisM/mock-microservice/internal/service/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/health"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"

	"github.com/FotiadisM/mock-microservice/pkg/version"
)

//go:generate mockery

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

	log := slog.New(ilog.NewHandler(nil))

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Error("received otel error", ilog.Err(err))
	}))

	store, err := store.New(config.Store)
	if err != nil {
		log.Error("failed to create store", ilog.Err(err))
		os.Exit(1)
	}

	svc := svcauthv1.NewService(store)
	healthSvc := health.NewService()

	server := server.New(config.Server, log)
	server.RegisterService(func(s *grpc.Server, m *runtime.ServeMux) {
		apiauthv1.RegisterAuthServiceServer(s, svc)
		healthv1.RegisterHealthServer(s, healthSvc)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err = apiauthv1.RegisterAuthServiceHandlerFromEndpoint(ctx, m, config.Server.GRPCAddr, opts); err != nil {
			log.Error("failed to register server", ilog.Err(err))
			os.Exit(1)
		}
	})
	server.Start()
	if err = server.GracefulStop(); err != nil {
		log.Error("server shutdown failed", ilog.Err(err))
	}
}

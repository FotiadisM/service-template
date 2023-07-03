package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	authv1 "github.com/FotiadisM/mock-microservice/api/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/server"
	servicev1 "github.com/FotiadisM/mock-microservice/internal/service/v1"
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

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// otel, err := otel.New(ctx,
	// 	otel.WithErrorHandlerFunc(func(err error) {
	// 		log.Error("otel error occurred", "error", err.Error())
	// 	}),
	// )
	// if err != nil {
	// 	log.Error("failed to initiliaze otel", "error", err.Error())
	// 	os.Exit(1)
	// }

	store, err := store.New(config.Store)
	if err != nil {
		log.Error("failed to create store", ilog.Err(err))
		os.Exit(1)
	}

	svc := servicev1.NewService(store)
	healthSvc := health.NewService(nil, nil, nil)

	server := server.New(config.Server, log)
	server.RegisterService(func(s *grpc.Server, m *runtime.ServeMux) {
		authv1.RegisterAuthServiceServer(s, svc)
		healthv1.RegisterHealthServer(s, healthSvc)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err = authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, m, config.Server.GRPCAddr, opts); err != nil {
			log.Error("failed to register server", ilog.Err(err))
			os.Exit(1)
		}
	})
	server.Start()
	if err = server.AwaitShutdown(ctx); err != nil {
		log.Error("server shutdown failed", ilog.Err(err))
	}

	// if err = otel.Shutdown(ctx); err != nil {
	// 	log.Error("otel shutdown failed", "erroor", err.Error())
	// }
}

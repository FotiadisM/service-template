package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	apiauthv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/db"
	"github.com/FotiadisM/mock-microservice/internal/server"
	srvauthv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/health"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
	"github.com/FotiadisM/mock-microservice/pkg/version"
)

//go:generate mockery

func main() {
	version.AddFlag(nil)
	flag.Parse()

	ctx := context.Background()
	config := config.NewConfig(ctx)

	log := ilog.NewLogger()
	slog.SetDefault(log)

	db, err := db.New(config.DB)
	if err != nil {
		log.Error("failed to create db", ilog.Err(err.Error()))
		os.Exit(1)
	}

	srv := srvauthv1.NewService(db)
	healthSrv := health.NewService()

	server := &server.Server{
		Log:    log,
		Config: config,
		ServerRegistrationFunc: func(s grpc.ServiceRegistrar, mux *runtime.ServeMux) error {
			authv1.RegisterAuthServiceServer(s, srv)
			healthv1.RegisterHealthServer(s, healthSrv)

			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			err = apiauthv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.Server.GRPC.Addr, opts)

			return err
		},
	}

	err = server.Start(ctx)
	if err != nil {
		log.Error("failed to start server", ilog.Err(err.Error()))
		os.Exit(1)
	}

	server.AwaitShutdown(ctx)
}

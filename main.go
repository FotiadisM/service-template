package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	server "github.com/FotiadisM/mock-microservice/internal"
	"github.com/FotiadisM/mock-microservice/internal/service"
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

	store, err := store.New(ctx, config.Store)
	if err != nil {
		log.Fatal("failed to create store", zap.Error(err))
	}

	svc := service.NewService(store)

	server := server.New(config.Server, log)
	server.Configure(svc)
	server.RegisterService(func(s *grpc.Server, m *runtime.ServeMux) {
		authv1.RegisterAuthServiceServer(s, svc)
		if err := authv1.RegisterAuthServiceHandlerServer(ctx, m, svc); err != nil {
			log.Fatal("failed to register server", zap.Error(err))
		}
	})
	server.Start(ctx)
}

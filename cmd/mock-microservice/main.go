package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/db"
	"github.com/FotiadisM/mock-microservice/internal/server"
	authv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
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

	if !config.Server.Inst.OtelSDKDisabled {
		var shutdownFunc otelShutDownFunc
		shutdownFunc, err = initializeOTEL(ctx, log, config.Server.Inst.OtelExporterAddr)
		if err != nil {
			log.Error("failed to initialize otel SDK", ilog.Err(err.Error()))
		}
		defer func() {
			err = shutdownFunc(ctx)
			if err != nil {
				log.Error("failed to shutdown otel", ilog.Err(err.Error()))
			}
		}()
	}

	checker := &checker{
		DB: db,
	}
	svc := authv1.NewService(db)

	server, err := server.NewServer(config, log, svc, checker)
	if err != nil {
		log.Error("failed to create server", ilog.Err(err.Error()))
	}

	server.Start()
	server.AwaitShutdown(ctx)
}

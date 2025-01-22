package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"

	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/database"
	"github.com/FotiadisM/mock-microservice/internal/server"
	authv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
	"github.com/FotiadisM/mock-microservice/pkg/version"
)

func main() {
	version.AddFlag(nil)
	flag.Parse()

	ctx := context.Background()
	config := config.NewConfig(ctx)

	log := ilog.NewLogger()
	slog.SetDefault(log)

	db, err := database.New(config.DB)
	if err != nil {
		log.Error("failed to create db", ilog.Err(err))
		os.Exit(1)
	}

	if !config.Server.Inst.OtelSDKDisabled {
		var shutdownFunc otelShutDownFunc
		shutdownFunc, err = initializeOTEL(ctx, log, config.Server.Inst.OtelExporterAddr)
		if err != nil {
			log.Error("failed to initialize otel SDK", ilog.Err(err))
		}
		defer func() {
			err = shutdownFunc(ctx)
			if err != nil {
				log.Error("failed to shutdown otel", ilog.Err(err))
			}
		}()
	}

	checker := &checker{
		DB: db.DB,
	}
	svc := authv1.NewService(db)

	interceptors := server.ChainMiddleware(config, log)
	authsvcPath, authsvcHanlder := authv1connect.NewAuthServiceHandler(svc,
		connect.WithInterceptors(interceptors...),
	)

	services := map[string]http.Handler{
		authsvcPath: authsvcHanlder,
	}

	server, err := server.NewServer(config, log, services, checker)
	if err != nil {
		log.Error("failed to create server", ilog.Err(err))
	}

	server.Start()
	server.AwaitShutdown(ctx)
}

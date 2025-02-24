package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"

	"github.com/FotiadisM/service-template/api/docs"
	"github.com/FotiadisM/service-template/api/gen/go/book/v1/bookv1connect"
	"github.com/FotiadisM/service-template/internal/config"
	"github.com/FotiadisM/service-template/internal/database"
	"github.com/FotiadisM/service-template/internal/server"
	bookv1 "github.com/FotiadisM/service-template/internal/services/book/v1"
	"github.com/FotiadisM/service-template/pkg/ilog"
	"github.com/FotiadisM/service-template/pkg/version"
)

func main() {
	version.AddFlag(nil)
	flag.Parse()

	ctx := context.Background()
	config := config.NewConfig(ctx)

	shutdownFunc, err := initializeOTEL(ctx, config.Inst)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to initialize otel SDK: %v\n", err)
		os.Exit(1)
	}

	log := ilog.NewLogger(
		ilog.WithLogLevel(slog.Level(config.Logging.Level)),
		ilog.WithJSON(config.Logging.JSON),
		ilog.WithAddSource(config.Logging.AddSource),
	)
	slog.SetDefault(log)
	defer func() {
		err = shutdownFunc(ctx)
		if err != nil {
			log.Error("failed to shutdown otel", ilog.Err(err))
		}
	}()

	db, err := database.New(config.DB)
	if err != nil {
		log.Error("failed to create db", ilog.Err(err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	healthChecker := &healthChecker{
		DB: db.DB,
	}
	mux.Handle(grpchealth.NewHandler(healthChecker))
	mux.Handle("/api/docs/", http.StripPrefix("/api/docs/", http.FileServerFS(docs.DocsFS)))

	svc := bookv1.NewService(db)
	interceptors := server.ChainMiddleware(config, log)
	booksvcPath, booksvcHanlder := bookv1connect.NewBookServiceHandler(svc,
		connect.WithInterceptors(interceptors...),
	)

	serverHandler := server.ChainHandlers(mux, config, log, map[string]http.Handler{
		booksvcPath: booksvcHanlder,
	})

	server, err := server.NewServer(config, log, serverHandler)
	if err != nil {
		log.Error("failed to create server", ilog.Err(err))
	}

	server.Start()
	server.AwaitShutdown(ctx)
}

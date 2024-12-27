package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
	"connectrpc.com/vanguard"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/db"
	authv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/logging"
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

	srv := authv1.NewService(db)

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		log.Error("failed to otel interceptor", ilog.Err(err.Error()))
		os.Exit(1)
	}

	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		log.Error("failed to create validate interceptor", ilog.Err(err.Error()))
		os.Exit(1)
	}

	loggingInterceptor := logging.NewInterceptor(log)

	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(srv,
		connect.WithInterceptors(otelInterceptor, loggingInterceptor, validationInterceptor),
	)
	vanguardSvc := vanguard.NewService(svcPath, svcHandler)
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{vanguardSvc})
	if err != nil {
		log.Error("failed to create vanguard transcoder", ilog.Err(err.Error()))
		os.Exit(1)
	}
	mux := http.NewServeMux()
	mux.Handle("/", transcoder)

	reflector := grpcreflect.NewStaticReflector(authv1connect.AuthServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	server := &http.Server{
		Addr:              config.Server.GRPC.Addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:       config.Server.HTTP.ReadTimeout,
		ReadHeaderTimeout: config.Server.HTTP.ReadHeaderTimeout,
		WriteTimeout:      config.Server.HTTP.WriteTimeout,
		IdleTimeout:       config.Server.HTTP.IdleTimeout,
		MaxHeaderBytes:    config.Server.HTTP.MaxHeaderBytes,
	}

	log.Info("http server is listening", "addr", config.Server.GRPC.Addr)
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server exited", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptSignal
	log.Info("shuting down")
	timer, cancel := context.WithTimeout(ctx, config.Server.HTTP.ShutdownTimeout)
	defer cancel()
	err = server.Shutdown(timer)
	if err != nil {
		log.Error("failed to shutdown http server", ilog.Err(err.Error()))
	}
}

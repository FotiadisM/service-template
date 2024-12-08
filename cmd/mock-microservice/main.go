package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	apiauthv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/config"
	svcauthv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/health"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/logging"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/recovery"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/validate"
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
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Error("open-telemtry error", ilog.Err(err.Error()))
	}))

	store, err := store.New(config.DB)
	if err != nil {
		log.Error("failed to create store", ilog.Err(err.Error()))
		os.Exit(1)
	}

	svc := svcauthv1.NewService(store)
	healthSvc := health.NewService()

	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	usi := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(log,
			logging.WithFilter(filter.Any(filter.Reflection(), filter.HealthCheck())),
		),
		recovery.UnaryServerInterceptor(),
		validate.UnaryServerInterceptor(validator),
	}

	ssi := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(log, logging.WithFilter(filter.Any(filter.Reflection(), filter.HealthCheck()))),
		recovery.StreamServerInterceptor(),
		validate.StreamServerInterceptor(validator),
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(usi...),
		grpc.ChainStreamInterceptor(ssi...),
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	if config.Server.GRPC.Reflection {
		log.Info("enabling grpc reflection")
		reflection.Register(grpcServer)
	}

	mux := runtime.NewServeMux()
	httpServer := &http.Server{
		Addr:              config.Server.HTTP.Addr,
		Handler:           mux,
		ReadTimeout:       config.Server.HTTP.ReadTimeout,
		ReadHeaderTimeout: config.Server.HTTP.ReadHeaderTimeout,
		WriteTimeout:      config.Server.HTTP.WriteTimeout,
		IdleTimeout:       config.Server.HTTP.IdleTimeout,
		MaxHeaderBytes:    config.Server.HTTP.MaxHeaderBytes,
	}

	apiauthv1.RegisterAuthServiceServer(grpcServer, svc)
	healthv1.RegisterHealthServer(grpcServer, healthSvc)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = apiauthv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.Server.GRPC.Addr, opts); err != nil {
		log.Error("failed to register server", ilog.Err(err.Error()))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", config.Server.GRPC.Addr)
	if err != nil {
		log.Error("failed to create net.Listener", ilog.Err(err.Error()))
		os.Exit(1)
	}

	log.Info("grpc server is listening", "port", config.Server.GRPC.Addr)
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Error("grpc serve failed", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	log.Info("http server is listening", "port", config.Server.HTTP.Addr)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http listen and serve failed", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptSignal
	log.Info("shuting down")
	timer, cancel := context.WithTimeout(context.Background(), config.Server.HTTP.ShutdownTimeout)
	defer cancel()
	err = httpServer.Shutdown(timer)
	if err != nil {
		log.Error("failed to shutdown http server", ilog.Err(err.Error()))
	}
	grpcServer.GracefulStop()
}

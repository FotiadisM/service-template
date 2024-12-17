package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	apiauthv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/db"
	srvauthv1 "github.com/FotiadisM/mock-microservice/internal/services/auth/v1"
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
	slog.SetDefault(log)

	db, err := db.New(config.DB)
	if err != nil {
		log.Error("failed to create db", ilog.Err(err.Error()))
		os.Exit(1)
	}

	srv := srvauthv1.NewService(db)
	healthSrv := health.NewService()

	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	if !config.Server.OTEL.SDKDisabled {
		var otelShutDown otelShutDownFunc
		otelShutDown, err = initializeOTEL(ctx, log, config.Server.OTEL.ExporterAddr)
		if err != nil {
			log.Error("failed to initialize otel sdk", ilog.Err(err.Error()))
			os.Exit(1)
		}
		defer func() {
			err = otelShutDown(context.Background())
			if err != nil {
				log.Error("failed to shutdown otel exporter", ilog.Err(err.Error()))
			}
		}()
		log.Info("enabled otel instrumentation")
	}

	otelStatsHandler := otelgrpc.NewServerHandler(otelgrpc.WithFilter(otelgrpcFilter))

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
		grpc.StatsHandler(otelStatsHandler),
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	if config.Server.GRPC.Reflection {
		reflection.Register(grpcServer)
		log.Info("enabled grpc reflection")
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

	apiauthv1.RegisterAuthServiceServer(grpcServer, srv)
	healthv1.RegisterHealthServer(grpcServer, healthSrv)
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
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Error("gRPC server exited", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	log.Info("http server is listening", "port", config.Server.HTTP.Addr)
	go func() {
		err = httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server exited", ilog.Err(err.Error()))
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

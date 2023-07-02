package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/logger"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/otelgrpc"
)

type Config struct {
	GRPCAddr   string `env:"GRPC_ADDR,default=:8080"`
	HTTPAddr   string `env:"HTTP_ADDR,default=:9090"`
	Debug      bool   `env:"DEBUG"`
	Reflection bool   `env:"REFLECTION,default=$DEBUG"`
}

type Server struct {
	config Config

	log *slog.Logger

	grpcServer *grpc.Server
	httpServer *http.Server
	mux        *runtime.ServeMux
}

func New(config Config, log *slog.Logger) *Server {
	return &Server{
		config: config,
		log:    log,
	}
}

func (s *Server) Configure() {
	// recoveryFunc := recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) error {
	// 	log := logger.FromContext(ctx)
	// 	log.LogAttrs(ctx, slog.LevelError, "PANIC", slog.Any("trace", p))
	// 	return status.Error(codes.Internal, "internal server error")
	// })

	usi := []grpc.UnaryServerInterceptor{
		logger.UnaryServerInterceptor(s.log),
		recovery.UnaryServerInterceptor(),
	}

	ssi := []grpc.StreamServerInterceptor{
		logger.StreamServerInterceptor(s.log),
		recovery.StreamServerInterceptor(),
		// logger.StreamServerInterceptor(s.log),
	}

	handler := otelgrpc.ServerStatsHandler(
		otelgrpc.WithFilter(
			filter.Any(
				filter.HealthCheck(),
				filter.ServiceName("grpc.reflection.v1alpha.ServerReflection"),
			),
		),
	)
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(usi...),
		grpc.ChainStreamInterceptor(ssi...),
		grpc.StatsHandler(handler),
	}
	s.grpcServer = grpc.NewServer(grpcOpts...)

	if s.config.Reflection {
		reflection.Register(s.grpcServer)
	}

	s.mux = runtime.NewServeMux()

	s.httpServer = &http.Server{
		Addr:              s.config.HTTPAddr,
		Handler:           s.mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (s *Server) RegisterService(rgFn func(s *grpc.Server, m *runtime.ServeMux)) {
	rgFn(s.grpcServer, s.mux)
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.config.GRPCAddr)
	if err != nil {
		s.log.Error("failed to create net.Listener", "err", err.Error())
		os.Exit(1)
	}

	s.log.Info("grpc server is listening", "port", s.config.GRPCAddr)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.log.Error("grpc serve failed", "err", err.Error())
			os.Exit(1)
		}
	}()

	s.log.Info("http server is listening", "port", s.config.HTTPAddr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("http listen and serve failed", "err", err.Error())
				os.Exit(1)
			}
		}
	}()
}

func (s *Server) AwaitShutdown(ctx context.Context) error {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal

	errs := make(chan error, 2)
	go func() {
		errs <- s.httpServer.Shutdown(ctx)
	}()
	go func() {
		s.grpcServer.GracefulStop()
		errs <- nil
	}()

	return errors.Join(<-errs, <-errs)
}

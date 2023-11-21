package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/logging"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/recovery"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/otelgrpc"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

var (
	// DefaultReadTimeout sets the maximum time a client has to fully stream a request (5s).
	DefaultReadTimeout = 5 * time.Second
	// DefaultWriteTimeout sets the maximum amount of time a handler has to fully process a request (10s).
	DefaultWriteTimeout = 10 * time.Second
	// DefaultIdleTimeout sets the maximum amount of time a Keep-Alive connection can remain idle before
	// being recycled (120s).
	DefaultIdleTimeout = 120 * time.Second
	// DefaultReadHeaderTimeout sets the maximum amount of time a client has to fully stream a request header (5s).
	DefaultReadHeaderTimeout = DefaultReadTimeout
	// DefaultShutdownTimeout defines how long Graceful will wait before forcibly shutting down.
	DefaultShutdownTimeout = 5 * time.Second
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
	s := &Server{
		config: config,
		log:    log,
	}

	usi := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(s.log),
		recovery.UnaryServerInterceptor(),
	}

	ssi := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(s.log),
		recovery.StreamServerInterceptor(),
	}

	handler := otelgrpc.ServerStatsHandler(
		otelgrpc.WithFilter(
			filter.Any(
				filter.Reflection(),
				filter.HealthCheck(),
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
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
	}

	return s
}

func (s *Server) RegisterService(rgFn func(s *grpc.Server, m *runtime.ServeMux)) {
	rgFn(s.grpcServer, s.mux)
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.config.GRPCAddr)
	if err != nil {
		s.log.Error("failed to create net.Listener", ilog.Err(err))
		os.Exit(1)
	}

	s.log.Info("grpc server is listening", "port", s.config.GRPCAddr)
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			s.log.Error("grpc serve failed", ilog.Err(err))
			os.Exit(1)
		}
	}()

	s.log.Info("http server is listening", "port", s.config.HTTPAddr)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http listen and serve failed", ilog.Err(err))
			os.Exit(1)
		}
	}()
}

func (s *Server) GracefulStop() error {
	return s.GracefulStopContext(context.Background())
}

func (s *Server) GracefulStopContext(ctx context.Context) error {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-interruptSignal:
	}

	errs := make(chan error, 2)
	go func() {
		timer, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()
		errs <- s.httpServer.Shutdown(timer)
	}()
	go func() {
		s.grpcServer.GracefulStop()
		errs <- nil
	}()

	return errors.Join(<-errs, <-errs)
}

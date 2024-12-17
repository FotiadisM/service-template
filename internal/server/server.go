package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServiceRegistrationFunc func(s grpc.ServiceRegistrar, mux *runtime.ServeMux) error

type Server struct {
	Log    *slog.Logger
	Config *config.Config

	ServerRegistrationFunc ServiceRegistrationFunc

	grpcServer *grpc.Server
	httpServer *http.Server

	otelShutDownFunc otelShutDownFunc
}

func (s *Server) Start(ctx context.Context) error {
	grpcOpts, err := s.newServerMiddleware(ctx)
	if err != nil {
		return fmt.Errorf("failed to create grpc server options: %w", err)
	}

	s.grpcServer = grpc.NewServer(grpcOpts...)
	if s.Config.Server.GRPC.Reflection {
		reflection.Register(s.grpcServer)
		s.Log.Info("enabled grpc reflection")
	}

	mux := runtime.NewServeMux()
	s.httpServer = &http.Server{
		Addr:              s.Config.Server.HTTP.Addr,
		Handler:           mux,
		ReadTimeout:       s.Config.Server.HTTP.ReadTimeout,
		ReadHeaderTimeout: s.Config.Server.HTTP.ReadHeaderTimeout,
		WriteTimeout:      s.Config.Server.HTTP.WriteTimeout,
		IdleTimeout:       s.Config.Server.HTTP.IdleTimeout,
		MaxHeaderBytes:    s.Config.Server.HTTP.MaxHeaderBytes,
	}

	err = s.ServerRegistrationFunc(s.grpcServer, mux)
	if err != nil {
		return fmt.Errorf("failed to register services %w", err)
	}

	lis, err := net.Listen("tcp", s.Config.Server.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("failed to create net.Listener: %w", err)
	}

	s.Log.Info("grpc server is listening", "port", s.Config.Server.GRPC.Addr)
	go func() {
		err = s.grpcServer.Serve(lis)
		if err != nil {
			s.Log.Error("gRPC server exited", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	s.Log.Info("http server is listening", "port", s.Config.Server.HTTP.Addr)
	go func() {
		err = s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error("http server exited", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()

	return nil
}

func (s *Server) AwaitShutdown(ctx context.Context) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptSignal
	s.Log.Info("shuting down")
	timer, cancel := context.WithTimeout(ctx, s.Config.Server.HTTP.ShutdownTimeout)
	defer cancel()
	err := s.httpServer.Shutdown(timer)
	if err != nil {
		s.Log.Error("failed to shutdown http server", ilog.Err(err.Error()))
	}
	s.grpcServer.GracefulStop()

	err = s.otelShutDownFunc(ctx)
	if err != nil {
		s.Log.Error("failed to shutdown otel exporter", ilog.Err(err.Error()))
	}
}

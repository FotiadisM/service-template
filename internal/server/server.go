package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/vanguard"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type Server struct {
	log *slog.Logger

	config *config.Config
	server *http.Server
}

func NewServer(config *config.Config, log *slog.Logger, svc authv1connect.AuthServiceHandler, checker grpchealth.Checker) (*Server, error) {
	interceptors, err := createInterceptors(log)
	if err != nil {
		return nil, fmt.Errorf("failed to create interceptors: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(grpchealth.NewHandler(checker))

	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(svc,
		connect.WithInterceptors(interceptors...),
	)

	if config.Server.HTTP.DisableRESTTranscoding {
		mux.Handle(svcPath, svcHandler)
	} else {
		vanguardSvc := vanguard.NewService(svcPath, svcHandler)
		transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{vanguardSvc})
		if err != nil {
			return nil, fmt.Errorf("failed to create vanguard transcoder: %w", err)
		}
		mux.Handle("/", transcoder)
		log.Info("enabled http rest transcoding")
	}

	if config.Server.HTTP.Reflection {
		reflector := grpcreflect.NewStaticReflector(authv1connect.AuthServiceName)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
		log.Info("enabled server reflection")
	}

	httpServer := &http.Server{
		Addr:              config.Server.HTTP.Addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:       config.Server.HTTP.ReadTimeout,
		ReadHeaderTimeout: config.Server.HTTP.ReadHeaderTimeout,
		WriteTimeout:      config.Server.HTTP.WriteTimeout,
		IdleTimeout:       config.Server.HTTP.IdleTimeout,
		MaxHeaderBytes:    config.Server.HTTP.MaxHeaderBytes,
	}

	server := &Server{
		log:    log,
		config: config,
		server: httpServer,
	}

	return server, nil
}

func (s *Server) Start() {
	s.log.Info("http server is listening", "addr", s.config.Server.HTTP.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server exited", ilog.Err(err.Error()))
			os.Exit(1)
		}
	}()
}

func (s *Server) AwaitShutdown(ctx context.Context) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptSignal
	s.log.Info("shuting down")
	timer, cancel := context.WithTimeout(ctx, s.config.Server.HTTP.ShutdownTimeout)
	defer cancel()
	err := s.server.Shutdown(timer)
	if err != nil {
		s.log.Error("failed to shutdown http server", ilog.Err(err.Error()))
	}
}

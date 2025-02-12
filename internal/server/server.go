package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/FotiadisM/service-template/internal/config"
	"github.com/FotiadisM/service-template/pkg/ilog"
)

type Server struct {
	log *slog.Logger

	config *config.Config
	server *http.Server
}

func NewServer(config *config.Config, log *slog.Logger, mux http.Handler) (*Server, error) {
	httpServer := &http.Server{
		Addr:              config.Server.Addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:       config.Server.ReadTimeout,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
		MaxHeaderBytes:    config.Server.MaxHeaderBytes,
	}

	server := &Server{
		log:    log,
		config: config,
		server: httpServer,
	}

	return server, nil
}

func (s *Server) Start() {
	s.log.Info("http server is listening", "addr", s.config.Server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server exited", ilog.Err(err))
			os.Exit(1)
		}
	}()
}

func (s *Server) AwaitShutdown(ctx context.Context) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptSignal
	s.log.Info("shuting down")
	timer, cancel := context.WithTimeout(ctx, s.config.Server.ShutdownTimeout)
	defer cancel()
	err := s.server.Shutdown(timer)
	if err != nil {
		s.log.Error("failed to shutdown http server", ilog.Err(err))
	}
}

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

	// grpczap "github.com/grpc-ecosystem/go-grpc-middleware/providers/zap/v2"
	// "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/FotiadisM/mock-microservice/pkg/logger"
	"github.com/FotiadisM/mock-microservice/pkg/otelgrpc"
	"github.com/FotiadisM/mock-microservice/pkg/otelgrpc/filters"
)

type Config struct {
	GRPCAddr   string `env:"GRPC_ADDR,default=:8080"`
	HTTPAddr   string `env:"HTTP_ADDR,default=:9090"`
	Debug      bool   `env:"DEBUG"`
	Reflection bool   `env:"REFLECTION,default=$DEBUG"`
}

type Server struct {
	config Config

	log *zap.Logger

	grpcServer *grpc.Server
	httpServer *http.Server
	mux        *runtime.ServeMux
}

func New(config Config, log *zap.Logger) *Server {
	return &Server{
		config: config,
		log:    log,
	}
}

func (s *Server) Configure() {
	// loggingOpts := []logging.Option{
	// 	logging.WithDecider(func(_ string, _ error) logging.Decision { return logging.LogFinishCall }),
	// }

	recoveryFunc := recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) error {
		log := logger.FromContext(ctx)
		log.Error("application paniced", zap.Any("trace", p))
		return status.Error(codes.Internal, "internal server error")
	})

	// error
	// log
	// recover

	usi := []grpc.UnaryServerInterceptor{
		// logging.UnaryServerInterceptor(grpczap.InterceptorLogger(s.log), loggingOpts...),
		recovery.UnaryServerInterceptor(recoveryFunc),
		logger.UnaryServerInterceptor(s.log),
	}

	ssi := []grpc.StreamServerInterceptor{
		// logging.StreamServerInterceptor(grpczap.InterceptorLogger(s.log), loggingOpts...),
		recovery.StreamServerInterceptor(recoveryFunc),
		logger.StreamServerInterceptor(s.log),
	}

	handler := otelgrpc.ServerStatsHandler(
		otelgrpc.WithFilter(
			filters.Any(
				filters.HealthCheck(),
				filters.ServiceName("grpc.reflection.v1alpha.ServerReflection"),
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
		s.log.Fatal("failed to create net.Listener", zap.Error(err))
	}

	s.log.Info("grpc server is listening", zap.String("port", s.config.GRPCAddr))
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.log.Fatal("grpc serve failed", zap.Error(err))
		}
	}()

	s.log.Info("http server is listening", zap.String("port", s.config.HTTPAddr))
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.log.Fatal("http listen and serve failed", zap.Error(err))
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

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

	grpcopenmetrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/openmetrics/v2"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/providers/zap/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/FotiadisM/mock-microservice/pkg/health"
	"github.com/FotiadisM/mock-microservice/pkg/logger"
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

	grpcServer    *grpc.Server
	httpServer    *http.Server
	mux           *runtime.ServeMux
	serverMetrics *grpcopenmetrics.ServerMetrics
}

func New(config Config, log *zap.Logger) *Server {
	return &Server{
		config: config,
		log:    log,
	}
}

func (s *Server) Configure(svc healthv1.HealthServer) {
	loggingOpts := []logging.Option{
		logging.WithDecider(func(_ string, _ error) logging.Decision { return logging.LogFinishCall }),
	}

	recoveryFunc := recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) error {
		log := logger.FromContext(ctx)
		log.Error("application paniced", zap.Any("trace", p))
		return status.Error(codes.Internal, "internal server error")
	})

	s.serverMetrics = grpcopenmetrics.NewServerMetrics()

	usi := []grpc.UnaryServerInterceptor{
		grpcopenmetrics.UnaryServerInterceptor(s.serverMetrics),
		logging.UnaryServerInterceptor(grpczap.InterceptorLogger(s.log), loggingOpts...),
		recovery.UnaryServerInterceptor(recoveryFunc),
		validator.UnaryServerInterceptor(false),
		logger.UnaryServerInterceptor(s.log),
	}

	ssi := []grpc.StreamServerInterceptor{
		grpcopenmetrics.StreamServerInterceptor(s.serverMetrics),
		logging.StreamServerInterceptor(grpczap.InterceptorLogger(s.log), loggingOpts...),
		recovery.StreamServerInterceptor(recoveryFunc),
		validator.StreamServerInterceptor(false),
		logger.StreamServerInterceptor(s.log),
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(usi...),
		grpc.ChainStreamInterceptor(ssi...),
	}
	s.grpcServer = grpc.NewServer(grpcOpts...)

	if s.config.Reflection {
		reflection.Register(s.grpcServer)
	}

	healthv1.RegisterHealthServer(s.grpcServer, svc)

	muxOptions := []runtime.ServeMuxOption{
		runtime.WithHealthzEndpoint(health.NewHealthClient(svc)),
	}

	s.mux = runtime.NewServeMux(muxOptions...)
	if err := s.mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		s.log.Fatal("mux.HandlePath() failed", zap.Error(err))
	}

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

func (s *Server) Start(ctx context.Context) {
	lis, err := net.Listen("tcp", s.config.GRPCAddr)
	if err != nil {
		s.log.Fatal("failed to create net.Listener", zap.Error(err))
	}

	// initialize metrics after all services have been registered
	s.serverMetrics.InitializeMetrics(s.grpcServer)

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

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal

	s.httpServer.Shutdown(ctx) //nolint:errcheck
	s.grpcServer.GracefulStop()
	lis.Close() //nolint:errcheck
}

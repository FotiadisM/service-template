package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/filter"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/logging"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/recovery"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/validate"
)

func otelgrpcFilter(ri *stats.RPCTagInfo) bool {
	fullName := strings.TrimLeft(ri.FullMethodName, "/")
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return true
	}
	service := parts[0]

	switch service {
	case "grpc.reflection.v1.ServerReflection":
		return false
	case "grpc.reflection.v1alpha.ServerReflection":
		return false
	case "grpc.health.v1.Health":
		return false
	}

	return true
}

func (s *Server) newServerMiddleware(ctx context.Context) ([]grpc.ServerOption, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create proto validattor: %w", err)
	}

	if !s.Config.Server.OTEL.SDKDisabled {
		s.otelShutDownFunc, err = initializeOTEL(ctx, s.Log, s.Config.Server.OTEL.ExporterAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize otel sdk: %w", err)
		}
		s.Log.Info("enabled otel instrumentation")
	} else {
		s.otelShutDownFunc = func(_ context.Context) error { return nil }
	}

	otelStatsHandler := otelgrpc.NewServerHandler(otelgrpc.WithFilter(otelgrpcFilter))

	usi := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(s.Log,
			logging.WithFilter(filter.Any(filter.Reflection(), filter.HealthCheck())),
		),
		recovery.UnaryServerInterceptor(),
		validate.UnaryServerInterceptor(validator),
	}

	ssi := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(s.Log, logging.WithFilter(filter.Any(filter.Reflection(), filter.HealthCheck()))), //nolint:contextcheck false positive
		recovery.StreamServerInterceptor(), //nolint:contextcheck false positive
		validate.StreamServerInterceptor(validator),
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(usi...),
		grpc.ChainStreamInterceptor(ssi...),
		grpc.StatsHandler(otelStatsHandler),
	}

	return grpcOpts, nil
}

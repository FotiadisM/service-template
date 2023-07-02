package health

import (
	"context"

	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type ProbeFunc func(context.Context) (healthv1.HealthCheckResponse_ServingStatus, error)

func DefaultProbeFunc(_ context.Context) (healthv1.HealthCheckResponse_ServingStatus, error) {
	return healthv1.HealthCheckResponse_SERVING, nil
}

type service struct {
	readiness ProbeFunc
	liveness  ProbeFunc
	startup   ProbeFunc

	healthv1.UnimplementedHealthServer
}

func NewService(readiness, liveness, startup ProbeFunc) healthv1.HealthServer {
	if readiness == nil {
		readiness = DefaultProbeFunc
	}
	if liveness == nil {
		liveness = DefaultProbeFunc
	}
	if startup == nil {
		startup = DefaultProbeFunc
	}
	return &service{
		readiness: readiness,
		liveness:  liveness,
		startup:   startup,
	}
}

func (s *service) Check(ctx context.Context, in *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	var f ProbeFunc
	switch in.Service {
	case "readiness":
		f = s.readiness
	case "liveness":
		f = s.liveness
	case "startup":
		f = s.startup
	default:
		return nil, status.Error(codes.NotFound, "unknown probe")
	}

	st, err := f(ctx)
	if err != nil {
		return nil, err
	}

	return &healthv1.HealthCheckResponse{Status: st}, nil
}

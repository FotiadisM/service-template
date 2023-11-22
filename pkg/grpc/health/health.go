package health

import (
	"context"

	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type service struct {
	readiness ProbeFunc
	liveness  ProbeFunc
	startup   ProbeFunc

	healthv1.UnimplementedHealthServer
}

func NewService(opts ...Option) healthv1.HealthServer {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &service{
		readiness: options.readiness,
		liveness:  options.liveness,
		startup:   options.startup,
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

func (s *service) Watch(in *healthv1.HealthCheckRequest, ws healthv1.Health_WatchServer) error {
	ctx := ws.Context()

	var f ProbeFunc
	switch in.Service {
	case "readiness":
		f = s.readiness
	case "liveness":
		f = s.liveness
	case "startup":
		f = s.startup
	default:
		return status.Error(codes.NotFound, "unknown probe")
	}

	st, err := f(ctx)
	if err != nil {
		return err
	}

	return ws.Send(&healthv1.HealthCheckResponse{Status: st})
}

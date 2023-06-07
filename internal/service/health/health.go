package health

import (
	"context"

	"github.com/FotiadisM/mock-microservice/internal/store"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type service struct {
	store store.Store

	healthv1.UnimplementedHealthServer
}

func NewService(st store.Store) healthv1.HealthServer {
	return &service{store: st}
}

func (s *service) Check(_ context.Context, _ *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	res := &healthv1.HealthCheckResponse{
		Status: healthv1.HealthCheckResponse_SERVING,
	}

	return res, nil
}

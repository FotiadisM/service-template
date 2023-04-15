package servicev1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	store store.Store

	authv1.UnimplementedAuthServiceServer
	health.UnimplementedHealthServer
}

func NewService(st store.Store) *Service {
	return &Service{store: st}
}

func (s *Service) Check(_ context.Context, _ *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	res := &health.HealthCheckResponse{}
	res.Status = health.HealthCheckResponse_SERVING
	return res, nil
}

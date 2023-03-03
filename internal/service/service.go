package service

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	DB store.Store

	authv1.UnimplementedAuthServiceServer
	health.UnimplementedHealthServer
}

func NewService(DB store.Store) *Service {
	return &Service{DB: DB}
}

func (s *Service) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	res := &health.HealthCheckResponse{}
	res.Status = health.HealthCheckResponse_SERVING
	return res, nil
}

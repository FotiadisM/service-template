package authv1

import (
	"github.com/FotiadisM/mock-microservice/internal/db"
)

type Service struct {
	db db.DB

	// health.UnimplementedHealthServer
}

func NewService(db db.DB) *Service {
	return &Service{db: db}
}

// func (s *Service) Check(_ context.Context, _ *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
// 	res := &health.HealthCheckResponse{}
// 	res.Status = health.HealthCheckResponse_SERVING
// 	return res, nil
// }

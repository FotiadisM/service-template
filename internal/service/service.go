package service

import (
	"context"
	"database/sql"
	"fmt"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var ErrStupidError = NewError(codes.Internal, "WHATS_UP", "what's up my boy")

func NewError(code codes.Code, reason, msg string) error {
	return NewErrorWithDomain(code, reason, msg, "auth-svc")
}

func NewErrorWithDomain(code codes.Code, reason, msg, domain string) error {
	st := status.New(code, msg)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: reason,
		Domain: domain,
	})
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return st.Err()
}

type Service struct {
	DB *sql.DB

	authv1.UnimplementedAuthServiceServer
	health.UnimplementedHealthServer
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	res := &health.HealthCheckResponse{}
	res.Status = health.HealthCheckResponse_SERVING
	return res, nil
}

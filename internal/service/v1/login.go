package servicev1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/errors"
	"google.golang.org/grpc/codes"
)

func (s *Service) Login(_ context.Context, _ *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	err := errors.NewDetailsError(codes.Internal, "MY_CUSTOM-CODE", "Unexpected error")

	return nil, err
}

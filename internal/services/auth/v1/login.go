package authv1

import (
	"context"

	"google.golang.org/grpc/codes"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/errors"
)

func (s *Service) Login(_ context.Context, _ *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	err := errors.NewInfoError(codes.Internal, "MY_CUSTOM-CODE", "Unexpected error", map[string]string{
		"one":   "two",
		"three": "four",
	})

	return nil, err
}

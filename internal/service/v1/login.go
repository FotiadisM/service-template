package servicev1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
)

func (s *Service) Login(_ context.Context, _ *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("not implemented") // TODO: Implement
}

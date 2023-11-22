package authv1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/go/auth/v1"
)

func (s *Service) RefreshToken(_ context.Context, _ *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	panic("not implemented") // TODO: Implement
}

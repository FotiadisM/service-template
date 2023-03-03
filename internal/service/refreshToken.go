package service

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
)

func (s *Service) RefreshToken(ctx context.Context, in *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	panic("not implemented") // TODO: Implement
}

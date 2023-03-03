package service

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
)

func (s *Service) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	out := &authv1.RegisterResponse{
		AccessToken:  "1234",
		RefreshToken: "5678",
	}
	return out, nil
}

package service

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
)

func (s *Service) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("not implemented") // TODO: Implement
}

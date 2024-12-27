package authv1

import (
	"context"

	"connectrpc.com/connect"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
)

func (s *Service) RefreshToken(_ context.Context, _ *connect.Request[authv1.RefreshTokenRequest]) (*connect.Response[authv1.RefreshTokenResponse], error) {
	panic("not implemented") // TODO: Implement
}

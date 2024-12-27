package authv1

import (
	"context"

	"connectrpc.com/connect"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

func (s *Service) Logout(ctx context.Context, _ *connect.Request[authv1.LogoutRequest]) (*connect.Response[authv1.LogoutResponse], error) {
	log := ilog.FromContext(ctx)
	log.Info("HELLO")

	panic("oh no this is terrible")
}

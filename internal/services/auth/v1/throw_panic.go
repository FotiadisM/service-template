package authv1

import (
	"context"

	"connectrpc.com/connect"
	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
)

func (s *Service) ThrowPanic(_ context.Context, _ *connect.Request[authv1.ThrowPanicRequest]) (*connect.Response[authv1.ThrowPanicResponse], error) {
	panic("this is a ranom panic")
}

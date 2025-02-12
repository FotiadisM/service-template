package bookv1

import (
	"context"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *Service) ThrowPanic(_ context.Context, _ *connect.Request[bookv1.ThrowPanicRequest]) (*connect.Response[bookv1.ThrowPanicResponse], error) {
	panic("this is a ranom panic")
}

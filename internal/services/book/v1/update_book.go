package bookv1

import (
	"context"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *Service) UpdateBook(_ context.Context, _ *connect.Request[bookv1.UpdateBookRequest]) (*connect.Response[bookv1.UpdateBookResponse], error) {
	panic("not implemented") // TODO: Implement
}

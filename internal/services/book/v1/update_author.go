package bookv1

import (
	"context"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/mock-microservice/api/gen/go/book/v1"
)

func (s *Service) UpdateAuthor(_ context.Context, _ *connect.Request[bookv1.UpdateAuthorRequest]) (*connect.Response[bookv1.UpdateAuthorResponse], error) {
	panic("not implemented") // TODO: Implement
}

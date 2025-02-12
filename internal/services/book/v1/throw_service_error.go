package bookv1

import (
	"context"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/errors"
)

func (s *Service) ThrowServiceError(_ context.Context, _ *connect.Request[bookv1.ThrowServiceErrorRequest]) (*connect.Response[bookv1.ThrowServiceErrorResponse], error) {
	return nil, errors.ErrMyError
}

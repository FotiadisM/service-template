package bookv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/google/uuid"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *Service) DeleteAuthor(ctx context.Context, req *connect.Request[bookv1.DeleteAuthorRequest]) (*connect.Response[bookv1.DeleteAuthorResponse], error) {
	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to parse author id: %w", err))
	}

	err = s.db.DeleteAuthor(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to delete author: %w", err)
	}

	res := connect.NewResponse(&bookv1.DeleteAuthorResponse{})
	return res, nil
}

package bookv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/google/uuid"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/encoder"
)

func (s *Service) GetAuthor(ctx context.Context, req *connect.Request[bookv1.GetAuthorRequest]) (*connect.Response[bookv1.GetAuthorResponse], error) {
	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to parse author id: %w", err))
	}

	author, err := s.db.GetAuthor(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	res := connect.NewResponse(&bookv1.GetAuthorResponse{
		Author: encoder.DBAuthorToAPI(author),
	})

	return res, nil
}

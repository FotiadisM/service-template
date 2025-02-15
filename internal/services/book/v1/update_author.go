package bookv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/encoder"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *Service) UpdateAuthor(ctx context.Context, req *connect.Request[bookv1.UpdateAuthorRequest]) (*connect.Response[bookv1.UpdateAuthorResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to parse book id: %w", err))
	}

	updateParams := queries.UpdateAuthorParams{ID: id, UpdatedAt: time.Now()}
	if req.Msg.Name != nil {
		updateParams.Name = sql.NullString{String: *req.Msg.Name, Valid: true}
	}
	if req.Msg.Bio != nil {
		updateParams.Bio = sql.NullString{String: *req.Msg.Bio, Valid: true}
	}

	author, err := s.db.UpdateAuthor(ctx, updateParams)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("author not found"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	res := connect.NewResponse(&bookv1.UpdateAuthorResponse{
		Author: encoder.DBAuthorToAPI(author),
	})

	return res, nil
}

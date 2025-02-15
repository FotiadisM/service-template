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

func (s *Service) UpdateBook(ctx context.Context, req *connect.Request[bookv1.UpdateBookRequest]) (*connect.Response[bookv1.UpdateBookResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to parse book id: %w", err))
	}

	updateParams := queries.UpdateBookParams{ID: id, UpdatedAt: time.Now()}
	if req.Msg.Title != nil {
		updateParams.Title = sql.NullString{String: *req.Msg.Title, Valid: true}
	}
	if req.Msg.Description != nil {
		updateParams.Description = sql.NullString{String: *req.Msg.Description, Valid: true}
	}

	book, err := s.db.UpdateBook(ctx, updateParams)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("book not found"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	res := connect.NewResponse(&bookv1.UpdateBookResponse{
		Book: encoder.DBBookToAPI(book),
	})

	return res, nil
}

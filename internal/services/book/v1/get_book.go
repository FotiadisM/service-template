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

func (s *Service) GetBook(ctx context.Context, req *connect.Request[bookv1.GetBookRequest]) (*connect.Response[bookv1.GetBookResponse], error) {
	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to parse book id: %w", err))
	}

	book, err := s.db.GetBook(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	res := connect.NewResponse(&bookv1.GetBookResponse{
		Book: encoder.DBBookToAPI(book),
	})

	return res, nil
}

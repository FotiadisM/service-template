package bookv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/mock-microservice/api/gen/go/book/v1"
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
		Book: &bookv1.Book{
			Id:        req.Msg.GetId(),
			Title:     book.Title,
			AuthorId:  book.AuthorID.String(),
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		},
	})

	return res, nil
}

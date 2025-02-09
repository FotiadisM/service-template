package bookv1

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/mock-microservice/api/gen/go/book/v1"
	"github.com/FotiadisM/mock-microservice/internal/services/book/v1/queries"
)

func (s *Service) CreateBook(ctx context.Context, req *connect.Request[bookv1.CreateBookRequest]) (*connect.Response[bookv1.CreateBookResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid %w", err)
	}

	now := time.Now()
	authorID, err := uuid.Parse(req.Msg.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse author id %w", err)
	}
	book := queries.CreateBookParams{
		ID:        id,
		Title:     req.Msg.Title,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = s.db.CreateBook(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("failed to create author %w", err)
	}

	res := connect.NewResponse(&bookv1.CreateBookResponse{
		Book: &bookv1.Book{
			Id:        book.ID.String(),
			Title:     book.Title,
			AuthorId:  book.AuthorID.String(),
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		},
	})

	return res, nil
}

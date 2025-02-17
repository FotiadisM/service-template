package bookv1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/google/uuid"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/encoder"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *Service) CreateBookReview(ctx context.Context, req *connect.Request[bookv1.CreateBookReviewRequest]) (*connect.Response[bookv1.CreateBookReviewResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to created uuid: %w", err)
	}
	bookID, err := uuid.Parse(req.Msg.BookId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	createParams := queries.CreateBookReviewParams{
		ID:     id,
		BookID: bookID,
		Rating: req.Msg.Rating,
		Text:   req.Msg.Text,
	}
	review, err := s.db.CreateBookReview(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create book review: %w", err)
	}

	res := connect.NewResponse(&bookv1.CreateBookReviewResponse{
		Review: encoder.DBBookReviewToAPI(review),
	})

	return res, nil
}

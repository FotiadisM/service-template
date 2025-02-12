package bookv1

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *Service) CreateAuthor(ctx context.Context, req *connect.Request[bookv1.CreateAuthorRequest]) (*connect.Response[bookv1.CreateAuthorResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid %w", err)
	}

	now := time.Now().UTC()
	author := queries.CreateAuthorParams{
		ID:        id,
		Name:      req.Msg.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = s.db.CreateAuthor(ctx, author)
	if err != nil {
		return nil, fmt.Errorf("failed to create author %w", err)
	}

	res := connect.NewResponse(&bookv1.CreateAuthorResponse{
		Author: &bookv1.Author{
			Id:        author.ID.String(),
			Name:      author.Name,
			CreatedAt: timestamppb.New(author.CreatedAt),
			UpdatedAt: timestamppb.New(author.UpdatedAt),
		},
	})

	return res, nil
}

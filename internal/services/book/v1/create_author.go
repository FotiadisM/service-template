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

func (s *Service) CreateAuthor(ctx context.Context, req *connect.Request[bookv1.CreateAuthorRequest]) (*connect.Response[bookv1.CreateAuthorResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid %w", err)
	}

	createParams := queries.CreateAuthorParams{
		ID:   id,
		Name: req.Msg.Name,
		Bio:  req.Msg.Bio,
	}
	author, err := s.db.CreateAuthor(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create author %w", err)
	}

	res := connect.NewResponse(&bookv1.CreateAuthorResponse{
		Author: encoder.DBAuthorToAPI(author),
	})

	return res, nil
}

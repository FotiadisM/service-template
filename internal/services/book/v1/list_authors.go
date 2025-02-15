package bookv1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/encoder"
)

func (s *Service) ListAuthors(ctx context.Context, _ *connect.Request[bookv1.ListAuthorsRequest]) (*connect.Response[bookv1.ListAuthorsResponse], error) {
	authors, err := s.db.ListAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	resAuthors := []*bookv1.Author{}
	for _, author := range authors {
		resAuthors = append(resAuthors, encoder.DBAuthorToAPI(author))
	}

	res := connect.NewResponse(&bookv1.ListAuthorsResponse{
		Authors: resAuthors,
	})

	return res, nil
}

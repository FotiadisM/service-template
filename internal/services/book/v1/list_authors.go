package bookv1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *Service) ListAuthor(ctx context.Context, _ *connect.Request[bookv1.ListAuthorRequest]) (*connect.Response[bookv1.ListAuthorResponse], error) {
	authors, err := s.db.ListAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	resAuthors := []*bookv1.Author{}
	for _, a := range authors {
		resAuthors = append(resAuthors, &bookv1.Author{
			Id:        a.ID.String(),
			Name:      a.Name,
			CreatedAt: timestamppb.New(a.CreatedAt),
			UpdatedAt: timestamppb.New(a.UpdatedAt),
		})
	}

	res := connect.NewResponse(&bookv1.ListAuthorResponse{
		Authors: resAuthors,
	})

	return res, nil
}

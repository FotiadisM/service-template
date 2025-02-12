package bookv1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *Service) ListBook(ctx context.Context, _ *connect.Request[bookv1.ListBookRequest]) (*connect.Response[bookv1.ListBookResponse], error) {
	books, err := s.db.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	resBooks := []*bookv1.Book{}
	for _, a := range books {
		resBooks = append(resBooks, &bookv1.Book{
			Id:        a.ID.String(),
			Title:     a.Title,
			AuthorId:  a.AuthorID.String(),
			CreatedAt: timestamppb.New(a.CreatedAt),
			UpdatedAt: timestamppb.New(a.UpdatedAt),
		})
	}

	res := connect.NewResponse(&bookv1.ListBookResponse{
		Books: resBooks,
	})

	return res, nil
}

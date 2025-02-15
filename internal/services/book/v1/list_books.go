package bookv1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/encoder"
)

func (s *Service) ListBooks(ctx context.Context, _ *connect.Request[bookv1.ListBooksRequest]) (*connect.Response[bookv1.ListBooksResponse], error) {
	books, err := s.db.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	resBooks := []*bookv1.Book{}
	for _, book := range books {
		resBooks = append(resBooks, encoder.DBBookToAPI(book))
	}

	res := connect.NewResponse(&bookv1.ListBooksResponse{
		Books: resBooks,
	})

	return res, nil
}

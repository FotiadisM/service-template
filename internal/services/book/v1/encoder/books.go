package encoder

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func DBBookToAPI(book queries.Book) *bookv1.Book {
	return &bookv1.Book{
		Id:          book.ID.String(),
		Title:       book.Title,
		AuthorId:    book.AuthorID.String(),
		Description: book.Description,
		CreatedAt:   timestamppb.New(book.CreatedAt),
		UpdatedAt:   timestamppb.New(book.UpdatedAt),
	}
}

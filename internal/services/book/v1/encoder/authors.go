package encoder

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func DBAuthorToAPI(author queries.Author) *bookv1.Author {
	return &bookv1.Author{
		Id:        author.ID.String(),
		Name:      author.Name,
		Bio:       author.Bio,
		CreatedAt: timestamppb.New(author.CreatedAt),
		UpdatedAt: timestamppb.New(author.UpdatedAt),
	}
}

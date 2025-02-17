package encoder

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func DBBookReviewToAPI(br queries.BookReview) *bookv1.BookReview {
	return &bookv1.BookReview{
		Id:        br.ID.String(),
		BookId:    br.BookID.String(),
		Rating:    br.Rating,
		Text:      br.Text,
		CreatedAt: timestamppb.New(br.CreatedAt),
		UpdatedAt: timestamppb.New(br.UpdatedAt),
	}
}

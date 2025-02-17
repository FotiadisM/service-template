package bookv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *EndpointTestingSuite) TestCreateBookReview(t *testing.T) {
	ctx := t.Context()

	req := &bookv1.CreateBookReviewRequest{
		BookId: s.Fixtures.Book1.ID.String(),
		Rating: 2,
		Text:   "this is review",
	}
	res, err := s.Client.CreateBookReview(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	require.NotNil(t, res.Msg.Review)
	assert.NotEmpty(t, res.Msg.Review.Id)
	assert.Equal(t, s.Fixtures.Book1.ID.String(), res.Msg.Review.BookId)
	assert.Equal(t, req.Rating, res.Msg.Review.Rating)
	assert.Equal(t, req.Text, res.Msg.Review.Text)
}

func (s *UnitTestingSuite) TestCreateBookReviewHTTP(t *testing.T) {
	ctx := t.Context()

	s.DB.EXPECT().CreateBookReview(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, in queries.CreateBookReviewParams) (queries.BookReview, error) {
		now := time.Now()
		book := queries.BookReview{
			ID:        in.ID,
			BookID:    in.BookID,
			Rating:    in.Rating,
			Text:      in.Text,
			CreatedAt: now,
			UpdatedAt: now,
		}
		return book, nil
	})

	req_buf := &bytes.Buffer{}
	req_body := &bookv1.CreateBookReviewRequest{
		BookId: "0194fee7-3d16-7703-b28a-5b5c6ff6ecf4",
		Rating: 2,
		Text:   "this is review",
	}
	err := json.NewEncoder(req_buf).Encode(req_body)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/books/%s/reviews", s.ServerURL, req_body.BookId),
		req_buf,
	)
	require.NoError(t, err)

	res, err := s.HTTPClint.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	res_body := &bookv1.CreateBookReviewResponse{}
	err = json.NewDecoder(res.Body).Decode(res_body)
	require.NoError(t, err)
	require.NotEmpty(t, res_body.Review)

	s.DB.AssertExpectations(t)
}

func (s *UnitTestingSuite) TestCreateBookReviewValidation(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		req *bookv1.CreateBookReviewRequest
	}{
		{&bookv1.CreateBookReviewRequest{}},
		{&bookv1.CreateBookReviewRequest{Text: "book_title"}},
		{&bookv1.CreateBookReviewRequest{Rating: 2}},
		{&bookv1.CreateBookReviewRequest{Text: "book_title", Rating: 4}},
		{&bookv1.CreateBookReviewRequest{BookId: "bad_book_id", Text: "book_title", Rating: 4}},
		{&bookv1.CreateBookReviewRequest{BookId: "0194fae3-33da-7464-9f78-b2d37a9a75ca", Text: "book_title", Rating: -1}},
		{&bookv1.CreateBookReviewRequest{BookId: "0194fae3-33da-7464-9f78-b2d37a9a75ca", Text: "book_title", Rating: 6}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := s.Client.CreateBookReview(ctx, connect.NewRequest(tt.req))
			require.Error(t, err)
			cErr := &connect.Error{}
			require.ErrorAs(t, err, &cErr)
			require.Nil(t, res)

			assert.Equal(t, connect.CodeInvalidArgument, cErr.Code())
		})
	}
}

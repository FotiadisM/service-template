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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *EndpointTestingSuite) TestCreateBook(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	authorID, err := uuid.NewV7()
	require.NoError(t, err)
	author := queries.CreateAuthorParams{
		ID:        authorID,
		Name:      "authro",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.Service.db.CreateAuthor(ctx, author)
	require.NoError(t, err)

	req := connect.NewRequest(&bookv1.CreateBookRequest{
		Title:    "book_title",
		AuthorId: authorID.String(),
	})
	res, err := s.Client.CreateBook(ctx, req)
	require.NoError(t, err)

	require.NotEmpty(t, res.Msg.Book)
	assert.NotEmpty(t, res.Msg.Book.Id)
	assert.Equal(t, req.Msg.Title, res.Msg.Book.Title)
	assert.Equal(t, req.Msg.AuthorId, res.Msg.Book.AuthorId)
	assert.NotEmpty(t, res.Msg.Book.CreatedAt)
	assert.NotEmpty(t, res.Msg.Book.UpdatedAt)

	bookID, err := uuid.Parse(res.Msg.Book.Id)
	require.NoError(t, err)
	book, err := s.Service.db.GetBook(ctx, bookID)
	require.NoError(t, err)

	assert.Equal(t, book.ID.String(), res.Msg.Book.Id)
	assert.Equal(t, book.Title, res.Msg.Book.Title)
	assert.Equal(t, book.AuthorID.String(), res.Msg.Book.AuthorId)
}

func (s *UnitTestingSuite) TestCreateBookHTTP(t *testing.T) {
	ctx := context.Background()

	s.DB.EXPECT().CreateBook(mock.Anything, mock.Anything).Return(nil).Once()

	req_buf := &bytes.Buffer{}
	req_body := &bookv1.CreateBookRequest{
		Title:    "book_title",
		AuthorId: "0194fee7-3d16-7703-b28a-5b5c6ff6ecf4",
	}
	err := json.NewEncoder(req_buf).Encode(req_body)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/books", s.ServerURL),
		req_buf,
	)
	require.NoError(t, err)

	res, err := s.HTTPClint.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	res_body := &bookv1.CreateBookResponse{}
	err = json.NewDecoder(res.Body).Decode(res_body)
	require.NoError(t, err)
	require.NotEmpty(t, res_body.Book)

	s.DB.AssertExpectations(t)
}

func (s *UnitTestingSuite) TestCreateBookValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		req *bookv1.CreateBookRequest
	}{
		{&bookv1.CreateBookRequest{}},
		{&bookv1.CreateBookRequest{Title: "book_title"}},
		{&bookv1.CreateBookRequest{AuthorId: "bad_author_id"}},
		{&bookv1.CreateBookRequest{AuthorId: "0194fae3-33da-7464-9f78-b2d37a9a75ca"}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req := connect.NewRequest(tt.req)
			res, err := s.Client.CreateBook(ctx, req)
			cErr := &connect.Error{}
			require.ErrorAs(t, err, &cErr)
			require.Nil(t, res)

			assert.Equal(t, connect.CodeInvalidArgument, cErr.Code())
		})
	}
}

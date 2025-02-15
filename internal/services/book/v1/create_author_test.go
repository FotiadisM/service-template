package bookv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func (s *EndpointTestingSuite) TestCreateAuthor(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	authorReq := &bookv1.CreateAuthorRequest{
		Name: "author_name",
		Bio:  "author_bio",
	}
	res, err := s.Client.CreateAuthor(ctx, connect.NewRequest(authorReq))
	require.NoError(t, err)

	require.NotEmpty(t, res.Msg.Author)
	assert.NotEmpty(t, res.Msg.Author.Id)
	assert.Equal(t, authorReq.Name, res.Msg.Author.Name)
	assert.Equal(t, authorReq.Bio, res.Msg.Author.Bio)
	assert.NotEmpty(t, res.Msg.Author.CreatedAt)
	assert.NotEmpty(t, res.Msg.Author.UpdatedAt)

	id, err := uuid.Parse(res.Msg.Author.Id)
	require.NoError(t, err)
	author, err := s.Service.db.GetAuthor(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, author.ID.String(), res.Msg.Author.Id)
	assert.Equal(t, author.Name, authorReq.Name)
	assert.Equal(t, author.Bio, authorReq.Bio)
}

func (s *UnitTestingSuite) TestCreateAuthorHTTP(t *testing.T) {
	ctx := context.Background()

	s.DB.EXPECT().CreateAuthor(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, cap queries.CreateAuthorParams) (queries.Author, error) {
		author := queries.Author{
			ID:        cap.ID,
			Name:      cap.Name,
			Bio:       cap.Bio,
			CreatedAt: cap.CreatedAt,
			UpdatedAt: cap.UpdatedAt,
		}
		return author, nil
	}).Once()

	reqBody := &bytes.Buffer{}
	authorReq := &bookv1.CreateAuthorRequest{
		Name: "author_name",
	}
	err := json.NewEncoder(reqBody).Encode(authorReq)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/authors", s.ServerURL),
		reqBody,
	)
	require.NoError(t, err)

	res, err := s.HTTPClint.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	res_body := &bookv1.CreateAuthorResponse{}
	err = json.NewDecoder(res.Body).Decode(res_body)
	require.NoError(t, err)
	require.NotEmpty(t, res_body.Author)

	s.DB.AssertExpectations(t)
}

func (s *UnitTestingSuite) TestCreateAuthorValidation(t *testing.T) {
	ctx := context.Background()

	req := connect.NewRequest(&bookv1.CreateAuthorRequest{
		Name: "",
	})
	res, err := s.Client.CreateAuthor(ctx, req)
	cErr := &connect.Error{}
	require.ErrorAs(t, err, &cErr)
	require.Nil(t, res)

	assert.Equal(t, connect.CodeInvalidArgument, cErr.Code())

	s.DB.AssertExpectations(t)
}

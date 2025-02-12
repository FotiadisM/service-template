package bookv1

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *EndpointTestingSuite) TestCreateAuthor(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	req := connect.NewRequest(&bookv1.CreateAuthorRequest{
		Name: "author_name",
	})
	res, err := s.Client.CreateAuthor(ctx, req)
	require.NoError(t, err)

	require.NotEmpty(t, res.Msg.Author)
	assert.NotEmpty(t, res.Msg.Author.Id)
	assert.Equal(t, req.Msg.Name, res.Msg.Author.Name)
	assert.NotEmpty(t, res.Msg.Author.CreatedAt)
	assert.NotEmpty(t, res.Msg.Author.UpdatedAt)

	id, err := uuid.Parse(res.Msg.Author.Id)
	require.NoError(t, err)
	author, err := s.Service.db.GetAuthor(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, author.ID.String(), res.Msg.Author.Id)
	assert.Equal(t, author.Name, res.Msg.Author.Name)
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

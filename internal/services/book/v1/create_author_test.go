package bookv1

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *EndpointTestingSuite) TestCreateAuthor(t *testing.T) {
	ctx := context.Background()

	req := &bookv1.CreateAuthorRequest{
		Name: "author_name",
	}
	res, err := s.Service.CreateAuthor(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	require.NotEmpty(t, res.Msg.Author)
	author := res.Msg.Author
	assert.NotEmpty(t, author.Id)
	assert.Equal(t, req.Name, author.Name)
	assert.NotEmpty(t, author.CreatedAt)
	assert.NotEmpty(t, author.UpdatedAt)
}

// func (s *UnitTestingSuite) TestCreateAuthorValidation(t *testing.T) {
// 	ctx := context.Background()
//
// 	req := &bookv1.CreateAuthorRequest{
// 		Name: "",
// 	}
// 	res, err := s.Service.CreateAuthor(ctx, connect.NewRequest(req))
// 	cErr := &connect.Error{}
// 	require.ErrorAs(t, err, cErr)
// 	require.Nil(t, res)
//
// 	assert.Equal(t, connect.CodeInvalidArgument, cErr.Code())
//
// 	s.DB.AssertExpectations(t)
// }

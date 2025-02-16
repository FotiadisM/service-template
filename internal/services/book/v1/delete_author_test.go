package bookv1

import (
	"database/sql"
	"testing"

	"connectrpc.com/connect"
	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *EndpointTestingSuite) TestDeleteAuthor(t *testing.T) {
	ctx := t.Context()

	req := &bookv1.DeleteAuthorRequest{
		Id: s.Fixtures.Author1.ID.String(),
	}
	_, err := s.Client.DeleteAuthor(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	_, err = s.Service.db.GetAuthor(ctx, s.Fixtures.Author1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

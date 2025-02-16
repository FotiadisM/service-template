package bookv1

import (
	"database/sql"
	"testing"

	"connectrpc.com/connect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

func (s *EndpointTestingSuite) TestDeleteBook(t *testing.T) {
	ctx := t.Context()

	req := &bookv1.DeleteBookRequest{
		Id: s.Fixtures.Book1.ID.String(),
	}
	_, err := s.Client.DeleteBook(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	_, err = s.Service.db.GetBook(ctx, s.Fixtures.Book1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

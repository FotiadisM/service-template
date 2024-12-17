package authv1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/mocks"
	"github.com/FotiadisM/mock-microservice/internal/db/repository"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	db := mocks.NewMockDB(t)
	srv := &Service{db: db}
	client := NewTestClient(t, srv)

	db.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, params repository.CreateUserParams) (repository.User, error) {
			return repository.User{
				ID:        params.ID,
				Email:     params.Email,
				Password:  params.Password,
				Scope:     params.Scope,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, nil
		})

	_, err := client.Register(ctx, &authv1.RegisterRequest{
		Email:    "test@mail.com",
		Password: "1234",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})
	assert.NoError(err)
}

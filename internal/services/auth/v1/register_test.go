package authv1

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/repository"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
)

func (s *UnitTestingSuite) TestRegister(t *suite.T) {
	ctx := context.Background()

	s.DB.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
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

	_, err := s.Client.Register(ctx, &authv1.RegisterRequest{
		Email:    "test@mail.com",
		Password: "1234",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})
	t.NoError(err)
}

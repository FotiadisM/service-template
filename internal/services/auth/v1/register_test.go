package authv1

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/FotiadisM/mock-microservice/internal/db/repository"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
)

func (s *UnitTestingSuite) TestRegister(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	require := require.New(t)

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

	req := connect.NewRequest(&authv1.RegisterRequest{
		Email:    "test@mail.com",
		Password: "0123456789",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})

	res, err := s.Client.Register(ctx, req)
	require.NoError(err)

	assert.NotZero(res.Msg.AccessToken)
	assert.NotZero(res.Msg.RefreshToken)
}

func (s *UnitTestingSuite) TestRegisterValidation(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name string
		body *authv1.RegisterRequest
	}{
		{
			name: "bad email",
			body: &authv1.RegisterRequest{
				Email:    "bademail",
				Password: "0123456789",
				UserType: authv1.UserType_USER_TYPE_APPLICANT,
			},
		},
		{
			name: "short password",
			body: &authv1.RegisterRequest{
				Email:    "test@email.com",
				Password: "short",
				UserType: authv1.UserType_USER_TYPE_APPLICANT,
			},
		},
		{
			name: "unspecified user type",
			body: &authv1.RegisterRequest{
				Email:    "test@email.com",
				Password: "0123456789",
				UserType: authv1.UserType_USER_TYPE_UNSPECIFIED,
			},
		},
		{
			name: "invalid user type",
			body: &authv1.RegisterRequest{
				Email:    "test@email.com",
				Password: "0123456789",
				UserType: 3,
			},
		},
		{
			name: "missing email",
			body: &authv1.RegisterRequest{
				Email:    "",
				Password: "0123456789",
				UserType: authv1.UserType_USER_TYPE_APPLICANT,
			},
		},
		{
			name: "missing password",
			body: &authv1.RegisterRequest{
				Email:    "test@email.com",
				Password: "",
				UserType: authv1.UserType_USER_TYPE_APPLICANT,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(tc.body)
			_, err := s.Client.Register(ctx, req)

			require.Error(t, err)
			connectErr := new(connect.Error)
			require.ErrorAs(t, err, &connectErr)
			assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
			assert.NotEmpty(t, connectErr.Details())
		})
	}
}

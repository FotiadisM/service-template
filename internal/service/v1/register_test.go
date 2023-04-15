package servicev1

import (
	"context"
	"database/sql"
	"time"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/service/errors"
	"github.com/FotiadisM/mock-microservice/internal/store"
	"github.com/FotiadisM/mock-microservice/internal/store/queries"
	"github.com/FotiadisM/mock-microservice/pkg/logger"

	"github.com/stretchr/testify/mock"
)

func (t *UnitTestSuit) TestRegister() {
	ctx := context.Background()
	ctx = logger.WithLogger(ctx, logger.New(true))

	mockStore := store.NewMockStore(t.T())
	mockStore.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
		return queries.User{
			ID:        params.ID,
			Email:     params.Email,
			Password:  params.Password,
			Scope:     params.Scope,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}, nil
	})

	svc := NewService(mockStore)
	_, err := svc.Register(ctx, &authv1.RegisterRequest{
		Email:    "mike@mail.com",
		Password: "1234",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})
	t.NoError(err)
}

func (t *UnitTestSuit) TestRegisterEmeilExists() {
	ctx := context.Background()
	ctx = logger.WithLogger(ctx, logger.New(true))

	mockStore := store.NewMockStore(t.T())
	mockStore.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(queries.User{}, sql.ErrNoRows)

	svc := NewService(mockStore)
	_, err := svc.Register(ctx, &authv1.RegisterRequest{
		Email:    "mike@mail.com",
		Password: "1234",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})

	t.ErrorIs(err, errors.ErrEmailExists)
}

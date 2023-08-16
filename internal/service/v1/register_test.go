package servicev1

import (
	"context"
	"time"

	authv1 "github.com/FotiadisM/mock-microservice/api/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store/mocks"
	"github.com/FotiadisM/mock-microservice/internal/store/queries"
	"github.com/stretchr/testify/mock"
)

func (t *UnitTestSuit) TestRegister() {
	ctx := context.Background()

	store := mocks.NewMockStore(t.T())
	store.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
		func(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
			return queries.User{
				ID:        params.ID,
				Email:     params.Email,
				Password:  params.Password,
				Scope:     params.Scope,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, nil
		})

	svc := NewService(store)
	_, err := svc.Register(ctx, &authv1.RegisterRequest{
		Email:    "mike@mail.com",
		Password: "1234",
		UserType: authv1.UserType_USER_TYPE_APPLICANT,
	})
	t.NoError(err)
}

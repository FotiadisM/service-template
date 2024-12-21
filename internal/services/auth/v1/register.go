package authv1

import (
	"context"

	"github.com/google/uuid"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/repository"
)

func (s *Service) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	_, err = s.db.CreateUser(ctx, repository.CreateUserParams{
		ID:       id,
		Email:    in.Email,
		Password: in.Password,
		Scope:    repository.UserScopeUser,
	})
	if err != nil {
		return nil, err
	}

	out := &authv1.RegisterResponse{
		AccessToken:  "1234",
		RefreshToken: "5678",
	}
	return out, nil
}

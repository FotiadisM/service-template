package authv1

import (
	"context"

	"github.com/google/uuid"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/repository"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

func (s *Service) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	log := ilog.FromContext(ctx)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	u, err := s.db.CreateUser(ctx, repository.CreateUserParams{
		ID:       id,
		Email:    in.Email,
		Password: in.Password,
		Scope:    repository.UserScopeUser,
	})
	if err != nil {
		return nil, err
	}

	log.Info("user info is", "email", u.Email, "password", u.Password)

	out := &authv1.RegisterResponse{
		AccessToken:  "1234",
		RefreshToken: "5678",
	}
	return out, nil
}

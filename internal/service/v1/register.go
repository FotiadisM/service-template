package servicev1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/store/queries"
	"github.com/FotiadisM/mock-microservice/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	log := logger.FromContext(ctx)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	u, err := s.store.CreateUser(ctx, queries.CreateUserParams{
		ID:       id,
		Email:    in.Email,
		Password: in.Password,
		Scope:    queries.UserScopeApplicant,
	})
	if err != nil {
		return nil, err
	}

	log.Info("user info is", zap.String("email", u.Email), zap.String("password", u.Password))

	out := &authv1.RegisterResponse{
		AccessToken:  "1234",
		RefreshToken: "5678",
	}
	return out, nil
}

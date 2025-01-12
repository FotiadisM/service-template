package authv1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/services/auth/v1/queries"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

func (s *Service) Register(ctx context.Context, req *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	log := ilog.FromContext(ctx)
	log.Info("hello, this is register")

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	_, err = s.db.CreateUser(ctx, queries.CreateUserParams{
		ID:       id,
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
		Scope:    queries.UserScopeUser,
	})
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&authv1.RegisterResponse{
		AccessToken:  "1234",
		RefreshToken: "5678",
	})

	return res, nil
}

package servicev1

import (
	"context"

	"github.com/FotiadisM/mock-microservice/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log := logger.FromContext(ctx)
	log.Info("HELLO")

	out := &emptypb.Empty{}
	return out, nil
}

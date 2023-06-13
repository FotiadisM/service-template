package servicev1

import (
	"context"

	authv1 "github.com/FotiadisM/mock-microservice/api/auth/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Login(_ context.Context, _ *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	st := status.New(codes.Internal, "Unexpected error")
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason:   "MY_CUSTOM-CODE",
		Domain:   "auth-svc",
		Metadata: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	return nil, st.Err()
}

package authv1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
)

var ErrHello = errors.New("hello I am an error")

func (s *Service) Login(_ context.Context, _ *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	err := connect.NewError(connect.CodeInternal, ErrHello)
	details, detailsErr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason:   "reason",
		Domain:   "domain",
		Metadata: map[string]string{"key": "value"},
	})
	if detailsErr != nil {
		panic(detailsErr)
	}
	err.AddDetail(details)

	return nil, err
}

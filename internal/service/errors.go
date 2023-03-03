package service

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrStupidError = NewError(codes.Internal, "WHATS_UP", "what's up my boy")

func NewError(code codes.Code, reason, msg string) error {
	return NewErrorWithDomain(code, reason, msg, "auth-svc")
}

func NewErrorWithDomain(code codes.Code, reason, msg, domain string) error {
	st := status.New(code, msg)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: reason,
		Domain: domain,
	})
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return st.Err()
}

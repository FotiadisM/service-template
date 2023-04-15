package errors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrEmailExists = NewError(codes.AlreadyExists, "auth-svc", "EMAIL_NOT_UNIQUE", "the email provided is already in use")

func NewError(code codes.Code, domain, reason, msg string) error {
	st := status.New(code, msg)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: reason,
		Domain: domain,
	})
	if err != nil {
		// If this errored, it will always error
		// here, better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata to status: %v", err))
	}

	return st.Err()
}

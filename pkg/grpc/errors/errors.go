package errors

import (
	"fmt"

	"github.com/FotiadisM/mock-microservice/pkg/version"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewInfoError(code codes.Code, reason, msg string, md map[string]string) error {
	st := status.New(code, msg)
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason:   reason,
		Domain:   version.ModulePath,
		Metadata: md,
	})
	if err != nil {
		// If this errored, it will always error
		// here, better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching details to status: %v", err))
	}

	return st.Err()
}

func NewBadRequestError(code codes.Code, msg string, fv ...*errdetails.BadRequest_FieldViolation) error {
	st := status.New(code, msg)
	if len(fv) == 0 {
		return st.Err()
	}

	st, err := st.WithDetails(&errdetails.BadRequest{
		FieldViolations: fv,
	})
	if err != nil {
		// If this errored, it will always error
		// here, better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching details to status: %v", err))
	}

	return st.Err()
}

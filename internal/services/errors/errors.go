package errors

import (
	"connectrpc.com/connect"

	bookv1 "github.com/FotiadisM/mock-microservice/api/gen/go/book/v1"
)

type ErrorCode string

var ErrorCodeMyErrorCode ErrorCode = "error_code_1"

type ServiceError struct {
	Err *bookv1.Error

	ConnectRPCCode connect.Code
}

func (e ServiceError) Error() string {
	if e.Err != nil {
		return e.Err.Name
	}

	return ""
}

var ErrMyError = ServiceError{
	Err: &bookv1.Error{
		Code:        string(ErrorCodeMyErrorCode),
		Name:        "my-error-name",
		Description: "my-error-description",
	},
	ConnectRPCCode: connect.CodeCanceled,
}

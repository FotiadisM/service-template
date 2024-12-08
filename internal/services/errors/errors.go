package errors

import (
	"google.golang.org/grpc/codes"

	"github.com/FotiadisM/mock-microservice/pkg/grpc/errors"
)

var ErrEmailExists = errors.NewInfoError(codes.AlreadyExists, "EMAIL_NOT_UNIQUE", "the email provided is already in use", nil)

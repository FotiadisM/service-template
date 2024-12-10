package validate

import (
	"errors"

	gerrors "github.com/FotiadisM/mock-microservice/pkg/grpc/errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func defaultErrWrapperFunc(err error) error {
	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		fvs := []*errdetails.BadRequest_FieldViolation{}
		for _, v := range valErr.Violations {
			fvs = append(fvs, &errdetails.BadRequest_FieldViolation{
				Field:       *v.FieldPath,
				Description: *v.Message,
			})
		}
		return gerrors.NewBadRequestError(codes.InvalidArgument, "Failed to validate request", fvs)
	}

	return status.Error(codes.InvalidArgument, err.Error())
}

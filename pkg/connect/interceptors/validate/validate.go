package validate

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
)

func DefaultErrorHanlder(err error) error {
	if err == nil {
		return nil
	}

	if tErr := new(protovalidate.ValidationError); errors.As(err, &tErr) {
		conErr := connect.NewError(connect.CodeInvalidArgument, err)

		var details *connect.ErrorDetail
		details, err = connect.NewErrorDetail(tErr.ToProto())
		if err != nil {
			panic(fmt.Sprintf("failed to create error details: %s", err))
		}
		conErr.AddDetail(details)

		return conErr
	}

	panic(fmt.Sprintf("error '%T' is not of type protovalidate.ValidationError: %s", err, err))
}

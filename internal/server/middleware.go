package server

import (
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"

	"github.com/bufbuild/protovalidate-go"

	"github.com/FotiadisM/mock-microservice/internal/config"
	svcErrors "github.com/FotiadisM/mock-microservice/internal/services/errors"
	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/errsanitizer"
	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/logging"
	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/recovery"
)

var errUnexpected = errors.New("unexpected error")

func OtelMiddleware() (connect.Interceptor, error) {
	m, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create otel middleware: %w", err)
	}
	return m, nil
}

func LoggingMiddleware(log *slog.Logger) connect.Interceptor {
	return logging.NewInterceptor(log)
}

func ValidationMiddleware() (connect.Interceptor, error) {
	m, err := validate.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create validation middleware: %w", err)
	}
	return m, nil
}

func RecoveryMiddleware() connect.Interceptor {
	return recovery.NewInterceptor()
}

func ErrSanitizerMiddleware() connect.Interceptor {
	errSanitizerFunc := func(err error) error {
		if tErr := new(protovalidate.ValidationError); errors.As(err, &tErr) {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		if tErr := new(connect.Error); errors.As(err, &tErr) {
			return err
		}

		if tErr := new(svcErrors.ServiceError); errors.As(err, &tErr) {
			cErr := connect.NewError(tErr.ConnectRPCCode, tErr)
			details, _ := connect.NewErrorDetail(tErr.Err)
			cErr.AddDetail(details)

			return cErr
		}

		return connect.NewError(connect.CodeInternal, errUnexpected)
	}

	return errsanitizer.NewInterceptor(errsanitizer.WithRecoveryFunc(errSanitizerFunc))
}

func ChainMiddleware(_ *config.Config, log *slog.Logger) []connect.Interceptor {
	otelInterceptor, err := OtelMiddleware()
	if err != nil {
		panic(err)
	}

	validationInterceptor, err := ValidationMiddleware()
	if err != nil {
		panic(err)
	}

	interceptors := []connect.Interceptor{
		ErrSanitizerMiddleware(),
		otelInterceptor,
		LoggingMiddleware(log),
		RecoveryMiddleware(),
		validationInterceptor,
	}

	return interceptors
}

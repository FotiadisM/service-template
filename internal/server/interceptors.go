package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/otelconnect"

	"github.com/FotiadisM/service-template/internal/config"
	svcErrors "github.com/FotiadisM/service-template/internal/services/errors"
	"github.com/FotiadisM/service-template/pkg/connect/interceptors/errsanitizer"
	"github.com/FotiadisM/service-template/pkg/connect/interceptors/logging"
	"github.com/FotiadisM/service-template/pkg/connect/interceptors/recovery"
	"github.com/FotiadisM/service-template/pkg/connect/interceptors/validate"
)

var errUnexpected = errors.New("unexpected error")

func OtelMiddleware() (connect.Interceptor, error) {
	filter := func(_ context.Context, spec connect.Spec) bool {
		name := strings.TrimLeft(spec.Procedure, "/")
		parts := strings.SplitN(name, "/", 2)
		if len(parts) != 2 {
			return false
		}

		switch parts[0] {
		case grpchealth.HealthV1ServiceName:
			return false
		}

		return true
	}

	m, err := otelconnect.NewInterceptor(otelconnect.WithFilter(filter))
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

package server

import (
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"

	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/logging"
)

func CreateInterceptors(log *slog.Logger) ([]connect.Interceptor, error) {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create otel interceptor: %w", err)
	}

	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create validation interceptor: %w", err)
	}
	loggingInterceptor := logging.NewInterceptor(log)

	interceptors := []connect.Interceptor{
		otelInterceptor,
		loggingInterceptor,
		validationInterceptor,
	}

	return interceptors, nil
}

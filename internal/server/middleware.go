package server

import (
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/pkg/connect/interceptors/logging"
)

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
		otelInterceptor,
		LoggingMiddleware(log),
		validationInterceptor,
	}

	return interceptors
}

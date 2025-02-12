package test

import (
	"log/slog"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	"github.com/FotiadisM/service-template/internal/config"
	"github.com/FotiadisM/service-template/internal/server"
)

func ChainMiddleware(t *testing.T, _ *config.Config) []connect.Interceptor {
	t.Helper()

	validationInterceptor, err := server.ValidationMiddleware()
	require.NoError(t, err, "failed to create validation interceptor")

	return []connect.Interceptor{
		server.LoggingMiddleware(slog.Default()),
		validationInterceptor,
	}
}

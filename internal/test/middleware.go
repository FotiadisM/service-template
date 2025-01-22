package test

import (
	"testing"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/stretchr/testify/require"
)

func NewMiddleware(t *testing.T) []connect.Interceptor {
	t.Helper()

	validationInterceptor, err := validate.NewInterceptor()
	require.NoError(t, err, "failed to create validation interceptor")

	return []connect.Interceptor{
		validationInterceptor,
	}
}

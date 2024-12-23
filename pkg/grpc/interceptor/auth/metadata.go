package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const headerAuthorize = "authorization"

// AuthFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func AuthFromMD(ctx context.Context, expectedScheme string) (string, error) { //nolint
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Failed to Authenticate")
	}

	v, ok := md[headerAuthorize]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Failed to Authenticate")
	}

	value := v[0] // TODO(FotiadisM): Do I have to check the length?
	if value == "" {
		return "", status.Error(codes.Unauthenticated, "Failed to Authenticate")
	}
	splits := strings.SplitN(value, " ", 2)
	if len(splits) < 2 {
		return "", status.Error(codes.Unauthenticated, "Failed to Authenticate")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Error(codes.Unauthenticated, "Failed to Authenticate")
	}
	return splits[1], nil
}

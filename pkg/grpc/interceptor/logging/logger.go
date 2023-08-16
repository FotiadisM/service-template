package logging

import (
	"log/slog"

	"google.golang.org/grpc/codes"
)

// DefaultServerCodeToLevel is the helper mapper that maps gRPC return codes to log levels for server side.
func DefaultServerCodeToLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK, codes.NotFound, codes.Canceled, codes.AlreadyExists, codes.InvalidArgument, codes.Unauthenticated:
		return slog.LevelInfo

	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted,
		codes.OutOfRange, codes.Unavailable:
		return slog.LevelWarn

	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
		return slog.LevelError

	default:
		return slog.LevelError
	}
}

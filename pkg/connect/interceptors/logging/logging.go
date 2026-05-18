package logging

import (
	"log/slog"
	"strconv"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// DefaultCodeToLevelFunc maps connectRPC return codes to slog levels.
func DefaultCodeToLevelFunc(code connect.Code) slog.Level {
	switch code {
	case connect.CodeNotFound, connect.CodeCanceled, connect.CodeAlreadyExists, connect.CodeInvalidArgument, connect.CodeUnauthenticated:
		return slog.LevelInfo

	case connect.CodeDeadlineExceeded, connect.CodePermissionDenied, connect.CodeResourceExhausted, connect.CodeFailedPrecondition, connect.CodeAborted,
		connect.CodeOutOfRange, connect.CodeUnavailable:
		return slog.LevelWarn

	case connect.CodeUnknown, connect.CodeUnimplemented, connect.CodeInternal, connect.CodeDataLoss:
		return slog.LevelError

	default:
		return slog.LevelError
	}
}

func DefaultErrorDetailsAttrFunc(details []*connect.ErrorDetail) []slog.Attr {
	returnAttrs := []slog.Attr{}

	for _, detail := range details {
		protoMsg, err := detail.Value()
		if err != nil {
			continue
		}

		switch t := protoMsg.(type) {
		case *errdetails.BadRequest:
			attrs := make([]slog.Attr, len(t.FieldViolations))
			for i, fv := range t.FieldViolations {
				attrs[i] = slog.Group(strconv.Itoa(i),
					slog.String("field", fv.Field),
					slog.String("description", fv.Description),
				)
			}
			returnAttrs = append(returnAttrs, slog.Attr{
				Key:   "error_badrequest",
				Value: slog.GroupValue(attrs...),
			})

		case *errdetails.DebugInfo:
			attrs := make([]slog.Attr, len(t.StackEntries))
			for i, v := range t.StackEntries {
				attrs[i] = slog.String(strconv.Itoa(i), v)
			}
			returnAttrs = append(returnAttrs, slog.Attr{
				Key: "error_debuginfo",
				Value: slog.GroupValue(
					slog.String("detail", t.Detail),
					slog.Attr{Key: "stack_entries", Value: slog.GroupValue(attrs...)}),
			})

		case *errdetails.ErrorInfo:
			index := 0
			attrs := make([]slog.Attr, len(t.Metadata))
			for k, v := range t.Metadata {
				attrs[index] = slog.String(k, v)
				index++
			}
			returnAttrs = append(returnAttrs, slog.Attr{
				Key: "error_errorinfo",
				Value: slog.GroupValue(
					slog.String("reason", t.Domain),
					slog.String("domain", t.Domain),
					slog.Attr{Key: "metadata", Value: slog.GroupValue(attrs...)}),
			})

		case *errdetails.PreconditionFailure:
			attrs := make([]slog.Attr, len(t.Violations))
			for i, v := range t.Violations {
				attrs[i] = slog.Attr{Key: strconv.Itoa(i), Value: slog.GroupValue(
					slog.String("type", v.Type),
					slog.String("subject", v.Subject),
					slog.String("description", v.Description),
				)}
			}
			returnAttrs = append(returnAttrs, slog.Attr{
				Key: "error_preconditionfailure",
				Value: slog.GroupValue(
					slog.Attr{Key: "violations", Value: slog.GroupValue(attrs...)}),
			})

		case *errdetails.RequestInfo:
			returnAttrs = append(returnAttrs, slog.Attr{
				Key: "error_preconditionfailure",
				Value: slog.GroupValue(
					slog.String("request_id", t.RequestId),
					slog.String("serving_data", t.ServingData),
				),
			})
		}
	}

	return returnAttrs
}

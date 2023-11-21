package validate

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var ErrUnsupportedMessageType = errors.New("unsupported message type")

func UnaryServerInterceptor(validator *protovalidate.Validator, opts ...Option) grpc.UnaryServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, ErrUnsupportedMessageType
		}

		err = validator.Validate(msg)
		if err != nil {
			return nil, options.errWrapperFn(err)
		}

		return handler(ctx, req)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream

	options   *options
	validator *protovalidate.Validator
}

func (ws *wrappedServerStream) RecvMsg(m any) error {
	msg, ok := m.(proto.Message)
	if !ok {
		return ErrUnsupportedMessageType
	}

	err := ws.validator.Validate(msg)
	if err != nil {
		return ws.options.errWrapperFn(err)
	}

	return ws.ServerStream.RecvMsg(m)
}

func StreamServerInterceptor(validator *protovalidate.Validator, opts ...Option) grpc.StreamServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ws := &wrappedServerStream{ServerStream: ss, options: options, validator: validator}

		return handler(srv, ws)
	}
}

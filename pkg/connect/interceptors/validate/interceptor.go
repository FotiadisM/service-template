package validate

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

type streamingClientInterceptor struct {
	connect.StreamingClientConn

	errHanlderFn ErrorHandlerFunc
	validator    *protovalidate.Validator
}

func (s *streamingClientInterceptor) Send(msg any) error {
	err := validate(s.validator, msg)
	if err != nil {
		return s.errHanlderFn(err)
	}

	return s.StreamingClientConn.Send(msg)
}

type streamingHandlerInterceptor struct {
	connect.StreamingHandlerConn

	validator *protovalidate.Validator
}

func (s *streamingHandlerInterceptor) Receive(msg any) error {
	if err := s.StreamingHandlerConn.Receive(msg); err != nil {
		return err
	}

	return validate(s.validator, msg)
}

func validate(validator *protovalidate.Validator, msg any) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		panic(fmt.Sprintf("expected proto.Message, got %T", msg))
	}

	err := validator.Validate(protoMsg)
	if err == nil {
		return nil
	}

	connectErr := connect.NewError(connect.CodeInvalidArgument, err)
	if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
		if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
			connectErr.AddDetail(detail)
		}
	}
	return connectErr
}

type Interceptor struct {
	validator    *protovalidate.Validator
	errHanlderFn ErrorHandlerFunc
}

var _ connect.Interceptor = &Interceptor{}

func NewInterceptor(opts ...Option) (*Interceptor, error) {
	options := defaultOptions()
	for _, fn := range opts {
		fn(options)
	}

	if options.validator == nil {
		var err error
		options.validator, err = protovalidate.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create validator: %w", err)
		}
	}

	interceprot := &Interceptor{
		validator:    options.validator,
		errHanlderFn: options.errHanlderFn,
	}

	return interceprot, nil
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := validate(i.validator, req.Any()); err != nil {
			return nil, i.errHanlderFn(err)
		}

		return next(ctx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return &streamingClientInterceptor{
			validator:           i.validator,
			errHanlderFn:        i.errHanlderFn,
			StreamingClientConn: next(ctx, spec),
		}
	}
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		err := next(ctx, &streamingHandlerInterceptor{
			validator:            i.validator,
			StreamingHandlerConn: conn,
		})
		if err != nil {
			return i.errHanlderFn(err)
		}

		return nil
	}
}

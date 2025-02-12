package validate

import "github.com/bufbuild/protovalidate-go"

type ErrorHandlerFunc func(err error) error

type options struct {
	validator    protovalidate.Validator
	errHanlderFn ErrorHandlerFunc
}

func defaultOptions() *options {
	return &options{
		errHanlderFn: DefaultErrorHanlder,
	}
}

type Option func(o *options)

func WithValidator(validator protovalidate.Validator) Option {
	return func(o *options) {
		o.validator = validator
	}
}

func WithErrHandlerFn(fn ErrorHandlerFunc) Option {
	return func(o *options) {
		o.errHanlderFn = fn
	}
}

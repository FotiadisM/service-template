package recovery

import "context"

type RecoveryFunc func(ctx context.Context, p any) (err error) //nolint

type options struct {
	recoveryFn RecoveryFunc
}

func defaultOptions() *options {
	return &options{
		recoveryFn: nil,
	}
}

type Option func(o *options)

func WithRecoveryFunc(fn RecoveryFunc) Option {
	return func(o *options) {
		o.recoveryFn = fn
	}
}

package errsanitizer

type SanitizerFunc func(err error) error

type options struct {
	snFunc SanitizerFunc
}

func defaultOptions() *options {
	return &options{
		snFunc: func(err error) error {
			return err
		},
	}
}

type Option func(o *options)

func WithRecoveryFunc(fn SanitizerFunc) Option {
	return func(o *options) {
		o.snFunc = fn
	}
}

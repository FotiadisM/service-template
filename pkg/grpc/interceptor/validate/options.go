package validate

type ErrWrapperFunc func(err error) error

type options struct {
	errWrapperFn ErrWrapperFunc
}

func defaultOptions() *options {
	return &options{
		errWrapperFn: defaultErrWrapperFunc,
	}
}

type Option func(o *options)

func WithErrWrapperFunc(fn ErrWrapperFunc) Option {
	return func(o *options) {
		o.errWrapperFn = fn
	}
}

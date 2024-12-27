package logging

type options struct{}

func defaultOptions() *options {
	return &options{}
}

type Option func(*options)

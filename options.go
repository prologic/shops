package main

type Options struct {
	ContinueOnError bool
}

func NewOptions() *Options {
	return &Options{
		ContinueOnError: false,
	}
}

type Option func(opts *Options) error

func WithContinueOnError(conf bool) Option {
	return func(opts *Options) error {
		opts.ContinueOnError = cont
		return nil
	}
}

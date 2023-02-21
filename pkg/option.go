package pkg

import "time"

type Options struct {
	baseURL string
	proxy   string
	timeout time.Duration
	debug   bool
	logger  Logger
}

type Option func(opts *Options)

func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func WithBaseURL(baseURL string) Option {
	return func(opts *Options) {
		opts.baseURL = baseURL
	}
}

func WithProxy(proxy string) Option {
	return func(opts *Options) {
		opts.proxy = proxy
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}

func WithDebug(debug bool) Option {
	return func(opts *Options) {
		opts.debug = debug
	}
}

func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.logger = logger
	}
}

package pkg

import "time"

type Options struct {
	baseURL string
	proxy   string
	timeout time.Duration
	debug   bool
}

type Option func(opts *Options)

func WithOptions(opt *Options) Option {
	return func(opts *Options) {
		*opts = *opt
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

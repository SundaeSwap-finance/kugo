// Copyright 2022 SundaeSwap Labs, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software
// is furnished to do so, subject to the following conditions:
//
// Licensed under the MIT License;
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://opensource.org/licenses/MIT
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package kugo

import (
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
)

// Options available to kugo client
type Options struct {
	endpoint string
	headers  map[string]string
	timeout  time.Duration
	logger   ogmigo.Logger
}

// Option to kugo client
type Option func(*Options)

// Keep the http connection open as long as possible
func WithoutTimeout() Option {
	return func(opts *Options) {
		opts.timeout = 0
	}
}

// Set a specific timeout for all requests
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}

// WithEndpoint allows kupo endpoint to be set; defaults to http://127.0.0.1:1442
func WithEndpoint(endpoint string) Option {
	return func(opts *Options) {
		opts.endpoint = endpoint
	}
}

// WithHeader allows extra header to be attached
func WithHeader(name string, value string) Option {
	return func(opts *Options) {
		opts.headers[name] = value
	}
}

// WithLogger allows custom logger to be specified
func WithLogger(logger ogmigo.Logger) Option {
	return func(opts *Options) {
		opts.logger = logger
	}
}

func buildOptions(opts ...Option) Options {
	var options Options
	options.timeout = 5 * time.Minute
	for _, opt := range opts {
		opt(&options)
	}
	if options.endpoint == "" {
		options.endpoint = "http://127.0.0.1:1442"
	}
	if options.logger == nil {
		options.logger = ogmigo.DefaultLogger
	}
	return options
}

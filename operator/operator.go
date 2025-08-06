// Package operator provides utilities for working with channels in Go.
// It offers a simple, configurable way to create typed channels with optional buffering,
// specifically designed to work with trx.Result types for type-safe error handling
// in concurrent operations.
package operator

import (
	"context"

	"github.com/foreveraloneT/trx"
)

// config holds configuration options for channel creation.
// This struct is used internally to store settings provided through functional options.
type config struct {
	bufferSize int  // Size of the channel buffer (0 = unbuffered)
	poolSize   int  // Number of worker goroutines in the pool (must be > 0)
	serialize  bool // Serialize output when poolSize >= 1
	ctx        context.Context
}

// Option represents an option for the channel utility.
// This follows the functional options pattern, providing a flexible way to configure
// channel creation with optional parameters.
type Option func(*config)

// WithBufferSize sets the buffer size of the channel.
// A buffer size of 0 creates an unbuffered channel (synchronous communication).
// A positive buffer size creates a buffered channel that can hold that many values
// before blocking senders. Negative values are ignored and the default (0) is used.
//
// Example:
//
//	WithBufferSize(100) // Creates a buffered channel with capacity 100
//	WithBufferSize(0)   // Creates an unbuffered channel (default)
func WithBufferSize(size int) Option {
	return func(c *config) {
		if size >= 0 {
			c.bufferSize = size
		}
	}
}

// WithPoolSize returns an Option that sets the pool size in the operator configuration.
// If the provided size is greater than 0, it updates the pool size; otherwise, it leaves it unchanged.
//
// Example:
//
//	WithPoolSize(5) // Sets the pool size to 5 worker goroutines
//	WithPoolSize(1) // Sets the pool size to 1 (default)
//	WithPoolSize(0) // Invalid, pool size remains unchanged (default is 1)
func WithPoolSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.poolSize = size
		}
	}
}

// WithSerialize returns an Option that enables serialization in the operator configuration.
//
// Example:
//
//	WithSerialize() // Enables serialization in the operator
func WithSerialize() Option {
	return func(c *config) {
		c.serialize = true
	}
}

func defaultConfig() *config {
	return &config{
		bufferSize: 0,
		poolSize:   1, // Default pool size is 1
		serialize:  false,
	}
}

func parseOption(opts ...Option) *config {
	c := defaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func makeResultChannel[T any](c *config) chan trx.Result[T] {
	return make(chan trx.Result[T], c.bufferSize)
}

func makePool(c *config) *pool {
	return newPool(c.poolSize, c.serialize)
}

func makeContext(c *config) context.Context {
	if c.ctx != nil {
		return c.ctx
	}

	return context.Background()
}

func prepareResources[T any](opts ...Option) (ctx context.Context, out chan trx.Result[T], pool *pool) {
	c := parseOption(opts...)
	ctx = makeContext(c)
	out = makeResultChannel[T](c)
	pool = makePool(c)

	return
}

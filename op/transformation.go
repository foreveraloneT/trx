package op

import (
	"time"

	"github.com/foreveraloneT/trx"
)

// Map applies the provided mapper function to each item received from the source channel,
// emitting the results to a new output channel. The mapper function receives the value and its
// index in the sequence, and may return an error. If an error occurs during mapping or when
// retrieving the value from the source, the error is sent downstream wrapped in a trx.Result.
//
// The function supports optional configuration via Option parameters, such as context control
// and concurrency settings. Mapping operations are performed concurrently using a worker pool,
// and the output channel is closed once all mapping operations are complete.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//	U - The type of output values after mapping.
//
// Parameters:
//
//	source - A receive-only channel of trx.Result[T] representing the input stream.
//	mapper - A function that maps each value and its index to a new value of type U, possibly returning an error.
//	options
//	    - WithBufferSize
//	    - WithPoolSize
//	    - WithSerialize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[U] containing the mapped results or errors.
//
// Example usage:
//
//	out := Map(source, func(v int, i int) (string, error) {
//	    return strconv.Itoa(v), nil
//	})
func Map[T, U any](source <-chan trx.Result[T], mapper func(value T, index int) (U, error), options ...Option) <-chan trx.Result[U] {
	conf := parseOption(options...)
	ctx := makeContext(conf)
	out := makeResultChannel[U](conf)
	pool := makePool(conf)

	go func() {
		defer close(out)

		i := 0
	LOOP:
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-source:
				if !ok {
					break LOOP
				}

				index := i
				result := v

				pool.submit(func() callback {
					value, err := result.Get()
					if err != nil {
						return func() {
							out <- trx.Err[U](err)
						}
					}

					mapped, err := mapper(value, index)
					if err != nil {
						return func() {
							out <- trx.Err[U](err)
						}
					}

					return func() {
						out <- trx.Ok(mapped)
					}
				})

				i++
			}
		}

		pool.wait()
	}()

	return out
}

// BufferWithCount collects items from the source channel into fixed-size buffers and emits them as slices.
// Each emitted slice contains up to 'count' items. If the source channel closes and there are remaining items
// that do not fill a complete buffer, the final slice will contain the remaining items.
//
// The function supports optional configuration via Option parameters, such as context control and buffer size.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//
// Parameters:
//
//	source  - A receive-only channel of trx.Result[T] representing the input stream.
//	count   - The number of items per buffer (must be > 0).
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[[]T] containing the buffered slices or errors.
//
// Example usage:
//
//	out := BufferWithCount(source, 3, WithBufferSize(10), WithContext(ctx))
func BufferWithCount[T any](source <-chan trx.Result[T], count int, options ...Option) <-chan trx.Result[[]T] {
	conf := parseOption(options...)
	ctx := makeContext(conf)
	out := makeResultChannel[[]T](conf)

	go func() {
		defer close(out)

		buffer := make([]T, 0, count)
	LOOP:
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-source:
				if !ok {
					break LOOP
				}

				value, err := v.Get()
				if err != nil {
					out <- trx.Err[[]T](err)

					return
				}

				buffer = append(buffer, value)
				if len(buffer) >= count {
					out <- trx.Ok(buffer)

					buffer = make([]T, 0, count)
				}
			}
		}

		if len(buffer) > 0 {
			out <- trx.Ok(buffer)
		}
	}()

	return out
}

// BufferWithTime collects items from the source channel into time-based buffers and emits them as slices.
// Each emitted slice contains items collected within the specified duration or up to 'maxSize' items.
// If 'maxSize' is 0, the buffer is emitted only based on the timer. If the source channel closes and there
// are remaining items that do not fill a complete buffer, the final slice will contain the remaining items.
//
// The function supports optional configuration via Option parameters, such as context control and buffer size.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//
// Parameters:
//
//	source  - A receive-only channel of trx.Result[T] representing the input stream.
//	d       - The duration to wait before emitting the buffer.
//	maxSize - The maximum number of items per buffer (if 0, only time is considered).
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[[]T] containing the buffered slices or errors.
//
// Example usage:
//
//	out := BufferWithTime(source, time.Second, 5, WithBufferSize(10), WithContext(ctx))
func BufferWithTime[T any](source <-chan trx.Result[T], d time.Duration, maxSize int, options ...Option) <-chan trx.Result[[]T] {
	conf := parseOption(options...)
	ctx := makeContext(conf)
	out := makeResultChannel[[]T](conf)

	go func() {
		defer close(out)

		buffer := make([]T, 0)

		timer := time.NewTicker(d)
		defer timer.Stop()

	LOOP:
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if len(buffer) > 0 {
					out <- trx.Ok(buffer)
					buffer = make([]T, 0)
				}
			case v, ok := <-source:
				if !ok {
					break LOOP
				}

				value, err := v.Get()
				if err != nil {
					out <- trx.Err[[]T](err)

					return
				}

				buffer = append(buffer, value)
				if maxSize > 0 && len(buffer) >= maxSize {
					out <- trx.Ok(buffer)
					buffer = make([]T, 0)
					timer.Reset(d)
				}
			}
		}

		if len(buffer) > 0 {
			out <- trx.Ok(buffer)
		}
	}()

	return out
}

// BufferWithTimeOrCount collects items from the source channel into buffers and emits them as slices
// either when the specified time duration has elapsed or when the buffer reaches the specified count, whichever comes first.
// If the source channel closes and there are remaining items that do not fill a complete buffer, the final slice will contain the remaining items.
//
// The function supports optional configuration via Option parameters, such as context control and buffer size.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//
// Parameters:
//
//	source  - A receive-only channel of trx.Result[T] representing the input stream.
//	d       - The duration to wait before emitting the buffer.
//	count   - The maximum number of items per buffer (must be > 0).
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[[]T] containing the buffered slices or errors.
//
// Example usage:
//
//	out := BufferWithTimeOrCount(source, time.Second, 5, WithBufferSize(10), WithContext(ctx))
func BufferWithTimeOrCount[T any](source <-chan trx.Result[T], d time.Duration, count int, options ...Option) <-chan trx.Result[[]T] {
	conf := parseOption(options...)
	ctx := makeContext(conf)
	out := makeResultChannel[[]T](conf)

	go func() {
		defer close(out)

		buffer := make([]T, 0)

		timer := time.NewTicker(d)
		defer timer.Stop()

	LOOP:
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if len(buffer) > 0 {
					out <- trx.Ok(buffer)
					buffer = make([]T, 0)
				}
			case v, ok := <-source:
				if !ok {
					break LOOP
				}

				value, err := v.Get()
				if err != nil {
					out <- trx.Err[[]T](err)

					return
				}

				buffer = append(buffer, value)
				if count > 0 && len(buffer) >= count {
					out <- trx.Ok(buffer)
					buffer = make([]T, 0)
				}
			}
		}

		if len(buffer) > 0 {
			out <- trx.Ok(buffer)
		}
	}()

	return out
}

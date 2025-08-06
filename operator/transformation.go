package operator

import "github.com/foreveraloneT/trx"

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
	ctx, out, pool := prepareResources[U](options...)

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

				pool.submit(func() {
					value, err := result.Get()
					if err != nil {
						out <- trx.Err[U](err)

						return
					}

					mapped, err := mapper(value, index)
					if err != nil {
						out <- trx.Err[U](err)

						return
					}

					out <- trx.Ok(mapped)
				})

				i++
			}
		}

		pool.wait()
	}()

	return out
}

// BufferCount collects items from the source channel into fixed-size buffers and emits them as slices.
// Each emitted slice contains up to 'n' items. If the source channel closes and there are remaining items
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
//	n       - The number of items per buffer (must be > 0).
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
//	out := BufferCount(source, 3, WithBufferSize(10), WithContext(ctx))
func BufferCount[T any](source <-chan trx.Result[T], n int, options ...Option) <-chan trx.Result[[]T] {
	ctx, out, _ := prepareResources[[]T](options...)

	go func() {
		defer close(out)

		buffer := make([]T, 0, n)
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
				if len(buffer) == n {
					out <- trx.Ok(buffer)

					buffer = make([]T, 0, n)
				}
			}
		}

		if len(buffer) > 0 {
			out <- trx.Ok(buffer)
		}
	}()

	return out
}

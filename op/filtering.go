package op

import "github.com/foreveraloneT/trx"

// Filter emits only those values from the source channel for which the predicate function returns true.
// The predicate receives each value and its index, and may return an error. If an error occurs during
// filtering or when retrieving the value from the source, the error is sent downstream wrapped in a trx.Result.
//
// The function supports optional configuration via Option parameters, such as context control and concurrency
// settings. Filtering operations are performed concurrently using a worker pool, and the output channel is
// closed once all filtering operations are complete.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//
// Parameters:
//
//	source   - A receive-only channel of trx.Result[T] representing the input stream.
//	predicate - A function that determines if a value and its index should be included, possibly returning an error.
//	options
//	    - WithBufferSize
//	    - WithPoolSize
//	    - WithSerialize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[T] containing the filtered results or errors.
//
// Example usage:
//
//	out := Filter(source, func(v int, i int) (bool, error) {
//	    return v%2 == 0, nil // filter even numbers
//	})
func Filter[T any](source <-chan trx.Result[T], predicate func(value T, index int) (bool, error), options ...Option) <-chan trx.Result[T] {
	ctx, out, pool := prepareResources[T](options...)

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
						out <- trx.Err[T](err)

						return
					}

					ok, err := predicate(value, index)
					if err != nil {
						out <- trx.Err[T](err)

						return
					}

					if ok {
						out <- trx.Ok(value)
					}
				})

				i++
			}
		}

		pool.wait()
	}()

	return out
}

// Take emits up to n values from the source channel and then stops.
// The function reads from the source channel of trx.Result[T] and forwards up to n successful values
// to the output channel. If an error is encountered in the source, it is sent downstream wrapped in a trx.Result,
// and iteration stops. The function also stops if the source channel is closed or the context is cancelled.
//
// The function supports optional configuration via Option parameters, such as context control.
//
// Type Parameters:
//
//	T - The type of input values from the source channel.
//
// Parameters:
//
//	source - A receive-only channel of trx.Result[T] representing the input stream.
//	n      - The maximum number of values to emit.
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[T] containing up to n results or errors.
//
// Example usage:
//
//	out := Take(source, 5)
//	for res := range out {
//	    // handle res
//	}
func Take[T any](source <-chan trx.Result[T], n int, options ...Option) <-chan trx.Result[T] {
	ctx, out, _ := prepareResources[T](options...)

	go func() {
		defer close(out)

		count := 0
		for count < n {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-source:
				if !ok {
					return
				}

				val, err := v.Get()
				if err != nil {
					out <- trx.Err[T](err)

					return
				}

				out <- trx.Ok(val)

				count++
			}
		}
	}()

	return out
}

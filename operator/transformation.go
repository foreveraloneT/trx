package operator

import "github.com/foreveraloneT/trx"

// Map transforms the values from the source channel using the mapper function
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

// BufferCount buffers the source channel values until the buffer size is reached, then emits the buffer and starts a new buffer.
// Unable to use with `WithPoolSize`
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

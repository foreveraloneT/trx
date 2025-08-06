package operator

import "github.com/foreveraloneT/trx"

// Filter emits values from the source channel that pass the predicate function
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

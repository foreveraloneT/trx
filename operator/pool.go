package operator

import (
	basePool "github.com/sourcegraph/conc/pool"
	"github.com/sourcegraph/conc/stream"
)

type pool struct {
	pool   *basePool.Pool
	stream *stream.Stream
}

func (p *pool) submit(fn func()) {
	if p.pool != nil {
		p.pool.Go(fn)

		return
	}

	if p.stream != nil {
		p.stream.Go(func() stream.Callback {
			return fn
		})

		return
	}

	fn()
}

func (p *pool) wait() {
	if p.pool != nil {
		p.pool.Wait()

		return
	}

	if p.stream != nil {
		p.stream.Wait()

		return
	}
}

func newPool(size int, serialize bool) *pool {
	if size <= 1 {
		return &pool{}
	}

	if !serialize {
		return &pool{
			pool: basePool.New().WithMaxGoroutines(size),
		}
	}

	return &pool{
		stream: stream.New().WithMaxGoroutines(size),
	}
}

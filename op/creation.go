package op

import (
	"time"

	"github.com/foreveraloneT/trx"
)

// Timer emits a single trx.Result[int] after the specified duration d has elapsed.
// If the context is cancelled before the duration, the channel is closed without emitting.
//
// Type Parameters:
//
//	None.
//
// Parameters:
//
//	d       - The duration to wait before emitting a value.
//	options
//			- WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[int] that emits 0 after the specified duration or closes early if cancelled.
//
// Example usage:
//
//	out := Timer(2 * time.Second)
func Timer(d time.Duration, options ...Option) <-chan trx.Result[int] {
	ctx, out, _ := prepareResources[int]()

	go func() {
		defer close(out)

		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
			out <- trx.Ok(0)
		}
	}()

	return out
}

// Interval emits a trx.Result[int] at each interval specified by the duration d, incrementing the value each time.
// If the context is cancelled, the channel is closed without emitting further values.
//
// Type Parameters:
//
//	None.
//
// Parameters:
//
//	d       - The duration between emissions.
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[int] that emits incrementing integers at each interval.
//
// Example usage:
//
//	out := Interval(1 * time.Second)
func Interval(d time.Duration, options ...Option) <-chan trx.Result[int] {
	ctx, out, _ := prepareResources[int]()

	go func() {
		defer close(out)

		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				out <- trx.Ok(i)
			}
		}
	}()

	return out
}

// FormSlice emits each element of the provided slice source as a trx.Result[T] on the returned channel.
// If the context is cancelled, the channel is closed without emitting further values.
//
// Type Parameters:
//
//	T - The type of elements in the input slice.
//
// Parameters:
//
//	source   - The slice of values to emit.
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[T] that emits each element of source in order.
//
// Example usage:
//
//	out := FormSlice([]int{1, 2, 3})
func FormSlice[T any](source []T, options ...Option) <-chan trx.Result[T] {
	ctx, out, _ := prepareResources[T](options...)

	go func() {
		defer close(out)

		for _, v := range source {
			select {
			case <-ctx.Done():
				return
			default:
				out <- trx.Ok(v)
			}
		}
	}()

	return out
}

// FormChannel creates a new output channel of trx.Result[T] from the given source channel.
// It applies the provided options to configure the channel behavior, such as buffer size.
// The function launches a goroutine that reads values from the source channel and sends
// them as trx.Ok results to the output channel. If the context is cancelled or the source
// channel is closed, the output channel is closed as well.
//
// Parameters:
//
//	source: The input channel of type T to read values from.
//	options
//			- WithBufferSize
//			- WithContext
//
// Returns:
//   - A receive-only channel of trx.Result[T] containing the wrapped values from the source channel.
func FormChannel[T any](source <-chan T, options ...Option) <-chan trx.Result[T] {
	opts := append([]Option{WithBufferSize(cap(source))}, options...)

	ctx, out, _ := prepareResources[T](opts...)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-source:
				if !ok {
					return
				}
				out <- trx.Ok(v)
			}
		}
	}()

	return out
}

// Range emits a sequence of trx.Result[int], starting from 'start' and producing 'count' consecutive values.
// If the context is cancelled, the channel is closed without emitting further values.
//
// Type Parameters:
//
//	None.
//
// Parameters:
//
//	start    - The starting integer value of the sequence.
//	count    - The number of consecutive values to emit.
//	options
//	    - WithBufferSize
//	    - WithContext
//
// Returns:
//
//	A receive-only channel of trx.Result[int] that emits integers from start to start+count-1.
//
// Example usage:
//
//	out := Range(0, 5)
func Range(start int, count int, options ...Option) <-chan trx.Result[int] {
	ctx, out, _ := prepareResources[int](options...)

	go func() {
		defer close(out)

		for i := start; i < start+count; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				out <- trx.Ok(i)
			}
		}
	}()

	return out
}

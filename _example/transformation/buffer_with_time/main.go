package main

import (
	"fmt"
	"time"

	"github.com/foreveralonet/trx/op"
)

func main() {
	withoutMaxBufferSize()
	withMaxBufferSize()
	thereIsSomeError()
}

func withoutMaxBufferSize() {
	fmt.Println("Without Max Buffer Size:")

	source := op.Take(op.Interval(400*time.Millisecond), 10)

	// NOTE: maxSize = 0 mean no max buffer size
	out := op.BufferWithTime(source, 1*time.Second, 0)

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Buffered: %v\n", v)
	}
}

func withMaxBufferSize() {
	fmt.Println("With Max Buffer Size:")

	source := op.Take(op.Interval(400*time.Millisecond), 10)

	// NOTE: maxSize = 3
	out := op.BufferWithTime(source, 1*time.Second, 2)

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Buffered: %v\n", v)
	}
}

func thereIsSomeError() {
	fmt.Println("There Is Some Error:")

	source := op.Map(op.Take(op.Interval(400*time.Millisecond), 10), func(v int, _ int) (int, error) {
		if v == 6 {
			return 0, fmt.Errorf("error")
		}

		return v, nil
	})

	out := op.BufferWithTime(source, 1*time.Second, 0)

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Buffered: %v\n", v)
	}
}

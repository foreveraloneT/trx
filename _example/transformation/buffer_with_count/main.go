package main

import (
	"fmt"
	"time"

	"github.com/foreveralonet/trx/op"
)

func main() {
	basicUsage()
}

func basicUsage() {
	fmt.Println("Basic Usage:")

	source := op.Take(op.Interval(400*time.Millisecond), 10)

	out := op.BufferWithCount(source, 3)

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Buffered: %v\n", v)
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/trx/op"
)

func main() {
	basicUsage()
}

func basicUsage() {
	source := op.Take(op.Interval(400*time.Millisecond), 10)

	out := op.Map(source, func(v int, index int) (string, error) {
		return fmt.Sprintf("Value: %d, Index: %d", v, index), nil
	})

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Println(v)
	}
}

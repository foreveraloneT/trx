package main

import (
	"context"
	"fmt"
	"time"

	"github.com/foreveraloneT/trx"
	"github.com/foreveraloneT/trx/operator"
)

func main() {
	source := make(chan trx.Result[int])

	go func() {
		defer close(source)

		for i := 0; i < 50; i++ {
			source <- trx.Ok(i)
			<-time.After(1 * time.Second)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(5 * time.Second)
		cancel() // Cancel the context after 10 seconds
	}()

	out := operator.Map(source, func(v int, index int) (string, error) {
		return fmt.Sprintf("Value: %d, Index: %d", v, index), nil
	}, operator.WithBufferSize(10), operator.WithPoolSize(3), operator.WithSerialize(), operator.WithContext(ctx))

	for i := range out {
		v, err := i.Get()
		if err != nil {
			panic(err)
		}

		fmt.Println(v)
	}
}

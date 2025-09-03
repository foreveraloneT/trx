package main

import (
	"fmt"
	"time"

	"github.com/foreveralonet/trx/op"
)

func main() {
	basicUsage()
	usingWorkerPool()
	usingSerializedPool()
}

func basicUsage() {
	fmt.Println("Basic Usage:")

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

func usingWorkerPool() {
	fmt.Println("Using Worker Pool:")

	source := op.Range(0, 10, op.WithBufferSize(3))

	out := op.Map(source, func(v int, index int) (string, error) {
		<-time.After(1 * time.Second) // Simulate some processing delay

		return fmt.Sprintf("Value: %d, Index: %d", v, index), nil
	}, op.WithPoolSize(3))

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Println(v)
	}
}

func usingSerializedPool() {
	fmt.Println("Using Serialized Pool:")

	source := op.Range(0, 10, op.WithBufferSize(3))

	out := op.Map(source, func(v int, index int) (string, error) {
		<-time.After(1 * time.Second) // Simulate some processing delay

		return fmt.Sprintf("Value: %d, Index: %d", v, index), nil
	}, op.WithPoolSize(3), op.WithSerialize())

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}

		fmt.Println(v)
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/trx/op"
)

func main() {
	basicTakeUsage()
}

func basicTakeUsage() {
	fmt.Println("Take Example:")

	source := op.Interval(300 * time.Millisecond)
	out := op.Take(source, 5)

	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

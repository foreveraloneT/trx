package main

import (
	"fmt"
	"time"

	"github.com/foreveralonet/trx/op"
)

func main() {
	exampleTimer()
}

func exampleTimer() {
	fmt.Println("Timer Example:")
	out := op.Timer(1 * time.Second)
	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

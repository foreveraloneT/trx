package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/trx/op"
)

func main() {
	exampleInterval()
}

func exampleInterval() {
	fmt.Println("Interval Example:")
	out := op.Take(op.Interval(300*time.Millisecond), 5)
	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

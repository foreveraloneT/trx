package main

import (
	"fmt"

	"github.com/foreveralonet/trx/op"
)

func main() {
	exampleRange()
}

func exampleRange() {
	fmt.Println("Range Example:")
	out := op.Range(5, 3)
	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

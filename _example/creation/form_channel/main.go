package main

import (
	"fmt"

	"github.com/foreveralonet/trx/op"
)

func main() {
	exampleFormChannel()
}

func exampleFormChannel() {
	fmt.Println("FormChannel Example:")
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	out := op.FormChannel(ch)
	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

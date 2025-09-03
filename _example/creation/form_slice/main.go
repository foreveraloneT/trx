package main

import (
	"fmt"

	"github.com/foreveralonet/trx/op"
)

func main() {
	exampleFormSlice()
}

func exampleFormSlice() {
	fmt.Println("FormSlice Example:")
	slice := []string{"a", "b", "c"}
	out := op.FormSlice(slice)
	for val := range out {
		v, err := val.Get()
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}

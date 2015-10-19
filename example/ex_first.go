package main

import (
	"fmt"
	"log"

	mc "github.com/iostrovok/go-memorycache/memorycache"
)

func main() {
	// 10 x 1000 = 10000 is total cache size
	mc.NewSingleton(10, 1000)

	for i := 0; i < 1000; i++ {
		mc.Put(i, fmt.Sprintf("==> %d", i))
	}

	for i := 0; i < 1000; i++ {
		k, ok := mc.Get(fmt.Sprintf("key : %d", i))
		if !ok {
			log.Fatalf("Nof found i : %d\n", i)
		}
		kInt, ok := k.(int)
		if !ok {
			log.Fatalf("Bad type for i : %d\n", i)
		}

		fmt.Printf("Find ==> %d : %d\n", i, kInt)
	}

}

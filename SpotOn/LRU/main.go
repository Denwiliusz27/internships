package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("LRU implementation")

	val := 2000

	lru := newLru(100)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		start := time.Now()
		x := fib(nil, val)
		elapsed := time.Since(start)

		fmt.Printf("NO LRU -- fib(%d): %d time(%v)\n", val, x, elapsed)
	}()

	fmt.Println("********")

	go func() {
		defer wg.Done()
		start := time.Now()
		x := fib(lru, val)
		elapsed := time.Since(start)

		fmt.Printf("LRU -- fib(%d): %d time(%v)\n", val, x, elapsed)
	}()

	wg.Wait()
}

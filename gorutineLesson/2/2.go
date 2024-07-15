package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	channel := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			n := rand.Intn(10)
			channel <- n
		}
		close(channel)
	}()

	wg.Add(1) // why 1????
	go func() {
		defer wg.Done()
		sum := 0
		for num := range channel {
			sum += num
		}
		fmt.Printf("The amount: %d\n", sum)
	}()

	wg.Wait()
}

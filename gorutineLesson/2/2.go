package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	channel := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			channel <- i
			time.Sleep(time.Second)
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

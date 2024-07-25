package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	strings := []string{"flower", "sun", "cloud"}

	channel := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, str := range strings {
			channel <- str
			time.Sleep(time.Second)
		}
		close(channel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for str := range channel {
			fmt.Printf("-> %s\n", str)
			time.Sleep(time.Second)
		}

	}()
	wg.Wait()
}

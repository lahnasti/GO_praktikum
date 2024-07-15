/*
Создайте программу, которая использует атомарные операции для
безопасного увеличения значения переменной в нескольких горутинах.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

var counter int64

func main() {
	numGoroutines := 3
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("Conclusion: %d", counter)
}

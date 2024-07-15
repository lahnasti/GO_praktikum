/*
	напишите программу которая использует атомарные операции

для подсчета общего количества выполненных операций
несколькими горутинами
*/
package main

import (
	"sync"
	"sync/atomic"
	"time"
	"fmt"
)


var sum int64

var wg sync.WaitGroup

func worker() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		atomic.AddInt64(&sum, 1)
		time.Sleep(time.Second)
	}
}

func main() {
	num := 3
	wg.Add(num)
	for i := 0; i < num; i++ {
		go worker()
	}
	wg.Wait()
	fmt.Printf("Total operations: %d\n", sum)
}
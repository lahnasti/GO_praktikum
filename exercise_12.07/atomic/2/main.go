/*
Реализуйте программу, которая использует атомарные операции для управления
флагом завершения выполнения горутины.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var done int32

var wg sync.WaitGroup

func worker() {
	for {
		if atomic.LoadInt32(&done) == 1 {
			fmt.Println("Worker closed")
			wg.Done()
			return
		}
		fmt.Println("It's a work!")
		time.Sleep(time.Second)
	}
}

func main() {
	wg.Add(1)
	go worker()
	time.Sleep(10 * time.Second)
	atomic.StoreInt32(&done, 1)
	time.Sleep(1 * time.Second)
	fmt.Println("Programm done")
	wg.Wait()
}

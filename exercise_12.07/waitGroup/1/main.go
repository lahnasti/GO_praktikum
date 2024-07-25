/*
Реализуйте программу, которая использует WaitGroup для синхронизации
выполнения нескольких горутин, выполняющих вычисления.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(9)
	for i := 0; i < 10; i++ {
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Printf("Worker %d started\n", id)
			time.Sleep(2 * time.Second)
			fmt.Printf("Worker %d finished\n", id)
		}(i, &wg)
	}
	wg.Wait()
	fmt.Println("\nAll workers have finished")
}

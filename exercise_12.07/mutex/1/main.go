/*
Напишите программу, которая использует мьютексы
для защиты доступа к общей карте (map).
*/
package main

import (
	"fmt"
	"sync"
)

var (
	wg        = sync.WaitGroup{}
	mutex     = sync.Mutex{}
	sharedMap = make(map[int]string)
)

func writeToMap(key int, value string) {
	defer wg.Done()
	mutex.Lock()
	sharedMap[key] = value
	mutex.Unlock()
}

func readFromMap(key int) string {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()
	return sharedMap[key]
}

func main() {
	wg.Add(6)

	go writeToMap(1, "It's a first goroutine")
	go writeToMap(2, "It's a second goroutine")
	go writeToMap(3, "It's a thirty goroutine")

	go func() {
		fmt.Println("Value: ", readFromMap(1))
		fmt.Println("Value: ", readFromMap(2))
		wg.Done()
	}()

	wg.Wait()
}

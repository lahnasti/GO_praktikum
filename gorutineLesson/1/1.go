package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	channel := make(chan int)

	var wg sync.WaitGroup

	//первая горутина отправляет числа в канал
	wg.Add(1)
	go func() {
		for i := 1; i < 6; i++ {
			defer wg.Done()
			channel <- i
			time.Sleep(time.Second)
		}
		close(channel) //закрываем канал после отпр всех чисел
	}()

	//вторая горутина читает из канала и выводит в консоль
	wg.Add(4) // ? почему 4
	go func() {
		defer wg.Done()
		for num := range channel {
			fmt.Println(num)
		}
	}()
	wg.Wait()
}

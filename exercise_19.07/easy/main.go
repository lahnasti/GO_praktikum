/*Создайте программу, которая запускает две горутины,
каждая из которых отправляет данные в свой канал с различной задержкой.
Используйте select для получения данных из любого из каналов и выводите полученные данные на экран.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))
			ch1 <- "Data one"
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(2000) * int(time.Millisecond)))
			ch2 <- "Data two"
		}
	}()

	for {
		select {
		case message1 := <-ch1:
			fmt.Printf("First message: %s\n", message1)
		case message2 := <-ch2:
			fmt.Printf("Second message: %s\n", message2)
		}
	}
}

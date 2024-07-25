/*
Напишите программу, которая запускает горутину для выполнения длительной работы
(например, ожидание нескольких секунд) и использует канал для получения результата.
Добавьте тайм-аут для чтения из канала, чтобы программа завершалась, если работа не будет выполнена
за отведенное время.

*/

package main

import (
	"fmt"
	"time"
)

func writer(c chan<- int) {
	time.Sleep(time.Second * 1)
	c <- 10
}

func main() {
	c := make(chan int)

	go writer(c)

	select {
	case good := <-c:
		fmt.Printf("Result: %d", good)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout. The work was not completed within the allotted time.")
	}
}

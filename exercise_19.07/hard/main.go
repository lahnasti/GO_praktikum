/*Напишите программу, которая использует канал для ограничения количества одновременных горутин.
Пусть программа запускает 10 горутин, каждая из которых выполняет какую-либо работу (например, ждет одну секунду).
Используйте канал для ограничения количества одновременно работающих горутин до 3.
*/

package main

import "sync"

func worker(id int, wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	

}

func main() {
	const maxGoroutines = 3 // одновременно работающие горутины
	const numGoroutines = 10 // всего горутин

	ch := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go worker(i, &wg, ch)
	}

	wg.Wait()
}


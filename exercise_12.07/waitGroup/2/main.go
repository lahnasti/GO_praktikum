/*
Напишите программу, которая использует WaitGroup для ожидания завершения
нескольких операций чтения из файла.
*/

package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	files := []string{"file_1.txt", "file_2.txt", "file_3.txt"}

	wg.Add(len(files))
	for _, file := range files {
		go func(filename string) {
			defer wg.Done()

			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", filename, err)
				return
			}
			fmt.Printf("- %s\n", data)

		}(file)

	}
	wg.Wait()
}

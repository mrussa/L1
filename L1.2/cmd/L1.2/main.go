package main

import (
	"fmt"
	"sync"
)

var numbers [5]int = [5]int{2, 4, 6, 8, 10}

func work(id int, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	square := n * n
	fmt.Printf("Горутина %d: %d^2 = %d\n", id, n, square)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(len(numbers))

	for i, n := range numbers {
		go work(i+1, n, &wg)
	}

	wg.Wait()
	fmt.Println("Горутины завершили выполнение")
}

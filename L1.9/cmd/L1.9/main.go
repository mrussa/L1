package main

import (
	"fmt"
)

func main() {
	input := []int{1, 2, 3, 4, 5}

	chIn := make(chan int)
	chOut := make(chan int)

	go func() {
		for _, x := range input {
			chIn <- x
		}
		close(chIn)
	}()

	go func() {
		for x := range chIn {
			chOut <- x * 2
		}
		close(chOut)
	}()

	for y := range chOut {
		fmt.Println(y)
	}
}

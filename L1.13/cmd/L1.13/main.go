package main

import (
	"fmt"
)

func swapAddSub(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}

func main() {
	a, b := 10, 42
	fmt.Println("before:", a, b)

	x, y := swapAddSub(a, b)
	fmt.Println("after:", x, y)
}

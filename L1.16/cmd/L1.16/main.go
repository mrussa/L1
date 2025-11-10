package main

import (
	"fmt"
)

func quickSort(a []int) []int {
	if len(a) < 2 {
		return a
	}

	var qs func(l, r int)
	qs = func(l, r int) {
		i, j := l, r
		pivot := a[(l+r)/2]

		for i <= j {
			for a[i] < pivot {
				i++
			}
			for a[j] > pivot {
				j--
			}
			if i <= j {
				a[i], a[j] = a[j], a[i]
				i++
				j--
			}
		}
		if l < j {
			qs(l, j)
		}
		if i < r {
			qs(i, r)
		}
	}

	qs(0, len(a)-1)
	return a
}

func main() {
	data := []int{5, 2, 9, 1, 5, 6, 3, 7}
	fmt.Println("before:", data)
	quickSort(data)
	fmt.Println("after: ", data)
}

package main

import (
	"fmt"
)

func binSearch(a []int, target int) int {
	l, r := 0, len(a)-1
	for l <= r {
		mid := l + (r-l)/2
		switch {
		case a[mid] == target:
			return mid
		case a[mid] < target:
			l = mid + 1
		default:
			r = mid - 1
		}
	}
	return -1
}

func main() {
	arr := []int{1, 3, 5, 7, 9, 11, 13}
	fmt.Println(binSearch(arr, 1))   // 0
	fmt.Println(binSearch(arr, 7))   // 3
	fmt.Println(binSearch(arr, 13))  // 6
	fmt.Println(binSearch(arr, 2))   // -1
	fmt.Println(binSearch(arr, 100)) // -1
}

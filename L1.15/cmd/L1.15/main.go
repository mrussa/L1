package main

import "strings"

var justString string

func someFunc() {
	v := createHugeString(1 << 20)
	if len(v) > 100 {
		justString = strings.Clone(v[:100])
	} else {
		justString = strings.Clone(v)
	}
}

func main() {
	someFunc()
}

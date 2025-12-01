package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverseWordsInPlace(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}

	reverse(r, 0, len(r)-1)

	start := 0
	for i := 0; i <= len(r); i++ {
		if i == len(r) || r[i] == ' ' {
			reverse(r, start, i-1)
			start = i + 1
		}
	}
	return string(r)
}

func reverse(r []rune, i, j int) {
	for i < j {
		r[i], r[j] = r[j], r[i]
		i++
		j--
	}
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	if in.Scan() {
		fmt.Println(reverseWordsInPlace(in.Text()))
	}
	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
}

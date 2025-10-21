package main

import (
	"fmt"
	"sort"
)

func main() {
	src := []string{"cat", "cat", "dog", "cat", "tree"}

	uniq := make(map[string]struct{}, len(src))
	for _, s := range src {
		uniq[s] = struct{}{}
	}

	keys := make([]string, 0, len(uniq))
	for k := range uniq {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println(keys)
}

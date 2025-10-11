package main

import (
	"fmt"
	"sort"
	"strings"
)

func formatSlice(fs []float64) string {
	out := make([]string, len(fs))
	for i, v := range fs {
		out[i] = fmt.Sprintf("%.1f", v)
	}
	return "[" + strings.Join(out, ", ") + "]"
}

func main() {
	temps := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	groups := make(map[int][]float64)
	for _, t := range temps {
		bucket := int(t/10.0) * 10
		groups[bucket] = append(groups[bucket], t)
	}

	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%d: %s\n", k, formatSlice(groups[k]))
	}
}

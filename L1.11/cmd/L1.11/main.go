package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseLine(line string) ([]int, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	line = strings.ReplaceAll(line, ",", " ")
	fields := strings.Fields(line)

	out := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("некорректное число %q", f)
		}
		out = append(out, n)
	}
	return out, nil
}

func readListStrict(r *bufio.Reader, prompt string) ([]int, error) {
	for {
		fmt.Print(prompt)
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		nums, err := parseLine(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка:", err, "— попробуйте ещё раз.")
			continue
		}
		return nums, nil
	}
}

func intersectKeepOrder(a, b []int) []int {
	inA := make(map[int]struct{}, len(a))
	for _, x := range a {
		inA[x] = struct{}{}
	}
	seen := make(map[int]struct{})
	res := make([]int, 0, len(b))
	for _, y := range b {
		if _, ok := inA[y]; ok {
			if _, dup := seen[y]; !dup {
				seen[y] = struct{}{}
				res = append(res, y)
			}
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)

	a, err := readListStrict(in, "Введите A (через пробел или запятую): ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "не удалось прочитать A:", err)
		return
	}
	b, err := readListStrict(in, "Введите B (через пробел или запятую): ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "не удалось прочитать B:", err)
		return
	}

	c := intersectKeepOrder(a, b)

	fmt.Print("{")
	for i, v := range c {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Print(v)
	}
	fmt.Println("}")
}

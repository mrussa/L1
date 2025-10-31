package main

import (
	"fmt"
	"reflect"
)

func detectType(v any) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	}
	if reflect.TypeOf(v).Kind() == reflect.Chan {
		return "chan"
	}
	return fmt.Sprintf("unknown (%T)", v)
}

func main() {
	tests := []any{
		42,
		"hello",
		true,
		make(chan int),
		make(chan struct{}),
		3.14,
	}

	for _, x := range tests {
		fmt.Printf("%v -> %s\n", x, detectType(x))
	}
}

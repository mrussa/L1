package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func setBit(n int64, i uint, one bool) (int64, error) {
	if i > 63 {
		return 0, errors.New("bit index must be 0..63")
	}
	m := int64(1) << i
	if one {
		return n | m, nil
	}
	return n &^ m, nil
}

func main() {
	n := flag.Int64("n", 5, "int64 value")
	i := flag.Uint("i", 1, "bit index (0..63)")
	v := flag.Int("v", 0, "bit value (0 or 1)")
	flag.Parse()

	if *v != 0 && *v != 1 {
		fmt.Fprintln(os.Stderr, "error: -v must be 0 or 1")
		os.Exit(2)
	}

	out, err := setBit(*n, *i, *v == 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	fmt.Printf("n=%d i=%d v=%d -> %d  (bin %08b -> %08b)\n",
		*n, *i, *v, out, uint8(*n), uint8(out))
}

package main

import (
	"flag"
	"fmt"
	"sync"
)

type Counter struct {
	mu sync.Mutex
	n  int64
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *Counter) Value() int64 {
	c.mu.Lock()
	v := c.n
	c.mu.Unlock()
	return v
}

func main() {
	workers := flag.Int("workers", 16, "goroutines")
	iters := flag.Int("iters", 100000, "increments per goroutine")
	flag.Parse()

	var c Counter
	var wg sync.WaitGroup
	wg.Add(*workers)

	for w := 0; w < *workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < *iters; i++ {
				c.Inc()
			}
		}()
	}

	wg.Wait()
	got := c.Value()
	want := int64((*workers) * (*iters))
	fmt.Printf("count=%d (expected %d)\n", got, want)
	if got != want {
		fmt.Println("mismatch")
	}
}

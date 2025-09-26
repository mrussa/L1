package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type CounterMap struct {
	mu sync.Mutex
	m  map[string]int
}

func NewCounterMap() *CounterMap {
	return &CounterMap{m: make(map[string]int)}
}

func (c *CounterMap) Add(key string, delta int) {
	c.mu.Lock()
	c.m[key] += delta
	c.mu.Unlock()
}

func (c *CounterMap) Get(key string) (int, bool) {
	c.mu.Lock()
	v, ok := c.m[key]
	c.mu.Unlock()
	return v, ok
}

func (c *CounterMap) Len() int {
	c.mu.Lock()
	n := len(c.m)
	c.mu.Unlock()
	return n
}

func main() {
	workers := flag.Int("workers", 16, "горутин")
	keys := flag.Int("keys", 100, "ключей")
	iters := flag.Int("iters", 1000, "инкрементов на горутину")
	flag.Parse()

	store := NewCounterMap()
	fmt.Printf("run: workers=%d keys=%d iters=%d\n", *workers, *keys, *iters)

	var wg sync.WaitGroup
	wg.Add(*workers)

	for w := 0; w < *workers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < *iters; i++ {
				k := "k" + strconv.Itoa(i%*keys)
				store.Add(k, 1)
				if i%257 == 0 {
					time.Sleep(time.Microsecond)
				}
			}
		}(w)
	}

	wg.Wait()

	total := 0
	for i := 0; i < *keys; i++ {
		if v, ok := store.Get("k" + strconv.Itoa(i)); ok {
			total += v
		}
	}
	expect := (*workers) * (*iters)

	fmt.Printf("done: distinct=%d, sum=%d (expected %d)\n", store.Len(), total, expect)
	for i := 0; i < 3 && i < *keys; i++ {
		k := "k" + strconv.Itoa(i)
		v, _ := store.Get(k)
		fmt.Printf("  %s=%d\n", k, v)
	}
	if total != expect {
		fmt.Println("mismatch")
	} else {
		fmt.Println("OK")
	}
}

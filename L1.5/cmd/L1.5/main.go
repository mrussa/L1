package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"
)

const defaultSecond = 5

func main() {
	n := flag.Int("n", defaultSecond, "time working (seconds)")
	flag.Parse()

	if *n <= 0 {
		log.Printf("start, N=%ds (<=0) — exiting immediately", *n)
		return
	}

	timeout := time.Duration(*n) * time.Second
	log.Printf("start, N=%v", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dataCh := make(chan int)
	dataC := dataCh

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	tickC := ticker.C

	var wg sync.WaitGroup
	wg.Add(1)
	go func(out chan<- int) {
		defer wg.Done()
		defer close(out)

		pace := time.NewTicker(150 * time.Millisecond)
		defer pace.Stop()

		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				return
			case <-pace.C:
			}
			select {
			case <-ctx.Done():
				return
			case out <- i:
			}
		}
	}(dataCh)

	valCount := 0
	tickCount := 0

	for {
		select {
		case t := <-tickC:
			tickCount++
			log.Printf("tick #%d at %s", tickCount, t.Format("15:04:05.000"))

		case v, ok := <-dataC:
			if !ok {
				log.Println("dataCh closed — waiting for timeout")
				dataC = nil
				tickC = nil
				continue
			}
			valCount++
			log.Printf("got value #%d from dataCh: %d", valCount, v)

		case <-ctx.Done():
			log.Println("timeout, shutting down")
			wg.Wait()
			return
		}
	}

}

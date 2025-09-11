package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "workerpool",
		Usage:   "run N workers that read from a shared channel",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "workers",
				Aliases: []string{"w"},
				Usage:   "number of workers to start",
				Value:   3,
			},
			&cli.IntFlag{
				Name:  "buffer",
				Usage: "buffer size of the jobs channel (0 = unbuffered)",
				Value: 0,
			},
			&cli.DurationFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Usage:   "delay between generated items (e.g. 100ms, 1s)",
				Value:   0,
			},
		},
		Action: func(c *cli.Context) error {
			n := c.Int("workers")
			if n <= 0 {
				return cli.Exit("workers must be > 0", 1)
			}

			buf := c.Int("buffer")
			interval := c.Duration("interval")

			sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			jobs := make(chan int, buf)

			var wg sync.WaitGroup
			wg.Add(n)
			for i := 1; i <= n; i++ {
				id := i
				go func() {
					defer wg.Done()
					for job := range jobs {
						fmt.Printf("worker-%d: %d\n", id, job)
					}
				}()
			}

			fmt.Printf("starting %d workers (buffer=%d, interval=%v)… press Ctrl+C to stop\n", n, buf, interval)

			counter := 1
			for {
				select {
				case <-sigCtx.Done():
					close(jobs)
					wg.Wait()
					return nil
				case jobs <- counter:
					counter++
					if interval > 0 {
						time.Sleep(interval)
					}
				}
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		if exitErr, ok := err.(cli.ExitCoder); ok {
			fmt.Fprintln(os.Stderr, exitErr.Error())
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

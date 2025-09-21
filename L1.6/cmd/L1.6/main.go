package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func sayHello() { fmt.Println("[hello] started") }

// 1) Остановка через закрытие канала done (закрывает отправитель)
func worker(done <-chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("[worker] stop signal received")
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

// 2) Остановка через контекст (таймаут/ручная отмена)
func longRunningTask(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[long] cancelled")
			return ctx.Err()
		default:
			time.Sleep(time.Second)
		}
	}
}

// 3) Каскадная отмена: родительский ctx гасит дочернюю задачу
func parentTask(parent context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		childCtx, _ := context.WithCancel(parent) // без локального cancel — отмена придёт сверху
		childTask(childCtx)
	}()
}

func childTask(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[child] cancelled by parent")
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

// 4) Возврат «контекстной» ошибки как есть — часто удобно в вызывающем коде
func robustTask(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("robust: %w", ctx.Err())
	default:
		return nil
	}
}

// Просто печатаем с паузой — используется в нескольких местах
func run(iteration int, name string) {
	for i := 0; i < iteration; i++ {
		time.Sleep(time.Second)
		fmt.Println("[run]", name)
	}
}

// 5) Демонстрация runtime.Goexit: завершает ТЕКУЩУЮ горутину, выполняя её defer
func runGoexit(name string) {
	defer fmt.Println("[goexit]", name, "defer executed")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("[goexit]", name, "calling runtime.Goexit()")
	runtime.Goexit()
	// сюда не дойдём
}

// 6) Остановка через atomic-флаг (без каналов/контекста)
func runAtomic(stop *atomic.Bool, iteration int, name string) {
	for i := 0; i < iteration; i++ {
		if stop.Load() {
			fmt.Println("[atomic]", name, "stopped via flag")
			return
		}
		time.Sleep(time.Second)
		fmt.Println("[atomic]", name)
	}
}

// 7) Остановка по таймеру (time.NewTimer)
func runTimer(iteration int, name string, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	for i := 0; i < iteration; i++ {
		select {
		case <-timer.C:
			fmt.Println("[timer]", name, "timeout")
			return
		case <-time.After(time.Second):
			fmt.Println("[timer]", name)
		}
	}
}

func main() {
	// 0) просто отдельная горутина «для разогрева»
	go sayHello()
	time.Sleep(200 * time.Millisecond)

	// (1) done-канал
	done := make(chan struct{})
	go worker(done)
	time.Sleep(1800 * time.Millisecond)
	close(done) // сигнал остановки

	// (2)+(3) общий контекст с таймаутом (каскадная отмена дочерней задачи)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wgCtx sync.WaitGroup
	parentTask(ctx, &wgCtx)

	if err := longRunningTask(ctx); err != nil {
		fmt.Println("[main] long:", err)
	}
	if err := robustTask(ctx); err != nil {
		fmt.Println("[main]", err)
	}

	// (4) Goexit — ТОЛЬКО внутри дочерней горутины
	var wgExit sync.WaitGroup
	wgExit.Add(1)
	go func() {
		defer wgExit.Done()
		runGoexit("demo")
	}()

	// (5) atomic-flag
	var stop atomic.Bool
	var wgAtomic sync.WaitGroup
	wgAtomic.Add(1)
	go func() { defer wgAtomic.Done(); runAtomic(&stop, 10, "A") }()
	time.Sleep(2500 * time.Millisecond)
	stop.Store(true) // подаём флаг остановки

	// (6) таймер
	var wgTimer sync.WaitGroup
	wgTimer.Add(1)
	go func() { defer wgTimer.Done(); runTimer(10, "T", 2500*time.Millisecond) }()

	// (7) системный сигнал → контекст (Ctrl-C). Ниже всё равно завершим вручную.
	ctxSig, stopSig := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopSig()

	var wgSig sync.WaitGroup
	wgSig.Add(1)
	go func() {
		defer wgSig.Done()
		_ = longRunningTask(ctxSig)
		fmt.Println("[signal] task stopped (interrupt or program end)")
	}()

	// Ждём все демонстрации
	wgCtx.Wait()  // child завершится по таймауту родителя
	wgExit.Wait() // goexit-гороутина корректно завершилась (через defer)
	wgAtomic.Wait()
	wgTimer.Wait()

	// Если Ctrl-C не нажимали — аккуратно гасим signal-контекст
	stopSig()
	wgSig.Wait()

	fmt.Println("[main] clean exit")
}

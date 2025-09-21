Небольшой учебный проект, показывающий **базовые способы остановки горутин** в Go:

1. Остановка по **закрытию канала** (`done`).
2. Остановка через **контекст** (`WithTimeout` / каскадная отмена).
3. Возврат **контекстной ошибки** (ошибкоустойчивое завершение).
4. Демонстрация **`runtime.Goexit()`** (завершает текущую горутину, выполняя `defer`).
5. Остановка по **атомарному флагу** (`atomic.Bool`).
6. Остановка по **таймеру** (`time.NewTimer`).
7. Остановка по **системному сигналу** (Ctrl-C) через `signal.NotifyContext`.

Код устроен так, чтобы каждую технику было видно в логах, и программа завершалась «чисто».

---

## Быстрый старт

```bash
# запустить с детектором гонок
make race
```

---

## Пример вывода

```text
[hello] started
[worker] stop signal received
[child] cancelled by parent
[long] cancelled
[main] long: context deadline exceeded
[main] robust: context deadline exceeded
[goexit] demo calling runtime.Goexit()
[goexit] demo defer executed
[atomic] A
[atomic] A
[atomic] A
[atomic] A stopped via flag
[timer] T
[timer] T
[timer] T timeout
[long] cancelled
[signal] task stopped (interrupt or program end)
[main] clean exit
```

> Порядок строк в выводе может слегка отличаться из-за конкуренции — это нормально.

---

## Что демонстрируется

- **[done-канал]** `worker(done)` — завершение по закрытию канала отправителем.
- **[контекст: таймаут]** `longRunningTask(ctx)` — остановка по `ctx.Done()`.
- **[контекст: каскад]** `parentTask(ctx) → childTask(childCtx)` — родительский `ctx` гасит дочерний.
- **[контекст: ошибка]** `robustTask(ctx)` — возвращаем наверх `ctx.Err()`.
- **[Goexit]** `runGoexit("demo")` — внутри *дочерней* горутины вызывается `runtime.Goexit()`; видно, что `defer` отрабатывает.
- **[atomic]** `runAtomic(&stop)` — проверка `atomic.Bool`; внешняя горутина ставит `stop=true`.
- **[timer]** `runTimer(..., d)` — прерывание по срабатыванию `time.NewTimer(d)`.
- **[signal]** `signal.NotifyContext` — `ctx` завершается по `os.Interrupt` (Ctrl-C) или вручную в конце.

---
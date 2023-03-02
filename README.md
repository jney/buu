# Buu

![buu](https://user-images.githubusercontent.com/747/222115773-82075b15-7e68-47fe-9d39-c7536ebcac53.png)

## Debouncer

The debouncer is partly based on https://github.com/bep/debounce.
`context.Context` was added in order to ensure passed function would be called even if the program is stopped

```go
debouncer := NewDebouncer(context.Background(), 80*time.Millisecond)
debouncer.Add(myFunc)
debouncer.Add(myFunc)
```

the debouncer **will run** the last function on `<-ctx.Done`

## Throttler

```go
throttler := NewThrottler(context.Background(), 80*time.Millisecond)
throttler.Add(myFunc)
throttler.Add(myFunc)
```

the throttler **will not run** the queued functions on `<-ctx.Done`

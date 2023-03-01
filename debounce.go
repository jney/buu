package buu

import (
	"context"
	"sync"
	"time"
)

type debouncer struct {
	mu    sync.Mutex
	after time.Duration
	done  chan struct{}
	timer *time.Timer
	fn    func()
	added bool
}

func NewDebouncer(ctx context.Context, after time.Duration) *debouncer {
	t := time.NewTimer(-1)
	t.Stop()
	d := &debouncer{after: after, done: make(chan struct{}, 1), timer: t}
	go d.run(ctx)
	return d
}

func (d *debouncer) run(ctx context.Context) {
	for {
		select {
		case <-d.timer.C:
			d.mu.Lock()
			d.fn()
			d.added = false
			d.mu.Unlock()
		case <-ctx.Done():
			d.mu.Lock()
			defer d.mu.Unlock()
			if d.added {
				d.fn()
			}
			d = nil
			return
		case <-d.done:
			d.mu.Lock()
			defer d.mu.Unlock()
			if d.added {
				d.fn()
			}
			d = nil
			return
		}
	}
}

func (d *debouncer) Stop() {
	d.done <- struct{}{}
}

func (d *debouncer) Add(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.added = true
	if d.timer != nil {
		d.timer.Reset(d.after)
	} else {
		d.timer.C = time.After(d.after)
	}
	d.fn = fn
}

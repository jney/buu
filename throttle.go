package buu

import (
	"context"
	"sync"
	"time"
)

type throttler struct {
	after   time.Duration
	done    chan struct{}
	first   chan struct{}
	fns     []func()
	mu      sync.Mutex
	started bool
	stop    chan struct{}
	timer   *time.Ticker
}

func NewThrottler(ctx context.Context, after time.Duration) *throttler {
	timer := time.NewTicker(after)
	timer.Stop()
	t := &throttler{
		after: after,
		done:  make(chan struct{}, 1),
		first: make(chan struct{}, 1), // run immediatly
		fns:   make([]func(), 0),
		stop:  make(chan struct{}, 1),
		timer: timer,
	}
	go t.run(ctx)
	return t
}

func (t *throttler) exec() bool {
	if len(t.fns) == 0 {
		return false
	}
	var fn func()
	fn, t.fns = t.fns[0], t.fns[1:]
	fn()
	return len(t.fns) != 0
}

func (t *throttler) run(ctx context.Context) {
	defer func() {
		t.done <- struct{}{}
	}()
	for {
		select {
		case <-t.first:
			t.mu.Lock()
			t.exec()
			t.mu.Unlock()
		case <-t.timer.C:
			t.mu.Lock()
			t.started = t.exec()
			if !t.started {
				t.timer.Stop()
			}
			t.mu.Unlock()
		case <-ctx.Done():
			t.mu.Lock()
			defer t.mu.Unlock()
			return
		case <-t.stop:
			t.mu.Lock()
			defer t.mu.Unlock()
			return
		}
	}
}

func (t *throttler) Done() <-chan struct{} { return t.done }

func (t *throttler) Stop() { t.stop <- struct{}{} }

func (t *throttler) Add(fn func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fns = append(t.fns, fn)
	if !t.started {
		t.timer.Reset(t.after)
		t.started = true
		t.first <- struct{}{}
	}
}

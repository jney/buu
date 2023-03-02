package buu

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func incrementer() (func() int, func()) {
	var counter uint64
	return func() int { return int(atomic.LoadUint64(&counter)) },
		func() { atomic.AddUint64(&counter, 1) }
}

func TestNewThrottler(t *testing.T) {
	ctx := context.Background()
	throttler := NewThrottler(ctx, 50*time.Millisecond)
	get, set := incrementer()
	go func() {
		for i := 0; i < 3; i++ {
			for j := 0; j < 10; j++ {
				throttler.Add(set)
				time.Sleep(10 * time.Millisecond)
			}
			// counter should be +1 only there
		}
	}()
	time.Sleep(310 * time.Millisecond)
	throttler.Stop()

	if get() != 7 {
		t.Errorf("debounced value should be 7 got %d", get())
		t.FailNow()
	}
}

func TestNewThrottlerWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	get, set := incrementer()
	throttler := NewThrottler(ctx, 20*time.Millisecond)
	go func() {
		for i := 0; i < 10; i++ {
			throttler.Add(set)
			time.Sleep(20 * time.Millisecond)
		}
	}()
	<-throttler.Done()

	if get() != 8 {
		t.Errorf("debounced value should be 8 got %d", get())
		t.FailNow()
	}
}

package buu

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewDebouncer(t *testing.T) {
	var counter uint64

	f := func() {
		atomic.AddUint64(&counter, 1)
	}
	ctx := context.Background()
	debouncer := NewDebouncer(ctx, 80*time.Millisecond)

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			debouncer.Add(f)
		}
		time.Sleep(100 * time.Millisecond)
		// counter should be +1 only there
	}

	c := int(atomic.LoadUint64(&counter))

	if c != 3 {
		t.Errorf("debounced value should be 3 got %d", c)
		t.FailNow()
	}
	debouncer.Stop()

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*190)
	defer cancel()
	counter = 0
	debouncer = NewDebouncer(ctx, 20*time.Millisecond)
	for i := 0; i < 10; i++ {
		debouncer.Add(f)
		time.Sleep(50 * time.Millisecond)
		// counter should be +1 only there
	}
	c = int(atomic.LoadUint64(&counter))
	time.Sleep(500 * time.Millisecond)

	if c != 4 {
		t.Errorf("debounced value should be 4 got %d", c)
		t.FailNow()
	}
	debouncer.Stop()
}

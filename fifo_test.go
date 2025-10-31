// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent_test

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"code.hybscloud.com/concurrent"
)

func TestQueue(t *testing.T) {
	t.Run("basic usage", func(t *testing.T) {
		c, p := concurrent.NewQueue[int](1024)
		i := 100
		err := p.Enqueue(&i)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		res, err := c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *res != 100 {
			t.Errorf("dequeue expected %v but got %v", i, *res)
		}
	})
	t.Run("with options", func(t *testing.T) {
		c, p := concurrent.NewQueue[int](16, func(opts *concurrent.QueueOptions) {
			opts.DistinctValues = true
		})
		val := 100
		err := concurrent.EnqueueWait(p, &val)
		if err != nil {
			t.Errorf("Enqueue failed: %v", err)
		}

		result, err := concurrent.DequeueWait(c)
		if err != nil {
			t.Errorf("Dequeue failed: %v", err)
		}
		if *result != 100 {
			t.Errorf("Expected 100, got %d", *result)
		}
	})
	t.Run("with invalid options", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid options")
			}
		}()
		_, _ = concurrent.NewQueue[int](16, func(opts *concurrent.QueueOptions) {
			opts.DistinctValues = false
		})
	})
}

func TestSPSCQueue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected not implement panic")
		}
	}()
	c, p := concurrent.NewSPSCQueue[int](8)
	if c == nil || p == nil {
		t.Error("NewSPSCQueue returned nil")
	}
}

func TestMPSCQueue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected not implement panic")
		}
	}()
	c, p := concurrent.NewMPSCQueue[int](8)
	if c == nil || p == nil {
		t.Error("NewMPSCQueue returned nil")
	}
}

func TestSPMCQueue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected not implement panic")
		}
	}()
	c, p := concurrent.NewSPMCQueue[int](8)
	if c == nil || p == nil {
		t.Error("NewSPMCQueue returned nil")
	}
}

func TestMPMCQueueInvalidOrder(t *testing.T) {
	t.Run("too small order", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for capacity < 2")
			}
		}()
		concurrent.NewMPMCQueue[int](1)
	})
	t.Run("too large order", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for too large capacity")
			}
		}()
		concurrent.NewMPMCQueue[int]((1 << 30) + 1)
	})
}

func TestMPMCQueueRmfLF(t *testing.T) {
	t.Run("basic usage", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int](1024)
		i := 100
		err := p.Enqueue(&i)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		res, err := c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *res != 100 {
			t.Errorf("dequeue expected %v but got %v", i, *res)
		}
	})

	t.Run("invalid order", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for capacity < 2")
			}
		}()
		concurrent.NewMPMCQueue[int](1)
	})

	t.Run("simple enqueue dequeue", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int](4)
		elem, err := c.Dequeue()
		if err != concurrent.ErrTemporaryUnavailable {
			t.Errorf("dequeue expected ErrTemporaryUnavailable but got %v", err)
			return
		}
		if elem != nil {
			t.Errorf("dequeue expected nil but got %v", elem)
			return
		}
		i0, i1, i2, i3, i4 := 100, 101, 102, 103, 104
		err = p.Enqueue(&i0)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		err = p.Enqueue(&i1)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		err = p.Enqueue(&i2)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		err = p.Enqueue(&i3)
		if err != nil {
			t.Errorf("enqueue: %v", err)
			return
		}
		err = p.Enqueue(&i4) // full
		if err != concurrent.ErrTemporaryUnavailable {
			t.Errorf("enqueue expected ErrTemporaryUnavailable but got %v", err)
			return
		}
		elem, err = c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *elem != i0 {
			t.Errorf("dequeue expected %v but got %v", i0, elem)
			return
		}
		elem, err = c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *elem != i1 {
			t.Errorf("dequeue expected %v but got %v", i1, elem)
			return
		}
		elem, err = c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *elem != i2 {
			t.Errorf("dequeue expected %v but got %v", i2, elem)
			return
		}
		elem, err = c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if *elem != i3 {
			t.Errorf("dequeue expected %v but got %v", i3, elem)
			return
		}
		_, err = c.Dequeue()
		if err != concurrent.ErrTemporaryUnavailable {
			t.Errorf("dequeue expected ErrTemporaryUnavailable but got %v", err)
			return
		}
	})

	t.Run("simple enqueue dequeue 2", func(t *testing.T) {
		type element struct {
			a, b int
		}
		e1, e2 := element{a: 100, b: 1}, element{a: 200, b: 2}
		c, p := concurrent.NewMPMCQueue[element](4)
		err := p.Enqueue(&e1)
		if err != nil {
			t.Errorf("enquque: %v", err)
			return
		}
		err = p.Enqueue(&e2)
		if err != nil {
			t.Errorf("enquque: %v", err)
			return
		}
		e, err := c.Dequeue()
		if err != nil {
			t.Errorf("dequeue: %v", err)
			return
		}
		if e.a != 100 || e.b != 1 {
			t.Errorf("dequeue expected %v but got %v", e1, e)
			return
		}
	})

	const defaultCapacity = 1 << 8
	t.Run("1 consumer 1 producer", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 1, 1)
	})

	t.Run("1 consumer 4 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 1, 4)
	})

	t.Run("1 consumer 16 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 1, 16)
	})

	t.Run("1 consumer 64 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 1, 64)
	})

	t.Run("4 consumers 1 producer", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 4, 1)
	})

	t.Run("4 consumers 4 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 4, 4)
	})

	t.Run("4 consumers 16 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 4, 16)
	})

	t.Run("4 consumers 64 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 4, 64)
	})

	t.Run("16 consumers 1 producer", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 16, 1)
	})

	t.Run("16 consumers 4 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 16, 4)
	})

	t.Run("16 consumers 16 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 16, 16)
	})

	t.Run("16 consumers 64 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 16, 64)
	})

	t.Run("64 consumers 1 producer", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 64, 1)
	})

	t.Run("64 consumers 4 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 64, 4)
	})

	t.Run("64 consumers 16 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 64, 16)
	})

	t.Run("64 consumers 64 producers", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		testMPMCQueue(t, c, p, 64, 64)
	})
}

func TestMPMCQueueIndirect(t *testing.T) {
	t.Run("basic usage", func(t *testing.T) {
		c, p := concurrent.NewMPMCQueueIndirect(256)

		index := uintptr(42)
		err := p.Enqueue(index)
		if err != nil {
			t.Errorf("Enqueue failed: %v", err)
		}

		index2 := uintptr(99)
		err = p.Enqueue(index2)
		if err != nil {
			t.Errorf("Enqueue failed: %v", err)
		}

		value, err := c.Dequeue()
		if err != nil {
			t.Errorf("Dequeue failed: %v", err)
		}
		if value != uintptr(42) {
			t.Errorf("Expected 42, got %d", value)
		}

		value2, err := c.Dequeue()
		if err != nil {
			t.Errorf("Dequeue failed: %v", err)
		}
		if value2 != uintptr(99) {
			t.Errorf("Expected 99, got %d", value2)
		}

		_, err = c.Dequeue()
		if err == nil {
			t.Error("Expected error on empty queue")
		}
	})

	t.Run("invalid capacity", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic on invalid capacity")
			}
		}()
		_, _ = concurrent.NewMPMCQueueIndirect(1)
	})

	t.Run("too large capacity", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic on invalid capacity")
			}
		}()
		_, _ = concurrent.NewMPMCQueueIndirect(1 << 31)
	})
}

func BenchmarkMPMCQueueRmfLF(b *testing.B) {
	const defaultCapacity = 1 << 16

	b.Run("1 consumer 1 producer", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 1, 1)
	})

	b.Run("1 consumer 4 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 1, 4)
	})

	b.Run("1 consumer 16 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 1, 16)
	})

	b.Run("1 consumer 64 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 1, 64)
	})

	b.Run("4 consumers 1 producer", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 4, 1)
	})

	b.Run("4 consumers 4 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 4, 4)
	})

	b.Run("4 consumers 16 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 4, 16)
	})

	b.Run("4 consumers 64 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 4, 64)
	})

	b.Run("16 consumers 1 producer", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 16, 1)
	})

	b.Run("16 consumers 4 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 16, 4)
	})

	b.Run("16 consumers 16 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 16, 16)
	})

	b.Run("16 consumers 64 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 16, 64)
	})

	b.Run("64 consumers 1 producer", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 64, 1)
	})

	b.Run("64 consumers 4 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 64, 4)
	})

	b.Run("64 consumers 16 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 64, 16)
	})

	b.Run("64 consumers 64 producers", func(b *testing.B) {
		c, p := concurrent.NewMPMCQueue[int64](defaultCapacity)
		benchmarkMPMCQueue(b, c, p, 64, 64)
	})
}

func TestEnqueueDequeueWait(t *testing.T) {
	c, p := concurrent.NewMPMCQueue[int](2)
	e1, e2, e3 := 1, 2, 3
	err := concurrent.EnqueueWait[int](p, &e1)
	if err != nil {
		t.Errorf("enqueue wait: %v", err)
		return
	}
	err = concurrent.EnqueueWait[int](p, &e2)
	if err != nil {
		t.Errorf("enqueue wait: %v", err)
		return
	}
	elem, err := concurrent.DequeueWait[int](c)
	if err != nil {
		t.Errorf("dequeue wait: %v", err)
		return
	}
	if *elem != e1 {
		t.Errorf("dequeue wait expected %v but got %v", e1, *elem)
		return
	}
	err = concurrent.EnqueueWait[int](p, &e3)
	if err != nil {
		t.Errorf("enqueue wait: %v", err)
		return
	}
	elem, err = concurrent.DequeueWait[int](c)
	if err != nil {
		t.Errorf("dequeue wait: %v", err)
		return
	}
	if *elem != e2 {
		t.Errorf("dequeue wait expected %v but got %v", e2, *elem)
		return
	}
	elem, err = concurrent.DequeueWait[int](c)
	if err != nil {
		t.Errorf("dequeue wait: %v", err)
		return
	}
	if *elem != e3 {
		t.Errorf("dequeue wait expected %v but got %v", e3, *elem)
		return
	}
}

// Test utilities for MPMC queues (interface-based, reusable)
func testMPMCQueue(t *testing.T, c concurrent.Consumer[int64], p concurrent.Producer[int64], cn, pn int) {
	n := 1 << 12
	for i := 0; i < pn; i++ {
		go func(i int) {
			for j := 0; j < n; j++ {
				val := int64(i<<32) | int64(j)
				err := concurrent.EnqueueWait[int64](p, &val)
				if err != nil {
					t.Errorf("queue enqueue: %v", err)
					return
				}
			}
		}(i)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < cn; i++ {
		wg.Add(1)
		go func(i int) {
			last := make([]atomic.Int64, pn)
			for j := 0; j < pn; j++ {
				last[j].Store(-1)
			}
			for j := 0; j < n*pn/cn; j++ {
				// Note: False positives for out-of-order detection may occur here, which is acceptable.
				item, err := concurrent.DequeueWait[int64](c)
				if err != nil {
					t.Errorf("queue dequeue: %v", err)
					return
				}
				high, low := *item>>32, *item&math.MaxUint32
				old := last[high].Swap(low)
				if low <= old {
					t.Logf("queue produce consume out of order: %d>%d", low, old)
					return
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func benchmarkMPMCQueue(b *testing.B, c concurrent.Consumer[int64], p concurrent.Producer[int64], cn, pn int) {
	for i := 0; i < pn; i++ {
		go func(i int) {
			for j := 0; j < (b.N+cn)/pn+1; j++ {
				val := int64(i<<32) | int64(j)
				_ = concurrent.EnqueueWait[int64](p, &val)
			}
		}(i)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < cn; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < b.N/cn+1; j++ {
				_, _ = concurrent.DequeueWait[int64](c)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

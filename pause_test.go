// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestPause(t *testing.T) {
	// Test default (20 cycles)
	Pause()

	// Test single cycle
	Pause(1)

	// Test multiple iterations
	for i := 0; i < 100; i++ {
		Pause()
	}
}

func TestPauseWithCycles(t *testing.T) {
	// Test various cycle counts
	testCases := []int{1, 10, 20, 50, 100}

	for _, cycles := range testCases {
		t.Run(fmt.Sprintf("cycles=%d", cycles), func(t *testing.T) {
			// Should not crash or hang
			Pause(cycles)
		})
	}
}

func TestPauseDefault(t *testing.T) {
	// Test default (should be 20 cycles)
	Pause()
}

func TestPauseInSpinLoop(t *testing.T) {
	var counter atomic.Int32
	done := make(chan struct{})

	// Goroutine that increments counter after a delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		counter.Store(1)
		close(done)
	}()

	// Spin-wait with default Pause
	for counter.Load() == 0 {
		Pause()
	}

	<-done
	if counter.Load() != 1 {
		t.Errorf("expected counter to be 1, got %d", counter.Load())
	}
}

func TestPauseSingleCycleInLoop(t *testing.T) {
	var counter atomic.Int32
	done := make(chan struct{})

	// Goroutine that increments counter after a delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		counter.Store(1)
		close(done)
	}()

	// Spin-wait with single cycle Pause
	for counter.Load() == 0 {
		Pause(1)
	}

	<-done
	if counter.Load() != 1 {
		t.Errorf("expected counter to be 1, got %d", counter.Load())
	}
}

func BenchmarkPause1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pause(1)
	}
}

func BenchmarkPause10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pause(10)
	}
}

func BenchmarkPause20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pause(20)
	}
}

func BenchmarkPauseDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pause()
	}
}

func BenchmarkSpinLoopWithPauseDefault(b *testing.B) {
	var counter atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Store(0)
			go func() {
				runtime.Gosched()
				counter.Store(1)
			}()

			// Spin-wait with default Pause
			for counter.Load() == 0 {
				Pause()
			}
		}
	})
}

func BenchmarkSpinLoopWithPause1(b *testing.B) {
	var counter atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Store(0)
			go func() {
				runtime.Gosched()
				counter.Store(1)
			}()

			// Spin-wait with single cycle Pause
			for counter.Load() == 0 {
				Pause(1)
			}
		}
	})
}

func BenchmarkSpinLoopWithoutPause(b *testing.B) {
	var counter atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Store(0)
			go func() {
				runtime.Gosched()
				counter.Store(1)
			}()

			// Spin-wait without Pause (busy wait)
			for counter.Load() == 0 {
				// Busy wait
			}
		}
	})
}

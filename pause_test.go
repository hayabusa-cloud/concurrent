// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent_test

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"code.hybscloud.com/concurrent"
)

func TestPause(t *testing.T) {
	// Test default (20 cycles)
	concurrent.Pause()

	// Test a single cycle
	concurrent.Pause(1)

	// Test multiple iterations
	for i := 0; i < 100; i++ {
		concurrent.Pause()
	}
}

func TestPauseWithCycles(t *testing.T) {
	// Test various cycle counts
	testCases := []int{-1, 1, 10, 20, 50, 100}

	for _, cycles := range testCases {
		t.Run(fmt.Sprintf("cycles=%d", cycles), func(t *testing.T) {
			// Should not crash or hang
			concurrent.Pause(cycles)
		})
	}
}

func TestPauseDefault(t *testing.T) {
	// Test default (should be 20 cycles)
	concurrent.Pause()
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
		concurrent.Pause()
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
		concurrent.Pause(1)
	}

	<-done
	if counter.Load() != 1 {
		t.Errorf("expected counter to be 1, got %d", counter.Load())
	}
}

func BenchmarkPause(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause()
		}
	})
	b.Run("1 cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause(1)
		}
	})
	b.Run("10 cycles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause(10)
		}
	})
	b.Run("20 cycles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause(20)
		}
	})
	b.Run("50 cycles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause(50)
		}
	})
	b.Run("100 cycles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concurrent.Pause(100)
		}
	})
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
				concurrent.Pause()
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
				concurrent.Pause(1)
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

// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestSpinLock_BasicLockUnlock tests basic lock and unlock operations
func TestSpinLock_BasicLockUnlock(t *testing.T) {
	var sl SpinLock

	// Lock and unlock should not panic
	sl.Lock()
	sl.Unlock()

	// Should be able to lock again after unlock
	sl.Lock()
	sl.Unlock()
}

// TestSpinLock_MutualExclusion tests that the lock provides mutual exclusion
func TestSpinLock_MutualExclusion(t *testing.T) {
	var sl SpinLock
	var counter int64
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				sl.Lock()
				counter++
				sl.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	if counter != expected {
		t.Errorf("Expected counter to be %d, but got %d", expected, counter)
	}
}

// TestSpinLock_ConcurrentAccess tests concurrent lock/unlock operations
func TestSpinLock_ConcurrentAccess(t *testing.T) {
	var sl SpinLock
	var wg sync.WaitGroup
	var activeHolders atomic.Int32
	var maxConcurrent atomic.Int32

	numGoroutines := 50
	iterations := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				sl.Lock()

				// Track how many goroutines hold the lock simultaneously
				current := activeHolders.Add(1)
				if current > maxConcurrent.Load() {
					maxConcurrent.Store(current)
				}

				// Simulate some work
				runtime.Gosched()

				activeHolders.Add(-1)
				sl.Unlock()
			}
		}()
	}

	wg.Wait()

	// Only one goroutine should hold the lock at a time
	if maxConcurrent.Load() > 1 {
		t.Errorf("Expected at most 1 concurrent lock holder, but got %d", maxConcurrent.Load())
	}
}

// TestSpinLock_HighContention tests the lock under high contention
func TestSpinLock_HighContention(t *testing.T) {
	var sl SpinLock
	var wg sync.WaitGroup
	var operations atomic.Int64

	numGoroutines := runtime.NumCPU() * 2
	duration := 100 * time.Millisecond
	done := make(chan struct{})

	time.AfterFunc(duration, func() {
		close(done)
	})

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					sl.Lock()
					operations.Add(1)
					sl.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if operations.Load() == 0 {
		t.Error("Expected some operations to complete, but got 0")
	}
}

// TestSpinLock_MultipleLocksUnlocks tests multiple sequential lock/unlock cycles
func TestSpinLock_MultipleLocksUnlocks(t *testing.T) {
	var sl SpinLock

	iterations := 10000
	for i := 0; i < iterations; i++ {
		sl.Lock()
		sl.Unlock()
	}
}

// TestSpinLock_WithDeferredUnlock tests using defer for unlock
func TestSpinLock_WithDeferredUnlock(t *testing.T) {
	var sl SpinLock
	var counter int

	func() {
		sl.Lock()
		defer sl.Unlock()
		counter++
	}()

	// Should be able to lock again after the deferred unlock
	sl.Lock()
	counter++
	sl.Unlock()

	if counter != 2 {
		t.Errorf("Expected counter to be 2, but got %d", counter)
	}
}

// TestSpinLock_ProtectedSection tests that critical section is properly protected
func TestSpinLock_ProtectedSection(t *testing.T) {
	var sl SpinLock
	var wg sync.WaitGroup
	var sharedSlice []int

	numGoroutines := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(val int) {
			defer wg.Done()
			sl.Lock()
			// Critical section: non-atomic slice operation
			sharedSlice = append(sharedSlice, val)
			sl.Unlock()
		}(i)
	}

	wg.Wait()

	if len(sharedSlice) != numGoroutines {
		t.Errorf("Expected slice length to be %d, but got %d", numGoroutines, len(sharedSlice))
	}
}

// BenchmarkSpinLock_LockUnlock benchmarks basic lock/unlock operations
func BenchmarkSpinLock_LockUnlock(b *testing.B) {
	var sl SpinLock

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Lock()
			sl.Unlock()
		}
	})
}

// BenchmarkSpinLock_WithWork benchmarks lock/unlock with simulated work
func BenchmarkSpinLock_WithWork(b *testing.B) {
	var sl SpinLock
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Lock()
			counter++
			sl.Unlock()
		}
	})
}

// BenchmarkSpinLock_HighContention benchmarks under high contention
func BenchmarkSpinLock_HighContention(b *testing.B) {
	var sl SpinLock
	var counter int64

	b.SetParallelism(runtime.NumCPU() * 2)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Lock()
			counter++
			runtime.Gosched()
			sl.Unlock()
		}
	})
}

// BenchmarkSpinLock_LowContention benchmarks under low contention
func BenchmarkSpinLock_LowContention(b *testing.B) {
	var sl SpinLock

	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Lock()
			sl.Unlock()
		}
	})
}

// BenchmarkSpinLock_vs_Mutex compares SpinLock with standard Mutex
func BenchmarkSpinLock_vs_Mutex(b *testing.B) {
	b.Run("SpinLock", func(b *testing.B) {
		var sl SpinLock
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sl.Lock()
				sl.Unlock()
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		var mu sync.Mutex
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				mu.Unlock()
			}
		})
	})
}

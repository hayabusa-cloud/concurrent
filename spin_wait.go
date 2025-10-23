// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"runtime"
	"time"
	_ "unsafe"
)

// SpinWait is a lightweight synchronization type that
// you can use in low-level scenarios with lower cost.
// The zero value for SpinWait is ready to use.
type SpinWait struct {
	counter, n uint32
}

// Once performs a single spin
func (s *SpinWait) Once() {
	s.counter++
	if s.WillYield() {
		s.n++
		runtime.Gosched()
		return
	}
	pauseN(defaultPauseCycles)
}

// WillYield returns true if calling SpinOnce() will result
// in occurring a thread sleeping instead of a simply procyield()
func (s *SpinWait) WillYield() bool {
	if (s.counter+1)&(1<<(4-min(4, s.n>>2))-1) != 0 {
		return false
	}

	return true
}

// Reset resets the counter in SpinWait
func (s *SpinWait) Reset() {
	s.counter = 0
	s.n = 0
}

var yieldDuration = 250 * time.Microsecond

// Yield yields the processor by sleeping for a duration based on the backoff level lv (default: 1).
// Higher levels sleep longer with quadratic scaling: Yield(1), Yield(2)=4x, Yield(3)=9x, etc.
// For automatic adaptive backoff in tight loops, use SpinWait instead.
func Yield(lv ...int) {
	d := yieldDuration

	if len(lv) > 0 {
		d = time.Duration(lv[0]*lv[0]) * yieldDuration
	}
	if d > 0 {
		time.Sleep(d)
		return
	}
	runtime.Gosched()
}

// SetYieldDuration sets the base duration unit for Yield().
// Recommended: 50-250 microseconds for real-time systems, 1-4 ms for general workloads.
func SetYieldDuration(d time.Duration) {
	yieldDuration = max(0, d)
}

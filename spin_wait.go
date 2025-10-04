// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"runtime"
	"time"
	_ "unsafe"
)

const (
	spinWaitSleepDuration = 100 * time.Microsecond
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
		time.Sleep(spinWaitSleepDuration)
		return
	}
	runtime.Gosched()
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

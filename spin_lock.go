// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"runtime"
	"sync/atomic"
)

type SpinLock struct {
	_ noCopy
	n atomic.Uintptr
}

func (sl *SpinLock) Lock() {
	for {
		n := sl.n.Add(1)
		if n < 2 {
			return
		} else if n < 4 {
			pauseN(defaultPauseCycles)
			continue
		}
		runtime.Gosched()
	}
}

func (sl *SpinLock) Unlock() {
	sl.n.Store(0)
}

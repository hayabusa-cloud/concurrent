// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
)

type poster struct {
	entries  []uintptr
	order    int
	capacity uint64
	offers   atomic.Uint64
	_        cpu.CacheLinePad
	polls    atomic.Uint64
}

const posterNilFlag = 1 << 63

func newPoster(order int) *poster {
	if order < 1 || order > 24 {
		panic("bad capacity order")
	}
	ret := &poster{
		order:    order,
		capacity: 1 << order,
	}
	ret.entries = make([]uintptr, ret.capacity)
	nil0 := uintptr(posterNilFlag | 0)
	for i := 0; i < len(ret.entries); i++ {
		ret.entries[i] = nil0
	}

	return ret
}

func (po *poster) offer(elem uintptr) bool {
	for {
		o, p := po.offers.Load(), po.polls.Load()
		if o != po.offers.Load() {
			continue
		}
		if o == p+po.capacity {
			return false
		}
		i := o & (po.capacity - 1)
		round := (o >> po.order) & (posterNilFlag - 1)
		success := atomic.CompareAndSwapUintptr(&po.entries[i], posterNilFlag|uintptr(round), elem)
		po.offers.CompareAndSwap(o, o+1)

		if success {
			return true
		}
	}
}

func (po *poster) poll() (elem uintptr, ok bool) {
	for {
		p, o := po.polls.Load(), po.offers.Load()
		i := p & (po.capacity - 1)
		e := po.entries[i]
		if p != po.polls.Load() {
			continue
		}
		if p == o {
			return 0, false
		}
		nextRound := uintptr((p>>po.order)+1) & (posterNilFlag - 1)
		if e == posterNilFlag|nextRound {
			po.polls.CompareAndSwap(p, p+1)
			continue
		}
		success := atomic.CompareAndSwapUintptr(&po.entries[i], e, posterNilFlag|nextRound)
		po.polls.CompareAndSwap(p, p+1)
		if success {
			return e, true
		}
	}
}

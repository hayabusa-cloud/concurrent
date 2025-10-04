// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"sync/atomic"

	"golang.org/x/sys/cpu"
)

type rmfLF struct {
	entries   []uintptr
	order     int
	capacity  uint64
	indexSkip uint64
	offers    atomic.Uint64
	_         cpu.CacheLinePad
	polls     atomic.Uint64
}

const (
	rmfLFNilFlag    = 1 << 63
	rmfLFModuleBit  = 6
	rmfLFModuleMask = (1 << rmfLFModuleBit) - 1
)

func newRmfLF(order int) *rmfLF {
	if order < 1 || order > 30 {
		panic("bad capacity order")
	}
	ret := &rmfLF{
		order:    order,
		capacity: 1 << order,
	}
	ret.indexSkip = 1 << max(0, ret.order-rmfLFModuleBit)

	ret.entries = make([]uintptr, ret.capacity)
	nil0 := uintptr(rmfLFNilFlag | 0)
	for i := 0; i < len(ret.entries); i++ {
		ret.entries[i] = nil0
	}

	return ret
}

func (lf *rmfLF) offer(elem uintptr) bool {
	for {
		o, p := lf.offers.Load(), lf.polls.Load()
		if o != lf.offers.Load() {
			continue
		}
		if o == p+lf.capacity {
			return false
		}
		i := o & (lf.capacity - 1)
		entry := lf.entry(i)
		round := (o >> lf.order) & (rmfLFNilFlag - 1)
		success := atomic.CompareAndSwapUintptr(&lf.entries[entry], rmfLFNilFlag|uintptr(round), elem)
		lf.offers.CompareAndSwap(o, o+1)

		if success {
			return true
		}
	}
}

func (lf *rmfLF) poll() (elem uintptr, ok bool) {
	for {
		p, o := lf.polls.Load(), lf.offers.Load()
		i := p & (lf.capacity - 1)
		if p != lf.polls.Load() {
			continue
		}
		if p == o {
			return 0, false
		}
		entry := lf.entry(i)
		e := lf.entries[entry]
		nextRound := uintptr((p>>lf.order)+1) & (rmfLFNilFlag - 1)
		if e == rmfLFNilFlag|nextRound {
			lf.polls.CompareAndSwap(p, p+1)
			continue
		}
		success := atomic.CompareAndSwapUintptr(&lf.entries[entry], e, rmfLFNilFlag|nextRound)
		lf.polls.CompareAndSwap(p, p+1)
		if success {
			return e, true
		}
	}
}

func (lf *rmfLF) entry(index uint64) uint64 {
	p, q := index>>rmfLFModuleBit, index&rmfLFModuleMask
	return q*lf.indexSkip + p
}

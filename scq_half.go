// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"golang.org/x/sys/cpu"
)

// half-address size version
type scqHalf struct {
	entries   []uint64
	_         cpu.CacheLinePad
	head      uint32
	_         cpu.CacheLinePad
	tail      uint32
	_         cpu.CacheLinePad
	threshold int32

	n      uint32
	module uint32
}

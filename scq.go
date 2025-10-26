// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"golang.org/x/sys/cpu"
)

// double word CAS version
type scq struct {
	entries   [][]uint64
	_         cpu.CacheLinePad
	head      uint64
	_         cpu.CacheLinePad
	tail      uint64
	_         cpu.CacheLinePad
	threshold int64

	n      uint64
	module uint64
}

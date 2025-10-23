// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import "unsafe"

// DoubleUint64 creates and returns []uint64 with size 2 and the given values
// the address of u128 will be 16-bytes aligned
func DoubleUint64(first, second uint64) (u128 []uint64) {
	mem := make([]uint64, 3)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(mem)))
	off := ptr - ptr&^uintptr(0xf)
	u128 = mem[off/8 : off/8+2]
	u128[0] = first
	u128[1] = second

	return
}

// DoubleUintPtr creates and returns []uintptr with size 2 and the given values
// the address of dw will be 16-bytes aligned
func DoubleUintPtr(first, second uintptr) (dw []uintptr) {
	u128 := DoubleUint64(uint64(first), uint64(second))
	addr := unsafe.Pointer(&u128[0])

	return unsafe.Slice((*uintptr)(addr), 2)
}

type noCopy struct{}

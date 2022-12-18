// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"testing"
	"unsafe"
)

func TestDoubleUint64(t *testing.T) {
	for i := 0; i < 1<<16; i++ {
		u128 := DoubleUint64(1, 2)
		addr := unsafe.Pointer(&u128[0])
		if uintptr(addr)&0xf > 0 {
			t.Errorf("address of u128 not 16-bytes aligned")
			return
		}
		if u128[0] != 1 || u128[1] != 2 {
			t.Errorf("bad u128 value expected %v but got %v", [2]uint64{1, 2}, u128)
			return
		}
	}
}

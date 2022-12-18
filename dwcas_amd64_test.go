// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64

package concurrent

import (
	"testing"
	"unsafe"
)

func TestCas128(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dw := DoubleUint64(2, 4)
		ptr := (*uint64)(unsafe.Pointer(&dw[0]))
		ok := cas128(ptr, [2]uint64{2, 4}, [2]uint64{3, 6})
		if !ok {
			t.Errorf("double word cas expected swapped but failed")
			return
		}
		if dw[0] != 3 || dw[1] != 6 {
			t.Errorf("double word cas expected new val %v but got %v", [2]uint64{3, 6}, dw)
			return
		}
		ok = cas128(ptr, [2]uint64{2, 4}, [2]uint64{3, 6})
		if ok {
			t.Errorf("double word cas expected failed but swapped")
			return
		}
	})

	t.Run("failed", func(t *testing.T) {
		dw := DoubleUint64(2, 4)
		ptr := (*uint64)(unsafe.Pointer(&dw[0]))
		ok := cas128(ptr, [2]uint64{1, 4}, [2]uint64{3, 6})
		if ok {
			t.Errorf("double word cas expected failed but swapped")
			return
		}
		ok = cas128(ptr, [2]uint64{2, 3}, [2]uint64{3, 6})
		if ok {
			t.Errorf("double word cas expected failed but swapped")
			return
		}
		ok = cas128(ptr, [2]uint64{1, 3}, [2]uint64{3, 6})
		if ok {
			t.Errorf("double word cas expected failed but swapped")
			return
		}
		if dw[0] != 2 || dw[1] != 4 {
			t.Errorf("double word cas expected val %v but got %v", [2]uint64{2, 4}, dw)
			return
		}
	})
}

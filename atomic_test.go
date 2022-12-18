// ©Hayabusa Cloud Co., Ltd. 2022. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"testing"
)

func TestCompareAndSwapUint128(t *testing.T) {
	u128 := DoubleUint64(2, 4)
	newVal := [2]uint64{5, 10}
	ok := CompareAndSwapUint128(&u128[0], [2]uint64{1, 4}, newVal)
	if ok {
		t.Errorf("compare and swap uint128 expected failed but swapped")
		return
	}
	if u128[0] != 2 || u128[1] != 4 {
		t.Errorf("expected u128 with value %v but got %v", [2]uint64{2, 4}, u128)
		return
	}
	ok = CompareAndSwapUint128(&u128[0], [2]uint64{2, 3}, newVal)
	if ok {
		t.Errorf("compare and swap uint128 expected failed but swapped")
		return
	}
	if u128[0] != 2 || u128[1] != 4 {
		t.Errorf("expected u128 with value %v but got %v", [2]uint64{2, 4}, u128)
		return
	}
	ok = CompareAndSwapUint128(&u128[0], [2]uint64{2, 4}, newVal)
	if !ok {
		t.Errorf("compare and swap uint128 expected swapped but failed")
		return
	}
	if u128[0] != newVal[0] || u128[1] != newVal[1] {
		t.Errorf("expected u128 with value %v but got %v", newVal, u128)
		return
	}
}

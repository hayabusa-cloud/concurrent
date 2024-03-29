// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent_test

import (
	"hybscloud.com/concurrent"
	"testing"
)

func TestCompareAndSwapUint128(t *testing.T) {
	u128 := concurrent.DoubleUint64(2, 4)
	newVal := [2]uint64{5, 10}
	ok := concurrent.CompareAndSwapUint128(&u128[0], [2]uint64{1, 4}, newVal)
	if ok {
		t.Errorf("compare and swap uint128 expected failed but swapped")
		return
	}
	if u128[0] != 2 || u128[1] != 4 {
		t.Errorf("expected u128 with value %v but got %v", [2]uint64{2, 4}, u128)
		return
	}
	ok = concurrent.CompareAndSwapUint128(&u128[0], [2]uint64{2, 3}, newVal)
	if ok {
		t.Errorf("compare and swap uint128 expected failed but swapped")
		return
	}
	if u128[0] != 2 || u128[1] != 4 {
		t.Errorf("expected u128 with value %v but got %v", [2]uint64{2, 4}, u128)
		return
	}
	ok = concurrent.CompareAndSwapUint128(&u128[0], [2]uint64{2, 4}, newVal)
	if !ok {
		t.Errorf("compare and swap uint128 expected swapped but failed")
		return
	}
	if u128[0] != newVal[0] || u128[1] != newVal[1] {
		t.Errorf("expected u128 with value %v but got %v", newVal, u128)
		return
	}
}

func TestAndUint64(t *testing.T) {
	u64 := uint64(0xffffffff)
	concurrent.AndUint64(&u64, 0xff0ff0ff)
	if u64 != 0xff0ff0ff {
		t.Errorf("AndUint64 expected 0xff0ff0ff but got %x", u64)
		return
	}
	concurrent.AndUint64(&u64, 0x00f00f00)
	if u64 != 0 {
		t.Errorf("AndUint64 expected 0 but got %x", u64)
		return
	}
	concurrent.AndUint64(&u64, 0xffffffff)
	if u64 != 0 {
		t.Errorf("AndUint64 expected 0 but got %x", u64)
		return
	}
}

func TestOrUint64(t *testing.T) {
	u64 := uint64(0)
	concurrent.OrUint64(&u64, 0xff0ff0ff)
	if u64 != 0xff0ff0ff {
		t.Errorf("OrUint64 expected 0xff0ff0ff but got %x", u64)
		return
	}
	concurrent.OrUint64(&u64, 0x00f00f00)
	if u64 != 0xffffffff {
		t.Errorf("OrUint64 expected 0xffffffff but got %x", u64)
		return
	}
	concurrent.OrUint64(&u64, 0)
	if u64 != 0xffffffff {
		t.Errorf("OrUint64 expected 0xffffffff but got %x", u64)
		return
	}
}

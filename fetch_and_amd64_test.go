// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64 || arm64

package concurrent

import (
	"testing"
)

func TestAnd64(t *testing.T) {
	u64 := uint64(0xffffffff)
	and64(&u64, 0xffffffff)
	if u64 != 0xffffffff {
		t.Errorf("fetch and expected 0xffffffff but got %x", u64)
		return
	}
	and64(&u64, 0xff0ff0ff)
	if u64 != 0xff0ff0ff {
		t.Errorf("fetch and expected 0xff0ff0ff but got %x", u64)
		return
	}
	and64(&u64, 0x00f00f00)
	if u64 != 0 {
		t.Errorf("fetch and expected 0 but got %x", u64)
		return
	}
	and64(&u64, 0xffffffff)
	if u64 != 0 {
		t.Errorf("fetch and expected 0 but got %x", u64)
		return
	}
}

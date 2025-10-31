// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent_test

import (
	"testing"
	"time"

	"code.hybscloud.com/concurrent"
)

func TestSpinWait(t *testing.T) {
	sw := concurrent.SpinWait{}
	for i := 0; i < 15; i++ {
		if sw.WillYield() {
			t.Errorf("spin wait expected will not yield but got will yield")
		}
		sw.Once()
	}
	if !sw.WillYield() {
		t.Errorf("spin wait expected will yield but got will not yield")
	}
	sw.Reset()
	for i := 0; i < 15; i++ {
		if sw.WillYield() {
			t.Errorf("spin wait expected will not yield but got will yield")
		}
		sw.Once()
	}
	sw.Reset()
	if sw.WillYield() {
		t.Errorf("spin wait expected will not yield but got will yield")
	}
}

func TestYield(t *testing.T) {
	concurrent.Yield()
	concurrent.Yield(-1)
	concurrent.Yield(20)
}

func TestSetYieldDuration(t *testing.T) {
	concurrent.SetYieldDuration(100 * time.Millisecond)
	concurrent.Yield()
	concurrent.SetYieldDuration(time.Duration(0))
	concurrent.Yield(100)
	concurrent.SetYieldDuration(-time.Millisecond)
	concurrent.Yield()
}

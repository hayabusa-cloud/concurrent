// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import "testing"

func TestSpinWait(t *testing.T) {
	sw := SpinWait{}
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

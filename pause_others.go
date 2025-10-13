// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64

package concurrent

import "runtime"

func pause1() {
	runtime.Gosched()
}

func pauseN(cycles int) {
	for i := 0; i < cycles; i++ {
		runtime.Gosched()
	}
}

func pause(cycles int) {
	if cycles == 1 {
		pause1()
	} else {
		pauseN(cycles)
	}
}

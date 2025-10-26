// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

const defaultPauseCycles = 20

// Pause executes CPU pause instructions to reduce energy consumption in spin-wait loops.
//
// Defaults to 20 cycles if not specified. Uses optimized assembly on amd64/arm64.
//
// Usage:
//
//	Pause()     // 20 cycles (default)
//	Pause(1)    // 1 cycle
//	Pause(40)   // 40 cycles
func Pause(cycles ...int) {
	n := defaultPauseCycles
	if len(cycles) > 0 && cycles[0] > 0 {
		n = cycles[0]
	}
	pause(n)
}

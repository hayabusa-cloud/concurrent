// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64 || arm64

package concurrent

//go:noescape
func pause1()

//go:noescape
func pauseN(cycles int)

func pause(cycles int) {
	if cycles == 1 {
		pause1()
	} else {
		pauseN(cycles)
	}
}

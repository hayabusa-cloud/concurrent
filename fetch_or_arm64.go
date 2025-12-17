// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

package concurrent

import "sync/atomic"

func or64(ptr *uint64, val uint64) {
	atomic.OrUint64(ptr, val)
}

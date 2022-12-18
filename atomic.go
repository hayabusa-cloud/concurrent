// ©Hayabusa Cloud Co., Ltd. 2022. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

// CompareAndSwapUint128 atomically compares u128 with old,
// and if they're equal, swaps value of u128 with new.
// It reports whether the swap ran.
var CompareAndSwapUint128 func(u128 *uint64, old, new [2]uint64) bool = cas128

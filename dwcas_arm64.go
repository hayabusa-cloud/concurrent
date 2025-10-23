// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

package concurrent

//go:noescape
func cas128(ptr *uint64, old, new [2]uint64) bool

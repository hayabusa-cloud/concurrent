// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64
// +build amd64

#include "textflag.h"

// void Or64(addr *uint64, v uint64)
// Atomically:
//   *addr |= v
TEXT ·or64(SB), NOSPLIT, $0-16
	MOVQ	ptr+0(FP), BX
	MOVQ	val+8(FP), AX
	LOCK
	ORQ AX, (BX)
	RET

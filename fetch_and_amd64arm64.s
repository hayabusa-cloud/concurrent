// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// void And64(ptr *uint64, val uint64)
// Atomically:
//   *ptr &= val
TEXT ·and64(SB), NOSPLIT, $0-16
	MOVQ	ptr+0(FP), BX
	MOVQ	val+8(FP), AX
	LOCK
	ANDQ    AX, (BX)
	RET

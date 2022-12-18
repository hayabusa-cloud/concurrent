// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// https://cdrdv2.intel.com/v1/dl/getContent/671200
// Vol.2A 3-207~3-209
// bool	·Cas128(uint64[2] ptr, uint64[2] old, uint64[2] new)
// Atomically:
//	if(val[0] == old[0] && val[1] == old[1]){
//		val[0] = new[0];
//		val[1] = new[1];
//		return 1;
//	} else {
//		return 0;
//	}
TEXT ·cas128(SB), NOSPLIT, $0-41
    MOVQ    ptr+0(FP), BP
    MOVQ    old+8(FP), AX
    MOVQ    old+16(FP), DX
    MOVQ    new+24(FP), BX
    MOVQ    new+32(FP), CX
    LOCK
    CMPXCHG16B (BP)
    SETEQ   swapped+40(FP)
    RET

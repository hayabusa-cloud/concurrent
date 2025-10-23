// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// https://developer.arm.com/documentation/ddi0596/latest/
// CASPAL (Compare and Swap Pair with Acquire-release semantics)
// bool	·cas128(uint64[2] ptr, uint64[2] old, uint64[2] new)
// Atomically:
//	if(val[0] == old[0] && val[1] == old[1]){
//		val[0] = new[0];
//		val[1] = new[1];
//		return 1;
//	} else {
//		return 0;
//	}
// Uses CASPAL for sequential consistency (equivalent to x86 LOCK CMPXCHG16B)
TEXT ·cas128(SB), NOSPLIT, $0-41
	MOVD	ptr+0(FP), R8         // R8 = ptr (address must be 16-byte aligned)
	MOVD	old+8(FP), R0         // R0 = old[0] (Rs must be even register)
	MOVD	old+16(FP), R1        // R1 = old[1] (Rs+1)
	MOVD	new+24(FP), R2        // R2 = new[0] (Rt must be even register)
	MOVD	new+32(FP), R3        // R3 = new[1] (Rt+1)

	// CASPAL: Compare and Swap Pair with Acquire-release + Load semantics
	// After CASPAL, if comparison fails, Rs and Rs+1 contain the actual value from memory
	CASPAL	(R0, R1), (R2, R3), (R8)

	// Compare if the swap succeeded
	// If swap succeeded, R0:R1 still equals old values
	// If swap failed, R0:R1 now contains the actual memory value
	MOVD	old+8(FP), R4         // Reload old[0] into R4
	CMP	R0, R4                    // Compare R0 with old[0]
	BNE	fail
	MOVD	old+16(FP), R4        // Reload old[1] into R4
	CMP	R1, R4                    // Compare R1 with old[1]
	BNE	fail

	// Success
	MOVD	$1, R5
	MOVB	R5, swapped+40(FP)
	RET

fail:
	// Failed
	MOVB	ZR, swapped+40(FP)
	RET

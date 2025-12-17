// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// bool ·cas128(ptr *uint64, old [2]uint64, new [2]uint64) bool
//
// AcqRel semantics (success):
// - Release: order prior writes before the conditional store of the CAS.
// - Acquire: order subsequent reads/writes after a successful CAS only.
TEXT ·cas128(SB), NOSPLIT, $0-41
	MOVD	ptr+0(FP), R8
	MOVD	old+8(FP), R0
	MOVD	old+16(FP), R1
	MOVD	new+24(FP), R2
	MOVD	new+32(FP), R3

	// Preserve expected pair (CASPD clobbers R0:R1 with memory value on failure).
	MOVD	R0, R4
	MOVD	R1, R5

	// Release ordering for prior writes (ISHST).
	// Required because we can't encode release directly in CASP* mnemonics here.
	DMB	$0xA

	// Pair CAS: if *R8 == R0:R1 then *R8 = R2:R3 else R0:R1 = *R8
	CASPD	(R0, R1), (R8), (R2, R3)

	// swapped = (*observed == expected)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	cas128_fail

	// Acquire ordering on success only (ISHLD).
	DMB	$0x9

	MOVD	$1, R6
	MOVB	R6, swapped+40(FP)
	RET

cas128_fail:
	MOVB	ZR, swapped+40(FP)
	RET

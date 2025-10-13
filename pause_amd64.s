// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT, $0-0
	PAUSE
	RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT, $0-8
	MOVQ cycles+0(FP), CX
	TESTQ CX, CX
	JLE done
loop:
	PAUSE
	DECQ CX
	JNZ loop
done:
	RET

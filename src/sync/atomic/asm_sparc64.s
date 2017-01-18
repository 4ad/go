// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "asm_sparc64.h"

TEXT ·SwapInt32(SB),NOSPLIT|NOFRAME,$0-20
	JMP	·SwapUint32(SB)

TEXT ·SwapUint32(SB),NOSPLIT,$0-20
	MOVD	addr+0(FP), I3
	MOVUW	new+8(FP), I1
again:
	MEM_SYNC
	MOVUW	(I3), I5
	CASW	(I3), I5, I1
	CMP	I1, I5
	BNEW	again
	MEM_SYNC
	MOVUW	I5, old+16(FP)
	RET

TEXT ·SwapInt64(SB),NOSPLIT|NOFRAME,$0-24
	JMP	·SwapUint64(SB)

TEXT ·SwapUint64(SB),NOSPLIT,$0-24
	MOVD	addr+0(FP), I3
	MOVD	new+8(FP), I1
again:
	MEM_SYNC
	MOVD	(I3), I5
	CASD	(I3), I5, I1
	CMP	I1, I5
	BNED	again
	MEM_SYNC
	MOVD	I5, old+16(FP)
	RET

TEXT ·SwapUintptr(SB),NOSPLIT|NOFRAME,$0-24
	JMP	·SwapUint64(SB)

TEXT ·CompareAndSwapInt32(SB),NOSPLIT|NOFRAME,$0-17
	JMP	·CompareAndSwapUint32(SB)

TEXT ·CompareAndSwapUint32(SB),NOSPLIT,$0-17
	MOVD	addr+0(FP), I1
	MOVUW	old+8(FP), I3
	MOVUW	new+12(FP), I5
	MEM_SYNC
	CASW	(I1), I3, I5
	CMP	I5, I3
	MOVD	$0, I3
	MOVE	ICC, $1, I3
	MEM_SYNC
	MOVB	I3, swapped+16(FP)
	RET

TEXT ·CompareAndSwapUintptr(SB),NOSPLIT|NOFRAME,$0-25
	JMP	·CompareAndSwapUint64(SB)

TEXT ·CompareAndSwapInt64(SB),NOSPLIT|NOFRAME,$0-25
	JMP	·CompareAndSwapUint64(SB)

TEXT ·CompareAndSwapUint64(SB),NOSPLIT,$0-25
	MOVD	addr+0(FP), I1
	MOVD	old+8(FP), I3
	MOVD	new+16(FP), I5
	MEM_SYNC
	CASD	(I1), I3, I5
	CMP	I5, I3
	MOVD	$0, I3
	MOVE	XCC, $1, I3
	MEM_SYNC
	MOVB	I3, swapped+24(FP)
	RET

TEXT ·AddInt32(SB),NOSPLIT|NOFRAME,$0-20
	JMP	·AddUint32(SB)

TEXT ·AddUint32(SB),NOSPLIT,$0-20
	MOVD	addr+0(FP), I4
	MOVUW	delta+8(FP), I3
	MOVUW	(I4), I1
	MEM_SYNC
retry:
	ADD	I1, I3, I5
	CASW	(I4), I1, I5
	CMP	I1, I5
	MOVNE	ICC, I5, I1
	BNEW	retry
	ADD	I1, I3, I5
	MEM_SYNC
	MOVUW	I5, new+16(FP)
	RET

TEXT ·AddUintptr(SB),NOSPLIT|NOFRAME,$0-24
	JMP	·AddUint64(SB)

TEXT ·AddInt64(SB),NOSPLIT|NOFRAME,$0-24
	JMP	·AddUint64(SB)

TEXT ·AddUint64(SB),NOSPLIT,$0-24
	MOVD	addr+0(FP), I4
	MOVD	delta+8(FP), I3
	MEM_SYNC
	MOVD	(I4), I1
retry:
	ADD	I1, I3, I5
	CASD	(I4), I1, I5
	CMP	I1, I5
	MOVNE	XCC, I5, I1
	BNED	retry
	ADD	I1, I3, I5
	MEM_SYNC
	MOVD	I5, new+16(FP)
	RET

TEXT ·LoadInt32(SB),NOSPLIT|NOFRAME,$0-12
	JMP	·LoadUint32(SB)

TEXT ·LoadUint32(SB),NOSPLIT,$0-12
	MOVD	addr+0(FP), I1
	MEM_SYNC
	LDUW	(I1), I1
	MEM_SYNC
	MOVUW	I1, val+8(FP)
	RET

TEXT ·LoadInt64(SB),NOSPLIT|NOFRAME,$0-16
	JMP	·LoadUint64(SB)

TEXT ·LoadUint64(SB),NOSPLIT,$0-16
	MOVD	addr+0(FP), I1
	MEM_SYNC
	LDD	(I1), I1
	MEM_SYNC
	MOVD	I1, val+8(FP)
	RET

TEXT ·LoadUintptr(SB),NOSPLIT|NOFRAME,$0-16
	JMP	·LoadPointer(SB)

TEXT ·LoadPointer(SB),NOSPLIT|NOFRAME,$0-16
	JMP	·LoadUint64(SB)

TEXT ·StoreInt32(SB),NOSPLIT|NOFRAME,$0-12
	JMP	·StoreUint32(SB)

TEXT ·StoreUint32(SB),NOSPLIT,$0-12
	MOVD	addr+0(FP), I3
	MOVUW	val+8(FP), I5
	MEM_SYNC
	STW	I5, (I3)
	MEM_SYNC
	RET

TEXT ·StoreInt64(SB),NOSPLIT|NOFRAME,$0-16
	JMP	·StoreUint64(SB)

TEXT ·StoreUint64(SB),NOSPLIT,$0-16
	MOVD	addr+0(FP), I3
	MOVD	val+8(FP), I5
	MEM_SYNC
	STD	I5, (I3)
	MEM_SYNC
	RET

TEXT ·StoreUintptr(SB),NOSPLIT|NOFRAME,$0-16
	JMP	·StoreUint64(SB)

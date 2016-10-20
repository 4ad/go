// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "asm_sparc64.h"

// bool cas(uint32 *ptr, uint32 old, uint32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT runtime∕internal∕atomic·Cas(SB), NOSPLIT, $0-17
	MOVD	ptr+0(FP), I1
	MOVUW	old+8(FP), I3
	MOVUW	new+12(FP), I5
	MEM_SYNC
	CASW	(I1), I3, I5
	CMP	I5, I3
	MOVD	$0, I3
	MOVE	ICC, $1, I3
	MEM_SYNC
	MOVB	I3, ret+16(FP)
	RET

// bool	runtime∕internal∕atomic·Cas64(uint64 *ptr, uint64 old, uint64 new)
// Atomically:
//	if(*val == *old){
//		*val = new;
//		return 1;
//	} else {
//		return 0;
//	}
TEXT runtime∕internal∕atomic·Cas64(SB), NOSPLIT, $0-25
	MOVD	ptr+0(FP), I1
	MOVD	old+8(FP), I3
	MOVD	new+16(FP), I5
	MEM_SYNC
	CASD	(I1), I3, I5
	CMP	I5, I3
	MOVD	$0, I3
	MOVE	XCC, $1, I3
	MEM_SYNC
	MOVB	I3, ret+24(FP)
	RET

TEXT runtime∕internal∕atomic·Casuintptr(SB), NOSPLIT|NOFRAME, $0-25
	JMP	runtime∕internal∕atomic·Cas64(SB)

TEXT runtime∕internal∕atomic·Loaduintptr(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Load64(SB)

TEXT runtime∕internal∕atomic·Loaduint(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Load64(SB)

TEXT runtime∕internal∕atomic·Storeuintptr(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT runtime∕internal∕atomic·Xadduintptr(SB), NOSPLIT|NOFRAME, $0-24
	JMP	runtime∕internal∕atomic·Xadd64(SB)

TEXT runtime∕internal∕atomic·Loadint64(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Load64(SB)

TEXT runtime∕internal∕atomic·Xaddint64(SB), NOSPLIT|NOFRAME, $0-24
	JMP	runtime∕internal∕atomic·Xadd64(SB)

// bool casp(void **val, void *old, void *new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT runtime∕internal∕atomic·Casp1(SB), NOSPLIT|NOFRAME, $0-25
	JMP runtime∕internal∕atomic·Cas64(SB)

// uint32 xadd(uint32 volatile *ptr, int32 delta)
// Atomically:
//	*val += delta;
//	return *val;
TEXT runtime∕internal∕atomic·Xadd(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), I4
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
	MOVUW	I5, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xadd64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), I4
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
	MOVD	I5, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), I3
	MOVUW	new+8(FP), I1
again:
	MEM_SYNC
	MOVUW	(I3), I5
	CASW	(I3), I5, I1
	CMP	I1, I5
	BNEW	again
	MEM_SYNC
	MOVUW	I5, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), I3
	MOVD	new+8(FP), I1
again:
	MEM_SYNC
	MOVD	(I3), I5
	CASD	(I3), I5, I1
	CMP	I1, I5
	BNED	again
	MEM_SYNC
	MOVD	I5, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchguintptr(SB), NOSPLIT|NOFRAME, $0-24
	JMP	runtime∕internal∕atomic·Xchg64(SB)


// TODO(shawn): verify this is performed without a write barrier;
// see #15270.
TEXT runtime∕internal∕atomic·StorepNoWB(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT runtime∕internal∕atomic·Store(SB), NOSPLIT, $0-12
	MOVD	ptr+0(FP), I3
	MOVUW	val+8(FP), I5
	MEM_SYNC
	STW	I5, (I3)
	MEM_SYNC
	RET

TEXT runtime∕internal∕atomic·Store64(SB), NOSPLIT, $0-16
	MOVD	ptr+0(FP), I3
	MOVD	val+8(FP), I5
	MEM_SYNC
	STD	I5, (I3)
	MEM_SYNC
	RET

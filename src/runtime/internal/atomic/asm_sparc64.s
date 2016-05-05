// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// bool cas(uint32 *ptr, uint32 old, uint32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT runtime∕internal∕atomic·Cas(SB), NOSPLIT, $0-17
	MOVD	ptr+0(FP), R3
	MOVUW	old+8(FP), R1
	MOVUW	new+12(FP), R2
	MEMBAR	$15
	CASW	(R3), R1, R2
	CMP	R2, R1
	MOVD	$0, R1
	MOVE	ICC, $1, R1
	MEMBAR	$15
	MOVB	R1, ret+16(FP)
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
	MOVD	ptr+0(FP), R3
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	MEMBAR	$15
	CASD	(R3), R1, R2
	CMP	R2, R1
	MOVD	$0, R1
	MOVE	XCC, $1, R1
	MEMBAR	$15
	MOVB	R1, ret+24(FP)
	RET

TEXT runtime∕internal∕atomic·Casuintptr(SB), NOSPLIT|NOFRAME, $0-25
	JMP	runtime∕internal∕atomic·Cas64(SB)

TEXT runtime∕internal∕atomic·Loaduintptr(SB),  NOSPLIT|NOFRAME, $0-16
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
	MOVD	ptr+0(FP), R4
	MOVUW	delta+8(FP), R1
	MOVUW	(R4), R3
	MEMBAR	$15
retry:
	ADD	R3, R1, R2
	CASW	(R4), R3, R2
	CMP	R3, R2
	MOVNE	ICC, R2, R3
	BNEW	retry
	ADD	R3, R1, R2
	MEMBAR	$15
	MOVUW	R2, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xadd64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R4
	MOVD	delta+8(FP), R1
	MEMBAR	$15
	MOVD	(R4), R3
retry:
	ADD	R3, R1, R2
	CASD	(R4), R3, R2
	CMP	R3, R2
	MOVNE	XCC, R2, R3
	BNED	retry
	ADD	R3, R1, R2
	MEMBAR	$15
	MOVD	R2, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), R1
	MOVUW	new+8(FP), R3
again:
	MEMBAR	$15
	MOVUW	(R1), R2
	CASW	(R1), R2, R3
	CMP	R3, R2
	BNEW	again
	MEMBAR	$15
	MOVUW	R2, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R1
	MOVD	new+8(FP), R3
again:
	MEMBAR	$15
	MOVD	(R1), R2
	CASD	(R1), R2, R3
	CMP	R3, R2
	BNED	again
	MEMBAR	$15
	MOVD	R2, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchguintptr(SB), NOSPLIT|NOFRAME, $0-24
	JMP	runtime∕internal∕atomic·Xchg64(SB)


TEXT runtime∕internal∕atomic·Storep1(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT runtime∕internal∕atomic·Store(SB), NOSPLIT, $0-12
	MOVD	ptr+0(FP), R1
	MOVUW	val+8(FP), R2
	MEMBAR	$12
	STW	R2, (R1)
	MEMBAR	$10
	RET

TEXT runtime∕internal∕atomic·Store64(SB), NOSPLIT, $0-16
	MOVD	ptr+0(FP), R1
	MOVD	val+8(FP), R2
	MEMBAR	$12
	STD	R2, (R1)
	MEMBAR	$10
	RET

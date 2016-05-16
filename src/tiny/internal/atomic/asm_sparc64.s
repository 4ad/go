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
	MOVD	ptr+0(FP), R25
	MOVUW	old+8(FP), R27
	MOVUW	new+12(FP), R29
	MEMBAR	$15
	CASW	(R25), R27, R29
	CMP	R29, R27
	MOVD	$0, R27
	MOVE	ICC, $1, R27
	MEMBAR	$15
	MOVB	R27, ret+16(FP)
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
	MOVD	ptr+0(FP), R25
	MOVD	old+8(FP), R27
	MOVD	new+16(FP), R29
	MEMBAR	$15
	CASD	(R25), R27, R29
	CMP	R29, R27
	MOVD	$0, R27
	MOVE	XCC, $1, R27
	MEMBAR	$15
	MOVB	R27, ret+24(FP)
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
	MOVD	ptr+0(FP), R28
	MOVUW	delta+8(FP), R27
	MOVUW	(R28), R25
	MEMBAR	$15
retry:
	ADD	R25, R27, R29
	CASW	(R28), R25, R29
	CMP	R25, R29
	MOVNE	ICC, R29, R25
	BNEW	retry
	ADD	R25, R27, R29
	MEMBAR	$15
	MOVUW	R29, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xadd64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R28
	MOVD	delta+8(FP), R27
	MEMBAR	$15
	MOVD	(R28), R25
retry:
	ADD	R25, R27, R29
	CASD	(R28), R25, R29
	CMP	R25, R29
	MOVNE	XCC, R29, R25
	BNED	retry
	ADD	R25, R27, R29
	MEMBAR	$15
	MOVD	R29, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg(SB), NOSPLIT, $0-20
	MOVD	ptr+0(FP), R27
	MOVUW	new+8(FP), R25
again:
	MEMBAR	$15
	MOVUW	(R27), R29
	CASW	(R27), R29, R25
	CMP	R25, R29
	BNEW	again
	MEMBAR	$15
	MOVUW	R29, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchg64(SB), NOSPLIT, $0-24
	MOVD	ptr+0(FP), R27
	MOVD	new+8(FP), R25
again:
	MEMBAR	$15
	MOVD	(R27), R29
	CASD	(R27), R29, R25
	CMP	R25, R29
	BNED	again
	MEMBAR	$15
	MOVD	R29, ret+16(FP)
	RET

TEXT runtime∕internal∕atomic·Xchguintptr(SB), NOSPLIT|NOFRAME, $0-24
	JMP	runtime∕internal∕atomic·Xchg64(SB)


TEXT runtime∕internal∕atomic·Storep1(SB), NOSPLIT|NOFRAME, $0-16
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT runtime∕internal∕atomic·Store(SB), NOSPLIT, $0-12
	MOVD	ptr+0(FP), R27
	MOVUW	val+8(FP), R29
	MEMBAR	$12
	STW	R29, (R27)
	MEMBAR	$10
	RET

TEXT runtime∕internal∕atomic·Store64(SB), NOSPLIT, $0-16
	MOVD	ptr+0(FP), R27
	MOVD	val+8(FP), R29
	MEMBAR	$12
	STD	R29, (R27)
	MEMBAR	$10
	RET

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "asm_sparc64.h"

// uint32 runtime∕internal∕atomic·Load(uint32 volatile* addr)
TEXT ·Load(SB),NOSPLIT,$0-12
	MOVD	ptr+0(FP), I1
	MEM_SYNC
	LDUW	(I1), I1
	MEM_SYNC
	MOVUW	I1, ret+8(FP)
	RET

// uint64 runtime∕internal∕atomic·Load64(uint64 volatile* addr)
TEXT ·Load64(SB),NOSPLIT,$0-16
	MOVD	ptr+0(FP), I1
	MEM_SYNC
	LDD	(I1), I1
	MEM_SYNC
	MOVD	I1, ret+8(FP)
	RET

// void *runtime∕internal∕atomic·Loadp(void *volatile *addr)
TEXT ·Loadp(SB),NOSPLIT|NOFRAME,$0-16
	JMP	runtime∕internal∕atomic·Load64(SB)

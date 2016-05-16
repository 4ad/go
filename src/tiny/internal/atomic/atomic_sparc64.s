// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// uint32 runtime∕internal∕atomic·Load(uint32 volatile* addr)
TEXT ·Load(SB),NOSPLIT|NOFRAME,$0-12
	MOVD	ptr+0(FP), R25
	MEMBAR	$3
	LDUW	(R25), R25
	MEMBAR	$5
	MOVUW	R25, ret+8(FP)
	RET

// uint64 runtime∕internal∕atomic·Load64(uint64 volatile* addr)
TEXT ·Load64(SB),NOSPLIT|NOFRAME,$0-16
	MOVD	ptr+0(FP), R25
	MEMBAR	$3
	LDD	(R25), R25
	MEMBAR	$5
	MOVD	R25, ret+8(FP)
	RET

// void *runtime∕internal∕atomic·Loadp(void *volatile *addr)
TEXT ·Loadp(SB),NOSPLIT|NOFRAME,$0-16
	JMP	runtime∕internal∕atomic·Load64(SB)

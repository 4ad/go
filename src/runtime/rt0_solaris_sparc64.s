// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "asm_sparc64.h"

TEXT _rt0_sparc64_solaris(SB),NOSPLIT|NOFRAME,$0
	MOVD	WINDOW_SIZE+0(BSP), O0 // argc
	MOVD	$(WINDOW_SIZE+8)(BSP), O1 // argv
	MOVD	$main(SB), O3
	JMPL	O3, ZR

TEXT main(SB),NOSPLIT|REGWIN,$0
	MOVW	I0, O0
	MOVD	I1, O1
	CALL	runtime·rt0_go(SB)
	RET

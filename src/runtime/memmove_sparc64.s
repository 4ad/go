// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// void runtime·memmove(void*, void*, uintptr)
TEXT runtime·memmove(SB), NOSPLIT, $-8-24
	// TODO(aram):
	MOVD	$50, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

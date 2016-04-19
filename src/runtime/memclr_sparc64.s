// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT runtime·memclr(SB),NOSPLIT|NOFRAME,$0-16
	// TODO(aram):
	MOVD	$60, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

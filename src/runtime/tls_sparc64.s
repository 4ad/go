// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

// save_g saves the g register into pthread-provided
// thread-local memory, so that we can call externally compiled
// arm64 code that will overwrite this register.
//
// If !iscgo, this is a no-op.
//
// NOTE: setg_gcc<> assume this clobbers only TMP.
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVB	runtime·iscgo(SB), TMP
	CMP	TMP, ZR
	BEW	nocgo

	MOVD	$runtime·tls_g(SB), TMP
	MOVD	g, (TMP)
nocgo:
	RET

// load_g loads the g register from pthread-provided
// thread-local memory, for use after calling externally compiled
// sparc64 code that overwrote those registers.
//
// This is never called directly from C code (it doesn't have to
// follow the C ABI), but it may be called from a C context, where the
// usual Go registers aren't set up.
//
// NOTE: _cgo_topofstack assumes this only clobbers g, and TMP.
TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
	MOVB	runtime·iscgo(SB), TMP
	CMP	TMP, ZR
	BEW	nocgo

	MOVD	$runtime·tls_g(SB), TMP
	MOVD	(TMP), g
nocgo:
	RET

GLOBL runtime·tls_g+0(SB), TLSBSS, $8

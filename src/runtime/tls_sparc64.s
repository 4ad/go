// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"
#include "asm_sparc64.h"

// save_g saves the g register into pthread-provided
// thread-local memory, so that we can call externally compiled
// sparc64 code that will overwrite this register.
//
// If !iscgo, this is a no-op.
//
// NOTE: setg_gcc<> assume this clobbers only RT1.
TEXT runtime·save_g(SB),NOSPLIT|NOFRAME,$0-0
// On Solaris we always use TLS, even without cgo.
#ifndef GOOS_solaris
	MOVB	runtime·iscgo(SB), RT1
	CMP	RT1, ZR
	BEW	nocgo
#endif

	MOVD	$runtime·tls_g(SB), RT1
	MOVD	g, (RT1)
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
// NOTE: _cgo_topofstack assumes this only clobbers g, and RT1.
TEXT runtime·load_g(SB),NOSPLIT|NOFRAME,$0-0
// On Solaris we always use TLS, even without cgo.
#ifndef GOOS_solaris
	MOVB	runtime·iscgo(SB), RT1
	CMP	RT1, ZR
	BEW	nocgo
#endif

	MOVD	$runtime·tls_g(SB), RT1
	MOVD	(RT1), g
nocgo:
	RET

TEXT runtime·do_cgo_init(SB),NOSPLIT,$0-0
	// if there is a _cgo_init, call it using the gcc ABI.
	MOVD	_cgo_init(SB), O4
	CMP	ZR, O4
	BED	nocgo

	MOVD	TLS, O3			// arg 3: TLS base pointer
	MOVD	$runtime·tls_g(SB), O2 	// arg 2: &tls_g
	MOVD	$setg_gcc<>(SB), O1	// arg 1: setg
	MOVD	g, O0			// arg 0: G
	CALL	(O4)
	MOVD	$runtime·g0(SB), g
nocgo:
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$16
	MOVD	O0, g
	MOVD	RT1, savedRT1-8(SP)
	CALL	runtime·save_g(SB)
	MOVD	savedRT1-8(SP), RT1
	RET

GLOBL runtime·tls_g+0(SB), TLSBSS, $8

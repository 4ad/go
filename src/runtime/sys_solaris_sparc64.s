// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for SPARC64, Solaris.
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"
#include "asm_sparc64.h"

// void libc_miniterrno(void *(*___errno)(void));
//
// Set the TLS errno pointer in M.
//
// Called using runtime·asmcgocall from os_solaris.c:/minit.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·miniterrno(SB),NOSPLIT,$0
	// asmcgocall will put first argument into I0.
	CALL	I0	// SysV ABI so returns in O0
	CALL	runtime·load_g(SB)
	MOVD	g_m(g), I3
	MOVD	O0,	(m_mOS+mOS_perrno)(I3)
	RET

// int64 runtime·nanotime1(void);
//
// clock_gettime(3c) wrapper because Timespec is too large for
// runtime·nanotime stack.
//
// Called using runtime·sysvicall6 from os_solaris.c:/nanotime.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·nanotime1(SB),NOSPLIT,$64
	MOVW	$3, O0	// CLOCK_REALTIME from <sys/time_impl.h>
	MOVD	$-16(BFP), O1
	MOVD	$libc_clock_gettime(SB), I3
	CALL	I3
	MOVD	-16(BFP), I3	// tv_sec from struct timespec
	MOVD	$1000000000, I1
	MULD	I1, I3	// multiply into nanoseconds
	MOVD	-8(BFP), I5	// tv_nsec, offset should be stable.
	ADD	I5, I3, I0
	RET

// pipe(3c) wrapper that returns fds in AX, DX.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·pipe1(SB),NOSPLIT,$16
	MOVD	$FIXED_FRAME(BSP), O0
	MOVD	$libc_pipe(SB), I3
	CALL	I3
	MOVW	(FIXED_FRAME+0)(BSP), I0
	MOVW	(FIXED_FRAME+4)(BSP), I1
	RET

// Call a library function with SysV calling conventions.
// The called function can take a maximum of 6 INTEGER class arguments,
// see 
// 	SYSTEM V APPLICATION BINARY INTERFACE
// 	SPARC Version 9 Processor Supplement
// section 3.2.2.
//
// Called by runtime·asmcgocall or runtime·cgocall.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·asmsysvicall6(SB),NOSPLIT,$0
	// asmcgocall will put first argument into I0.
	MOVD	I0, L7
	MOVD	libcall_fn(I0), I3
	MOVD	libcall_args(I0), L1
	MOVD	libcall_n(I0), L2

	CMP	ZR, g
	BED	skiperrno1
	MOVD	g_m(g), I5
	MOVD	(m_mOS+mOS_perrno)(I5), I1
	CMP	I1, ZR
	BED	skiperrno1
	MOVW	ZR, (I1)

skiperrno1:
	CMP	L1, ZR
	BED	skipargs
	// Load 6 args into correspondent registers.
	MOVD	0(L1), O0
	MOVD	8(L1), O1
	MOVD	16(L1), O2
	MOVD	24(L1), O3
	MOVD	32(L1), O4
	MOVD	40(L1), O5
skipargs:

	MOVD	g, L1
	// Call SysV function
	CALL	I3
	MOVD	L1, g

	// Return result
	MOVD	O0, libcall_r1(L7)
	MOVD	O1, libcall_r2(L7)

	CMP	g, ZR
	BED	skiperrno2
	MOVD	g_m(g), I5
	MOVD	(m_mOS+mOS_perrno)(I5), I1
	CMP	I1, ZR
	BED	skiperrno2
	MOVW	(I1), I4
	MOVD	I4, libcall_err(L7)

skiperrno2:	
	RET

// uint32 tstart_sysvicall(M *newm);
TEXT runtime·tstart_sysvicall(SB),NOSPLIT,$0
	// I0 contains first arg newm
	MOVD	m_g0(I0), g		// g
	MOVD	I0, g_m(g)

	CALL	runtime·save_g(SB)

	// Layout new m scheduler stack on os stack.
	MOVD	BSP, I3
	MOVD	I3, (g_stack+stack_hi)(g)
	SUB	$(0x100000), I3		// stack size
	MOVD	I3, (g_stack+stack_lo)(g)
	ADD	$const__StackGuard, I3
	MOVD	I3, g_stackguard0(g)
	MOVD	I3, g_stackguard1(g)

	CALL	runtime·stackcheck(SB)
	CALL	runtime·mstart(SB)

	MOVW	ZR, ret+8(FP)
	RET

// Careful, this is called by __sighndlr, a libc function. We must preserve
// registers as per AMD 64 ABI.
TEXT runtime·sigtramp(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$75, I3
	ADD	$'!', I3, I3
	MOVB	I3, dbgbuf(SB)
	MOVD	$2, O0
	MOVD	$dbgbuf(SB), O1
	MOVD	$2, O2
	MOVD	$libc_write(SB), I3
	CALL	I3
	UNDEF
	RET

// Runs on OS stack, called from runtime·usleep1_go.
TEXT runtime·usleep2(SB),NOSPLIT,$0
	MOVW	usec+0(FP), O0
	MOVD	$libc_usleep(SB), I3
	CALL	I3
	RET

// Runs on OS stack, called from runtime·osyield.
TEXT runtime·osyield1(SB),NOSPLIT,$0
	MOVD	$libc_sched_yield(SB), I3
	CALL	I3
	RET

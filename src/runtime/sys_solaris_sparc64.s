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
TEXT runtime·miniterrno(SB),NOSPLIT|REGWIN,$0
	// asmcgocall will put first argument into I0.
	MOVD	I0, I1
	CALL	I1	// SysV ABI so returns in O0
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
TEXT runtime·nanotime1(SB),NOSPLIT|REGWIN,$64
	MOVW	$3, O0	// CLOCK_REALTIME from <sys/time_impl.h>
	MOVD	$tv-16(SP), O1
	MOVD	$libc_clock_gettime(SB), I3
	CALL	I3
	MOVD	tv_sec-16(SP), I3	// tv_sec from struct timespec
	MOVD	$1000000000, I1
	MULD	I1, I3	// multiply into nanoseconds
	MOVD	tv_nsec-8(SP), I5	// tv_nsec, offset should be stable.
	ADD	I5, I3, I0
	RET

// pipe(3c) wrapper that returns fds in AX, DX.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·pipe1(SB),NOSPLIT|REGWIN,$16
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
TEXT runtime·asmsysvicall6(SB),NOSPLIT|REGWIN,$0
	// asmcgocall will put first argument into I0.
	MOVD	I0, L6
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
	MOVD	O0, libcall_r1(L6)
	MOVD	O1, libcall_r2(L6)
	MOVD	O0, I0
	MOVD	O1, I1

	CMP	g, ZR
	BED	skiperrno2
	MOVD	g_m(g), I5
	MOVD	(m_mOS+mOS_perrno)(I5), I4
	CMP	I4, ZR
	BED	skiperrno2
	MOVW	(I4), I4
	MOVD	I4, libcall_err(L6)

skiperrno2:	
	RET

// uint32 tstart_sysvicall(M *newm);
TEXT runtime·tstart_sysvicall(SB),NOSPLIT|REGWIN,$0
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

#define SIGTRAMP_FRAME 192

// Careful, this is called by __sighndlr, a libc function.
// We must preserve registers as per SPARC64 ABI.
TEXT runtime·sigtramp(SB),NOSPLIT|REGWIN,$SIGTRAMP_FRAME
	CMP	g, ZR
	BNED	allgood
	MOVD	I0, (FIXED_FRAME+0)(BSP)
	CALL	runtime·badsignal(SB)
	JMP	exit

allgood:
	// save g
	MOVD	g, (-8-0*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)

	// save registers
	MOVD	CTXT, (-8-1*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I0, (-8-2*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I1, (-8-3*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I2, (-8-4*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I3, (-8-5*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I4, (-8-6*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	I5, (-8-7*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)

	// Save m->libcall and m->scratch. We need to do this because we
	// might get interrupted by a signal in runtime·asmcgocall.

	// save m->libcall 
	MOVD	g_m(g), L1
	MOVD	$m_libcall(L1), L2
	MOVD	libcall_fn(L2), L3
	MOVD	L3, (-8-8*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	libcall_args(L2), L3
	MOVD	L3, (-8-9*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	libcall_n(L2), L3
	MOVD	L3, (-8-10*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	libcall_r1(L2), L3
	MOVD	L3, (-8-11*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	libcall_r2(L2), L3
	MOVD	L3, (-8-12*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)

	// save m->scratch
	MOVD	$(m_mOS+mOS_scratch)(L1), L2
	MOVD	0(L2), L3
	MOVD	L3, (-8-13*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	8(L2), L3
	MOVD	L3, (-8-14*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	16(L2), L3
	MOVD	L3, (-8-15*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	24(L2), L3
	MOVD	L3, (-8-16*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	32(L2), L3
	MOVD	L3, (-8-17*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)
	MOVD	40(L2), L3
	MOVD	L3, (-8-18*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)

	// save errno, it might be EINTR; stuff we do here might reset it.
	MOVD	(m_mOS+mOS_perrno)(L1), L2
	MOVW	0(L2), L3
	MOVD	L3, (-8-19*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP)

	MOVD	g, I3
	// g = m->gsignal
	MOVD	m_gsignal(L1), L4
	MOVD	L4, g
	CALL	runtime·save_g(SB)

	// TODO: If current SP is not in gsignal.stack, then adjust.

	// prepare call
	MOVW	I0, (8*0+FIXED_FRAME)(BSP)
	MOVD	I1, (8*1+FIXED_FRAME)(BSP)
	MOVD	I2, (8*2+FIXED_FRAME)(BSP)
	MOVD	I3, (8*3+FIXED_FRAME)(BSP)
	CALL	runtime·sighandler(SB)

	// restore errno
	MOVD	g_m(g), L1
	MOVD	(m_mOS+mOS_perrno)(L1), L2
	MOVD	(-8-19*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVW	L3, 0(L2)

	// restore scratch
	MOVD	$(m_mOS+mOS_scratch)(L1), L2
	MOVD	(-8-13*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 0(L2)
	MOVD	(-8-14*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 8(L2)
	MOVD	(-8-15*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 16(L2)
	MOVD	(-8-16*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 24(L2)
	MOVD	(-8-17*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 32(L2)
	MOVD	(-8-18*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, 40(L2)

	// restore libcall
	MOVD	$m_libcall(L1), L2
	MOVD	(-8-8*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, libcall_fn(L2)
	MOVD	(-8-9*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, libcall_args(L2)
	MOVD	(-8-10*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3	
	MOVD	L3, libcall_n(L2)
	MOVD	(-8-11*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, libcall_r1(L2)
	MOVD	(-8-12*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), L3
	MOVD	L3, libcall_r2(L2)

	// restore registers
	MOVD	(-8-1*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), CTXT
	MOVD	(-8-2*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I0
	MOVD	(-8-3*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I1
	MOVD	(-8-4*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I2
	MOVD	(-8-5*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I3
	MOVD	(-8-6*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I4
	MOVD	(-8-7*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), I5

	// restore g
	MOVD	(-8+0*8+SIGTRAMP_FRAME+FIXED_FRAME)(BSP), g
	CALL	runtime·save_g(SB)

exit:
	RET


// Runs on OS stack, called from runtime·usleep1_go.
TEXT runtime·usleep2(SB),NOSPLIT|REGWIN,$0
	MOVW	us+0(FP), O0
	MOVD	$libc_usleep(SB), I3
	CALL	I3
	RET

// Runs on OS stack, called from runtime·osyield.
TEXT runtime·osyield1(SB),NOSPLIT|REGWIN,$0
	MOVD	$libc_sched_yield(SB), I3
	CALL	I3
	RET

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

DATA dbgbuf(SB)/8, $"\n\n"
GLOBL dbgbuf(SB), $8

TEXT runtime·rt0_go(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$1, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

DATA	runtime·mainPC+0(SB)/8,$runtime·main(SB)
GLOBL	runtime·mainPC(SB),RODATA,$8

TEXT runtime·breakpoint(SB),NOSPLIT,$-8-0
	TA	$0x81
	RET

TEXT runtime·asminit(SB),NOSPLIT,$-8-0
	RET

TEXT runtime·reginit(SB),NOSPLIT,$-8-0
	// TODO(aram): do we need to initialize FP registers?
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), NOSPLIT, $-8-8
	MOVD	buf+0(FP), R3
	MOVD	RSP, R1
	MOVD	R1, gobuf_sp(R3)
	MOVD	LR, gobuf_pc(R3)
	MOVD	g, gobuf_g(R3)
	MOVD	ZR, gobuf_lr(R3)
	MOVD	ZR, gobuf_ret(R3)
	MOVD	ZR, gobuf_ctxt(R3)
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT, $-8-8
	MOVD	buf+0(FP), R5
	MOVD	gobuf_g(R5), g
	CALL	runtime·save_g(SB)

	MOVD	0(g), R4	// make sure g is not nil
	MOVD	gobuf_sp(R5), R1
	MOVD	R1, RSP
	MOVD	gobuf_lr(R5), LR
	MOVD	gobuf_ret(R5), R1
	MOVD	gobuf_ctxt(R5), CTXT
	MOVD	ZR, gobuf_sp(R5)
	MOVD	ZR, gobuf_ret(R5)
	MOVD	ZR, gobuf_lr(R5)
	MOVD	ZR, gobuf_ctxt(R5)
	CMP	ZR, ZR // set condition codes for == test, needed by stack split
	MOVD	gobuf_pc(R5), R6
	JMPL	R6, ZR

// void mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT, $-8-8
	// Save caller state in g->sched
	MOVD	RSP, (g_sched+gobuf_sp)(g)
	MOVD	LR, (g_sched+gobuf_pc)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	g, (g_sched+gobuf_g)(g)

	// Switch to m->g0 & its stack, call fn.
	MOVD	g, R3
	MOVD	g_m(g), R8
	MOVD	m_g0(R8), g
	CALL	runtime·save_g(SB)
	CMP	g, R3
	BNED	ok
	JMP	runtime·badmcall(SB)
ok:
	MOVD	fn+0(FP), CTXT			// context
	MOVD	0(CTXT), R4			// code pointer
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, RSP	// sp = m->g0->sched.sp
	MOVD	R3, -8(RSP)
	MOVD	$0, -16(RSP)
	SUB	$16, RSP
	CALL	(R4)
	JMP	runtime·badmcall2(SB)

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
	UNDEF
	CALL	(LR)	// make sure this function is not leaf
	RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
	// TODO(aram):
	MOVD	$8, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
// Caller has already loaded:
// R3 prolog's LR (R30)
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$9, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	// TODO(aram):
	MOVD	$10, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·stackBarrier(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$11, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// reflectcall: call a function with the given argument list
// func call(argtype *_type, f *FuncVal, arg *byte, argsize, retoffset uint32).
// we don't have variable-sized frames, so we use a small number
// of constant-sized-frame functions to encode a few bits of size in the pc.
// Caution: ugly multiline assembly macros in your future!

TEXT reflect·call(SB), NOSPLIT, $0-0
	JMP	·reflectcall(SB)

TEXT ·reflectcall(SB), NOSPLIT, $-8-32
	// TODO(aram):
	MOVD	$12, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// AES hashing not implemented for SPARC64.
TEXT runtime·aeshash(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$13, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET
TEXT runtime·aeshash32(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$14, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET
TEXT runtime·aeshash64(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$15, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET
TEXT runtime·aeshashstr(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$16, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	// TODO(aram):
	MOVD	$17, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT, $-8-16
	// TODO(aram):
	MOVD	$18, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// Save state of caller into g->sched. Smashes R0.
TEXT gosave<>(SB),NOSPLIT,$-8
	// TODO(aram):
	MOVD	$19, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// func asmcgocall(fn, arg unsafe.Pointer) int32
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
TEXT ·asmcgocall(SB),NOSPLIT,$0-20
	// TODO(aram):
	MOVD	$20, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),NOSPLIT,$24-24
	// TODO(aram):
	MOVD	$21, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.go for more details.
TEXT ·cgocallback_gofunc(SB),NOSPLIT,$24-24
	// TODO(aram):
	MOVD	$22, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// Called from cgo wrappers, this function returns g->m->curg.stack.hi.
// Must obey the gcc calling convention.
TEXT _cgo_topofstack(SB),NOSPLIT,$24
	// TODO(aram):
	MOVD	$23, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// void setg(G*); set g. for use by needm.
TEXT runtime·setg(SB), NOSPLIT, $0-8
	// TODO(aram):
	MOVD	$24, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$8
	// TODO(aram):
	MOVD	$25, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·getcallerpc(SB),NOSPLIT,$8-16
	// TODO(aram):
	MOVD	$26, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·setcallerpc(SB),NOSPLIT,$8-16
	// TODO(aram):
	MOVD	$27, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·getcallersp(SB),NOSPLIT,$0-16
	// TODO(aram):
	MOVD	$28, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·abort(SB),NOSPLIT,$-8-0
	JMPL	ZR, ZR
	UNDEF

// func cputicks() int64
TEXT runtime·cputicks(SB),NOSPLIT,$0-0
	RD	TICK, R1
	MOVD	R1, ret+0(FP)
	RET

// memhash_varlen(p unsafe.Pointer, h seed) uintptr
// redirects to memhash(p, h, size) using the size
// stored in the closure.
TEXT runtime·memhash_varlen(SB),NOSPLIT,$40-24
	// TODO(aram):
	MOVD	$30, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// memequal(p, q unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT,$-8-25
	// TODO(aram):
	MOVD	$31, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$40-17
	// TODO(aram):
	MOVD	$32, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·cmpstring(SB),NOSPLIT|NOFRAME,$0-40
	// TODO(aram):
	MOVD	$33, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT bytes·Compare(SB),NOSPLIT|NOFRAME,$0-56
	// TODO(aram):
	MOVD	$34, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// On entry:
// R0 is the length of s1
// R1 is the length of s2
// R2 points to the start of s1
// R3 points to the start of s2
// R7 points to return value (-1/0/1 will be written here)
//
// On exit:
// R4, R5, and R6 are clobbered
TEXT runtime·cmpbody<>(SB),NOSPLIT|NOFRAME,$0-0
	// TODO(aram):
	MOVD	$35, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// eqstring tests whether two strings are equal.
// The compiler guarantees that strings passed
// to eqstring have equal length.
// See runtime_test.go:eqstring_generic for
// equivalent Go code.
TEXT runtime·eqstring(SB),NOSPLIT,$0-33
	// TODO(aram):
	MOVD	$36, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

//
// functions for other packages
//
TEXT bytes·IndexByte(SB),NOSPLIT,$0-40
	// TODO(aram):
	MOVD	$37, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT strings·IndexByte(SB),NOSPLIT,$0-32
	// TODO(aram):
	MOVD	$38, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// TODO: share code with memequal?
TEXT bytes·Equal(SB),NOSPLIT,$0-49
	// TODO(aram):
	MOVD	$39, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·fastrand1(SB),NOSPLIT,$-8-4
	// TODO(aram):
	MOVD	$40, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT runtime·return0(SB), NOSPLIT, $0
	// TODO(aram):
	MOVD	$41, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$42, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

// TODO(aram):
TEXT runtime·prefetcht0(SB),NOSPLIT,$0-8
	RET

TEXT runtime·prefetcht1(SB),NOSPLIT,$0-8
	RET

TEXT runtime·prefetcht2(SB),NOSPLIT,$0-8
	RET

TEXT runtime·prefetchnta(SB),NOSPLIT,$0-8
	RET

TEXT runtime·sigreturn(SB),NOSPLIT,$0-8
        RET

// This is called from .init_array and follows the platform, not Go, ABI.
TEXT runtime·addmoduledata(SB),NOSPLIT,$0-0
	// TODO(aram):
	MOVD	$43, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT ·checkASM(SB),NOSPLIT,$0-1
	// TODO(aram):
	MOVD	$44, R1
	ADD	$'!', R1, R1
	MOVB	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	UNDEF
	RET

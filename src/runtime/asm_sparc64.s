// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"
#include "asm_sparc64.h"

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
	MOVD	fn+0(FP), R3	// R3 = fn
	MOVD	R3, CTXT		// context
	MOVD	g_m(g), R4	// R4 = m

	MOVD	m_gsignal(R4), R5	// R5 = gsignal
	CMP	g, R5
	BED	noswitch

	MOVD	m_g0(R4), R5	// R5 = g0
	CMP	g, R5
	BED	noswitch

	MOVD	m_curg(R4), R6
	CMP	g, R6
	BED	switch

	// Bad: g is not gsignal, not g0, not curg. What is it?
	// Hide call from linker nosplit analysis.
	MOVD	$runtime·badsystemstack(SB), R3
	CALL	(R3)

switch:
	// save our state in g->sched. Pretend to
	// be systemstack_switch if the G stack is scanned.
	MOVD	$runtime·systemstack_switch(SB), R6
	ADD	$8, R6	// get past prologue
	MOVD	R6, (g_sched+gobuf_pc)(g)
	MOVD	RSP, (g_sched+gobuf_sp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	g, (g_sched+gobuf_g)(g)

	// switch to g0
	MOVD	R5, g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), R3
	// make it look like mstart called systemstack on g0, to stop traceback
	SUB	$16, R3
	AND	$~15, R3
	MOVD	$runtime·mstart(SB), R4
	MOVD	R4, 0(R3)
	MOVD	R3, RSP

	// call target function
	MOVD	0(CTXT), R3	// code pointer
	CALL	(R3)

	// switch back to g
	MOVD	g_m(g), R3
	MOVD	m_curg(R3), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), RSP
	MOVD	$0, (g_sched+gobuf_sp)(g)
	RET

noswitch:
	// already on m stack, just call directly
	MOVD	0(CTXT), R3	// code pointer
	CALL	(R3)
	RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
// Caller has already loaded:
// R3 prolog's LR
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB),NOSPLIT,$-8-0
	// Cannot grow scheduler stack (m->g0).
	MOVD	g_m(g), R8
	MOVD	m_g0(R8), R4
	CMP	g, R4
	BNED	2(PC)
	JMP	runtime·abort(SB)

	// Cannot grow signal stack (m->gsignal).
	MOVD	m_gsignal(R8), R4
	CMP	g, R4
	BNED	2(PC)
	JMP	runtime·abort(SB)

	// Called from f.
	// Set g->sched to context in f
	MOVD	CTXT, (g_sched+gobuf_ctxt)(g)
	MOVD	RSP, (g_sched+gobuf_sp)(g)
	MOVD	LR, (g_sched+gobuf_pc)(g)
	MOVD	R3, (g_sched+gobuf_lr)(g)

	// Called from f.
	// Set m->morebuf to f's callers.
	MOVD	R3, (m_morebuf+gobuf_pc)(R8)	// f's caller's PC
	MOVD	RSP, (m_morebuf+gobuf_sp)(R8)	// f's caller's RSP
	MOVD	g, (m_morebuf+gobuf_g)(R8)

	// Call newstack on m->g0's stack.
	MOVD	m_g0(R8), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), RSP
	CALL	runtime·newstack(SB)

	// Not reached, but make sure the return PC from the call to newstack
	// is still in this function, and not the beginning of the next.
	UNDEF

TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	ZR, CTXT
	JMP	runtime·morestack(SB)

TEXT runtime·stackBarrier(SB),NOSPLIT,$0
	// We came here via a RET to an overwritten LR.
	// R8 may be live (see return0). Other registers are available.

	// Get the original return PC, g.stkbar[g.stkbarPos].savedLRVal.
	MOVD	(g_stkbar+slice_array)(g), R4
	MOVD	g_stkbarPos(g), R5
	MOVD	$stkbar__size, R6
	MULD	R5, R6
	ADD	R4, R6
	MOVD	stkbar_savedLRVal(R6), R6
	// Record that this stack barrier was hit.
	ADD	$1, R5
	MOVD	R5, g_stkbarPos(g)
	// Jump to the original return PC.
	JMPL	R6, ZR

// reflectcall: call a function with the given argument list
// func call(argtype *_type, f *FuncVal, arg *byte, argsize, retoffset uint32).
// we don't have variable-sized frames, so we use a small number
// of constant-sized-frame functions to encode a few bits of size in the pc.
// Caution: ugly multiline assembly macros in your future!

#define DISPATCH(NAME,MAXSIZE)		\
	MOVD	$MAXSIZE, TMP;		\
	CMP	TMP, RT1;		\
	BGD	3(PC);			\
	MOVD	$NAME(SB), RT1;	\
	JMPL	RT1, ZR
// Note: can't just "B NAME(SB)" - bad inlining results.

TEXT reflect·call(SB), NOSPLIT, $0-0
	JMP	·reflectcall(SB)

TEXT ·reflectcall(SB), NOSPLIT, $-8-32
	MOVUW argsize+24(FP), RT1
	// NOTE(rsc): No call16, because CALLFN needs four words
	// of argument space to invoke callwritebarrier.
	DISPATCH(runtime·call32, 32)
	DISPATCH(runtime·call64, 64)
	DISPATCH(runtime·call128, 128)
	DISPATCH(runtime·call256, 256)
	DISPATCH(runtime·call512, 512)
	DISPATCH(runtime·call1024, 1024)
	DISPATCH(runtime·call2048, 2048)
	DISPATCH(runtime·call4096, 4096)
	DISPATCH(runtime·call8192, 8192)
	DISPATCH(runtime·call16384, 16384)
	DISPATCH(runtime·call32768, 32768)
	DISPATCH(runtime·call65536, 65536)
	DISPATCH(runtime·call131072, 131072)
	DISPATCH(runtime·call262144, 262144)
	DISPATCH(runtime·call524288, 524288)
	DISPATCH(runtime·call1048576, 1048576)
	DISPATCH(runtime·call2097152, 2097152)
	DISPATCH(runtime·call4194304, 4194304)
	DISPATCH(runtime·call8388608, 8388608)
	DISPATCH(runtime·call16777216, 16777216)
	DISPATCH(runtime·call33554432, 33554432)
	DISPATCH(runtime·call67108864, 67108864)
	DISPATCH(runtime·call134217728, 134217728)
	DISPATCH(runtime·call268435456, 268435456)
	DISPATCH(runtime·call536870912, 536870912)
	DISPATCH(runtime·call1073741824, 1073741824)
	MOVD	$runtime·badreflectcall(SB), R1
	JMPL	R1, ZR

#define CALLFN(NAME,MAXSIZE)			\
TEXT NAME(SB), WRAPPER, $MAXSIZE-24;		\
	NO_LOCAL_POINTERS;			\
	/* copy arguments to stack */		\
	MOVD	arg+16(FP), R3;			\
	MOVUW	argsize+24(FP), R4;			\
	MOVD	RSP, R5;				\
	ADD	$(8-1), R5;			\
	SUB	$1, R3;				\
	ADD	R5, R4;				\
	CMP	R5, R4;				\
	BED	6(PC);				\
	MOVUB	1(R3), R6;			\
	ADD	$1, R3;				\
	MOVUB	R6, 1(R5);			\
	ADD	$1, R5;				\
	JMP	-6(PC);				\
	/* call function */			\
	MOVD	f+8(FP), CTXT;			\
	MOVD	(CTXT), R1;			\
	PCDATA  $PCDATA_StackMapIndex, $0;	\
	CALL	(R1);				\
	/* copy return values back */		\
	MOVD	arg+16(FP), R3;			\
	MOVUW	n+24(FP), R4;			\
	MOVUW	retoffset+28(FP), R6;		\
	MOVD	RSP, R5;				\
	ADD	R6, R5; 			\
	ADD	R6, R3;				\
	SUB	R6, R4;				\
	ADD	$(8-1), R5;			\
	SUB	$1, R3;				\
	ADD	R5, R4;				\
loop:						\
	CMP	R5, R4;				\
	BED	end;				\
	MOVUB	1(R5), R6;			\
	ADD	$1, R5;				\
	MOVUB	R6, 1(R3);			\
	ADD	$1, R3;			\
	JMP	loop;				\
end:						\
	/* execute write barrier updates */	\
	MOVD	argtype+0(FP), R8;		\
	MOVD	arg+16(FP), R3;			\
	MOVUW	n+24(FP), R4;			\
	MOVUW	retoffset+28(FP), R6;		\
	MOVD	R8, 8(RSP);			\
	MOVD	R3, 16(RSP);			\
	MOVD	R4, 24(RSP);			\
	MOVD	R6, 32(RSP);			\
	CALL	runtime·callwritebarrier(SB);	\
	RET

// These have 8 added to make the overall frame size a multiple of 16,
// as required by the ABI. (There is another +8 for the saved LR.)
CALLFN(·call32, 40 )
CALLFN(·call64, 72 )
CALLFN(·call128, 136 )
CALLFN(·call256, 264 )
CALLFN(·call512, 520 )
CALLFN(·call1024, 1032 )
CALLFN(·call2048, 2056 )
CALLFN(·call4096, 4104 )
CALLFN(·call8192, 8200 )
CALLFN(·call16384, 16392 )
CALLFN(·call32768, 32776 )
CALLFN(·call65536, 65544 )
CALLFN(·call131072, 131080 )
CALLFN(·call262144, 262152 )
CALLFN(·call524288, 524296 )
CALLFN(·call1048576, 1048584 )
CALLFN(·call2097152, 2097160 )
CALLFN(·call4194304, 4194312 )
CALLFN(·call8388608, 8388616 )
CALLFN(·call16777216, 16777224 )
CALLFN(·call33554432, 33554440 )
CALLFN(·call67108864, 67108872 )
CALLFN(·call134217728, 134217736 )
CALLFN(·call268435456, 268435464 )
CALLFN(·call536870912, 536870920 )
CALLFN(·call1073741824, 1073741832 )

// AES hashing not implemented for SPARC64.
TEXT runtime·aeshash(SB),NOSPLIT,$-8-0
	MOVW	(ZR), R1
TEXT runtime·aeshash32(SB),NOSPLIT,$-8-0
	MOVW	(ZR), R1
TEXT runtime·aeshash64(SB),NOSPLIT,$-8-0
	MOVW	(ZR), R1
TEXT runtime·aeshashstr(SB),NOSPLIT,$-8-0
	MOVW	(ZR), R1
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	RD	CCR, R2
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT, $-8-16
	MOVD	0(RSP), R1
	SUB	$4, R1
	MOVD	R1, LR

	MOVD	fv+0(FP), CTXT
	MOVD	argp+8(FP), RSP
	SUB	$8, RSP
	MOVD	0(CTXT), R3
	JMPL	R3, ZR

// Save state of caller into g->sched.
TEXT gosave<>(SB),NOSPLIT,$-8
	MOVD	LR, (g_sched+gobuf_pc)(g)
	MOVD	RSP, (g_sched+gobuf_sp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	$0, (g_sched+gobuf_ret)(g)
	MOVD	$0, (g_sched+gobuf_ctxt)(g)
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
	MOVD	$fn+0(FP), R1
	MOVD	R1, (8+STACK_BIAS)(RSP)
	MOVD	frame+8(FP), R1
	MOVD	R1, (16+STACK_BIAS)(RSP)
	MOVD	framesize+16(FP), R1
	MOVD	R1, (24+STACK_BIAS)(RSP)
	MOVD	$runtime·cgocallback_gofunc(SB), R1
	CALL	R1
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
	// g and TMP might be clobbered by load_g. They
	// are callee-save in the gcc calling convention, so save them.
	MOVD	TMP, savedTMP-8(SP)
	MOVD	g, saveG-16(SP)

	CALL	runtime·load_g(SB)
	MOVD	g_m(g), R1
	MOVD	m_curg(R1), R1
	MOVD	(g_stack+stack_hi)(R1), R1

	MOVD	saveG-16(SP), g
	MOVD	savedTMP-8(SP), TMP
	RET

// void setg(G*); set g. for use by needm.
TEXT runtime·setg(SB), NOSPLIT, $0-8
	MOVD	gg+0(FP), g
	// This only happens if iscgo, so jump straight to save_g
	CALL	runtime·save_g(SB)
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$8
	MOVD	R8, g
	MOVD	TMP, savedTMP-8(SP)
	CALL	runtime·save_g(SB)
	MOVD	savedTMP-8(SP), TMP
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
	MOVD	argp+0(FP), R1
	SUB	$8, R1
	MOVD	R1, ret+8(FP)
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
	MOVW	ZR, R8
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT,$-8-0
	MOVD	R1, R1	// NOP
	CALL	runtime·goexit1(SB)	// does not return

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

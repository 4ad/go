// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"
#include "asm_sparc64.h"

DATA dbgbuf(SB)/8, $"\n\n"
GLOBL dbgbuf(SB), $8

// Note: define used in this file to avoid affecting registers.
// #MemIssue|#Sync|#LoadLoad|#StoreLoad|#LoadStore|#StoreStore
#define REGFLUSH	\
	MEMBAR	$111;	\
	FLUSHW;		\
	MEMBAR	$111

TEXT runtime·rt0_go(SB),NOSPLIT,$16-0
	// BSP = stack; O0 = argc; O1 = argv

	// initialize essential registers
	CALL	runtime·reginit(SB)

	MOVW	O0, FIXED_FRAME+0(BSP)	// copy argc
	MOVD	O1, FIXED_FRAME+8(BSP)	// copy argv

	// create istack out of the given (operating system) stack.
	// _cgo_init may update stackguard.
	MOVD	$runtime·g0(SB), g
	MOVD	BSP, RT1
	// must be larger than _StackSystem
	MOVD	$(-64*1024)(BSP), RT2
	MOVD	RT2, g_stackguard0(g)
	MOVD	RT2, g_stackguard1(g)
	MOVD	RT2, (g_stack+stack_lo)(g)
	MOVD	RT1, (g_stack+stack_hi)(g)

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
	// update stackguard after _cgo_init
	MOVD	(g_stack+stack_lo)(g), I1
	ADD	$const__StackGuard, I1
	MOVD	I1, g_stackguard0(g)
	MOVD	I1, g_stackguard1(g)

	// set the per-goroutine and per-mach "registers"
	MOVD	$runtime·m0(SB), I1

	// save m->g0 = g0
	MOVD	g, m_g0(I1)
	// save m0 to g0->m
	MOVD	I1, g_m(g)

	CALL	runtime·check(SB)

	// argc, argv already copied.
	CALL	runtime·args(SB)
	CALL	runtime·osinit(SB)
	CALL	runtime·schedinit(SB)

	// create a new goroutine to start program
	MOVD	ZR, FIXED_FRAME+0(BSP)
	MOVD	$runtime·mainPC(SB), RT1		// entry
	MOVD	RT1, FIXED_FRAME+8(BSP)
	CALL	runtime·newproc(SB)

	// start this M
	CALL	runtime·mstart(SB)

	MOVD	ZR, (ZR)	// boom
	UNDEF

DATA	runtime·mainPC+0(SB)/8,$runtime·main(SB)
GLOBL	runtime·mainPC(SB),RODATA,$8

TEXT runtime·breakpoint(SB),NOSPLIT|NOFRAME,$0-0
	TA	$0x81
	RET

TEXT runtime·asminit(SB),NOSPLIT|NOFRAME,$0-0
	RET

TEXT runtime·reginit(SB),NOSPLIT|NOFRAME,$0-0
	// initialize essential FP registers
	FMOVD	$2.0, D28
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), I1
	MOVD	BSP, I3
	MOVD	I3, gobuf_sp(I1)
	MOVD	OLR, gobuf_pc(I1)
	MOVD	g, gobuf_g(I1)
	MOVD	ZR, gobuf_lr(I1)
	MOVD	ZR, gobuf_ret(I1)
	MOVD	ZR, gobuf_ctxt(I1)
	MOVD	BFP, I3
	MOVD	I3, gobuf_bp(I1)
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), L6
	MOVD	gobuf_g(L6), g
	CALL	runtime·save_g(SB)

	MOVD	buf+0(FP), L6
	MOVD	gobuf_g(L6), g
	MOVD	0(g), I4	// make sure g is not nil
	MOVD	gobuf_lr(L6), OLR
	MOVD	gobuf_ret(L6), RT1
	MOVD	gobuf_ctxt(L6), CTXT
	MOVD	gobuf_sp(L6), I3
	MOVD	gobuf_bp(L6), I4
	// restore continuation's ILR before resetting the stack pointer
	// otherwise a spill will overwrite the saved link register.
	MOVD	120(I3), ILR
	MOVD	I3, BSP
	MOVD	I4, BFP

	MOVD	ZR, gobuf_sp(L6)
	MOVD	ZR, gobuf_ret(L6)
	MOVD	ZR, gobuf_lr(L6)
	MOVD	ZR, gobuf_ctxt(L6)
	MOVD	ZR, gobuf_bp(L6)
	CMP	ZR, ZR // set condition codes for == test, needed by stack split
	MOVD	gobuf_pc(L6), O0
	JMPL	$8(O0), ZR

// void mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT|NOFRAME, $0-8
	// Save caller state in g->sched
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	OLR, (g_sched+gobuf_pc)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	g, (g_sched+gobuf_g)(g)
	MOVD	BFP, I3
	MOVD	I3, (g_sched+gobuf_bp)(g)

	// Switch to m->g0 & its stack, call fn.
	MOVD	g, I1
	MOVD	g_m(g), O0
	MOVD	m_g0(O0), g
	CALL	runtime·save_g(SB)
	CMP	g, I1
	BNED	ok
	JMP	runtime·badmcall(SB)
ok:

	MOVD	fn+0(FP), CTXT			// context
	MOVD	0(CTXT), I4			// code pointer
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP	// sp = m->g0->sched.sp
	MOVD	TMP, BFP
	SUB	$FIXED_FRAME+16, BSP
	MOVD	I1, (FIXED_FRAME+0)(BSP)
	MOVD	$0, (FIXED_FRAME+8)(BSP)
	CALL	(I4)
	JMP	runtime·badmcall2(SB)

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
	UNDEF
	UNDEF
	UNDEF
	UNDEF
	UNDEF
	CALL	(ILR)	// make sure this function is not leaf
	RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
	MOVD	fn+0(FP), I1	// I1 = fn
	MOVD	I1, CTXT	// context
	MOVD	g_m(g), I4	// I4 = m

	MOVD	m_gsignal(I4), L6	// L6 = gsignal
	CMP	g, L6
	BED	noswitch

	MOVD	m_g0(I4), L6	// L6 = g0
	CMP	g, L6
	BED	noswitch

	MOVD	m_curg(I4), O0
	CMP	g, O0
	BED	switch

	// Bad: g is not gsignal, not g0, not curg. What is it?
	// Hide call from linker nosplit analysis.
	MOVD	$runtime·badsystemstack(SB), I1
	CALL	(I1)

switch:
	// save our state in g->sched. Pretend to
	// be systemstack_switch if the G stack is scanned.
	MOVD	$runtime·systemstack_switch(SB), O0
	ADD	$20, O0	// get past prologue
	MOVD	O0, (g_sched+gobuf_pc)(g)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	BFP, TMP
	MOVD	TMP, (g_sched+gobuf_bp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	$0, (g_sched+gobuf_ret)(g)
	MOVD	$0, (g_sched+gobuf_ctxt)(g)
	MOVD	g, (g_sched+gobuf_g)(g)

	// switch to g0

	MOVD	L6, g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), I1
	MOVD	I1, BFP	// subtle
	// make it look like mstart called systemstack on g0, to stop traceback.
	SUB	$FIXED_FRAME, I1
	MOVD	$runtime·mstart(SB), I4
	MOVD	I4, 120(I1)
	MOVD	I4, ILR
	MOVD	I1, BSP

	// call target function
	MOVD	0(CTXT), I1	// code pointer
	CALL	(I1)

	// switch back to g
	MOVD	g_m(g), I1
	MOVD	m_curg(I1), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_bp)(g), TMP
	MOVD	TMP, BFP
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP
	MOVD	ZR, (g_sched+gobuf_sp)(g)
	MOVD	ZR, (g_sched+gobuf_bp)(g)
	RET

noswitch:
	// already on m stack, just call directly
	MOVD	0(CTXT), I1	// code pointer
	CALL	(I1)
	RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
// Caller has already loaded:
// I1 prolog's LR
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB),NOSPLIT|NOFRAME,$0-0
	// Cannot grow scheduler stack (m->g0).
	MOVD	g_m(g), O0
	MOVD	m_g0(O0), I4
	CMP	g, I4
	BNED	2(PC)
	JMP	runtime·abort(SB)

	// Cannot grow signal stack (m->gsignal).
	MOVD	m_gsignal(O0), I4
	CMP	g, I4
	BNED	2(PC)
	JMP	runtime·abort(SB)

	// Called from f.
	// Set g->sched to context in f
	MOVD	CTXT, (g_sched+gobuf_ctxt)(g)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	BFP, TMP
	MOVD	TMP, (g_sched+gobuf_bp)(g)
	MOVD	OLR, (g_sched+gobuf_pc)(g)
	MOVD	I1, (g_sched+gobuf_lr)(g)

	// Called from f.
	// Set m->morebuf to f's callers.
	MOVD	I1, (m_morebuf+gobuf_pc)(O0)	// f's caller's PC
	MOVD	BSP, TMP
	MOVD	TMP, (m_morebuf+gobuf_sp)(O0)	// f's caller's BSP
	MOVD	BFP, TMP
	MOVD	TMP, (m_morebuf+gobuf_bp)(O0)	// f's caller's BFP
	MOVD	g, (m_morebuf+gobuf_g)(O0)

	// Call newstack on m->g0's stack.
	MOVD	m_g0(O0), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP

	CALL	runtime·newstack(SB)

	// Not reached, but make sure the return PC from the call to newstack
	// is still in this function, and not the beginning of the next.
	UNDEF

TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	ZR, CTXT
	JMP	runtime·morestack(SB)

TEXT runtime·stackBarrier(SB),NOSPLIT|NOFRAME,$0
	// We came here via a RET to an overwritten LR.
	// RT1 may be live (see return0). Other registers are available.

	// Get the original return PC, g.stkbar[g.stkbarPos].savedLRVal.
	MOVD	(g_stkbar+slice_array)(g), I4
	MOVD	g_stkbarPos(g), L6
	MOVD	$stkbar__size, O1
	MULD	L6, O1
	ADD	I4, O1
	MOVD	stkbar_savedLRVal(O1), O1
	// Record that this stack barrier was hit.
	ADD	$1, L6
	MOVD	L6, g_stkbarPos(g)
	// Jump to the original return PC.
	JMPL	O1, ZR

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

TEXT reflect·call(SB), NOSPLIT|NOFRAME, $0-0
	JMP	·reflectcall(SB)

TEXT ·reflectcall(SB), NOSPLIT|NOFRAME, $0-32
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
	MOVD	$runtime·badreflectcall(SB), I3
	JMPL	I3, ZR

#define CALLFN(NAME,MAXSIZE)			\
TEXT NAME(SB), WRAPPER, $MAXSIZE-24;		\
	NO_LOCAL_POINTERS;			\
	/* copy arguments to stack */		\
	MOVD	arg+16(FP), I1;			\
	MOVUW	argsize+24(FP), I4;		\
	MOVD	BSP, L6;			\
	ADD	$(FIXED_FRAME-1), L6;		\
	SUB	$1, I1;				\
	ADD	L6, I4;				\
	CMP	L6, I4;				\
	BED	6(PC);				\
	ADD	$1, I1;				\
	MOVUB	(I1), O1;			\
	ADD	$1, L6;				\
	MOVUB	O1, (L6);			\
	JMP	-6(PC);				\
	/* call function */			\
	MOVD	f+8(FP), CTXT;			\
	MOVD	(CTXT), I3;			\
	PCDATA	$PCDATA_StackMapIndex, $0;	\
	CALL	(I3);				\
	/* copy return values back */		\
	MOVD	arg+16(FP), I1;			\
	MOVUW	n+24(FP), I4;			\
	MOVUW	retoffset+28(FP), O1;		\
	MOVD	BSP, L6;			\
	ADD	O1, L6; 			\
	ADD	O1, I1;				\
	SUB	O1, I4;				\
	ADD	$(FIXED_FRAME-1), L6;		\
	SUB	$1, I1;				\
	ADD	L6, I4;				\
loop:						\
	CMP	L6, I4;				\
	BED	end;				\
	ADD	$1, L6;				\
	MOVUB	(L6), O1;			\
	ADD	$1, I1;				\
	MOVUB	O1, (I1);			\
	JMP	loop;				\
end:						\
	/* execute write barrier updates */	\
	MOVD	argtype+0(FP), O0;		\
	MOVD	arg+16(FP), I1;			\
	MOVUW	n+24(FP), I4;			\
	MOVUW	retoffset+28(FP), O1;		\
	MOVD	O0, (FIXED_FRAME+0)(BSP);	\
	MOVD	I1, (FIXED_FRAME+8)(BSP);	\
	MOVD	I4, (FIXED_FRAME+16)(BSP);	\
	MOVD	O1, (FIXED_FRAME+24)(BSP);	\
	CALL	runtime·callwritebarrier(SB);	\
	RET

CALLFN(·call32, 32)
CALLFN(·call64, 64)
CALLFN(·call128, 128)
CALLFN(·call256, 256)
CALLFN(·call512, 512)
CALLFN(·call1024, 1024)
CALLFN(·call2048, 2048)
CALLFN(·call4096, 4096)
CALLFN(·call8192, 8192)
CALLFN(·call16384, 16384)
CALLFN(·call32768, 32768)
CALLFN(·call65536, 65536)
CALLFN(·call131072, 131072)
CALLFN(·call262144, 262144)
CALLFN(·call524288, 524288)
CALLFN(·call1048576, 1048576)
CALLFN(·call2097152, 2097152)
CALLFN(·call4194304, 4194304)
CALLFN(·call8388608, 8388608)
CALLFN(·call16777216, 16777216)
CALLFN(·call33554432, 33554432)
CALLFN(·call67108864, 67108864)
CALLFN(·call134217728, 134217728)
CALLFN(·call268435456, 268435456)
CALLFN(·call536870912, 536870912)
CALLFN(·call1073741824, 1073741824)


// AES hashing not implemented for SPARC64.
TEXT runtime·aeshash(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), I3
TEXT runtime·aeshash32(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), I3
TEXT runtime·aeshash64(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), I3
TEXT runtime·aeshashstr(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), I3
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	RD	CCR, I5
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT|NOFRAME, $0-16
	// We're in the same stack frame and same register window as
	// our caller, deferreturn, so this retrieves the return address
	// to deferreturn's caller.
	// 
	// We need to subtract -8 from this value, because the deferred
	// functions returns to $8(ILR).
	MOVD	(8*15)(BSP), I3
	SUB	$8, I3
	// ILR will become OLR once we RESTORE, so the deferred function
	// will return to the CALL instruction, calling deferreturn
	// again.
	MOVD	I3, ILR

	// fv is the deferred function. I1 will become O1 once we
	// RESTORE. This affects registers in deferreturn's caller,
	// but that's ok, registers are caller-save in Go.
	MOVD	fv+0(FP), CTXT
	MOVD	0(CTXT), I1

	// We must RESTORE here, because the deferred function expects
	// to be called by deferreturn's caller, so deferreturn's
	// caller and the deferred functions must be in adjacent
	// register windows.
	// 
	// After the RESTORE, BSP and BFP will select deferreturn's
	// caller activation record, but that is exactly what we want,
	// logically deferreturn's caller calls the deferred function;
	// deferreturn will have put the arguments to the deferred
	// function in the correct place in the caller frame.
	RESTORE	$0, ZR, ZR
	JMPL	O1, ZR

// Save state of caller into g->sched.
TEXT gosave<>(SB),NOSPLIT|NOFRAME,$0
	MOVD	OLR, (g_sched+gobuf_pc)(g)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	$0, (g_sched+gobuf_ret)(g)
	MOVD	$0, (g_sched+gobuf_ctxt)(g)
	MOVD	BFP, TMP
	MOVD	TMP, (g_sched+gobuf_bp)(g)
	RET

// func asmcgocall(fn, arg unsafe.Pointer) int32
TEXT ·asmcgocall(SB),NOSPLIT,$16-20
	MOVD	fn+0(FP), O1
	MOVD	arg+8(FP), O0

	// save original stack pointer
	MOVD	BSP, I1
	// save g
	MOVD	g, I2

	MOVD	g_m(g), L1
	MOVD	m_g0(L1), L2
	CMP	g, L2
	BED	g0

	CALL	gosave<>(SB)
	MOVD	L2, g
	CALL	runtime·save_g(SB)

	MOVD	(g_sched+gobuf_sp)(g), L4
	MOVD	L4, BFP
	SUB	$(16+FIXED_FRAME), L4, L5
	MOVD	L5, BSP
	SUB	$STACK_BIAS, L4
	MOVD	L4, 112(BSP)
	MOVD	ILR, 120(BSP)

g0:
	// Now on a scheduling stack (a pthread-created stack).
	// save old g on stack
	MOVD	I2, (16+FIXED_FRAME-8)(BSP)
	// save depth in old g stack, can't just save SP, as stack
	// might be copied during a callback
	MOVD	(g_stack+stack_hi)(I2), L1
	SUB	I1, L1
	MOVD	L1, (16+FIXED_FRAME-16)(BSP)

	// call target function
//	CALL	runtime·save_g(SB)
	CALL	(O1)

	// Restore g
	MOVD	(16+FIXED_FRAME-8)(BSP), g
	CALL	runtime·save_g(SB)
	// Retrieve stack pointer
	MOVD	(g_stack+stack_hi)(g), L1
	MOVD	(16+FIXED_FRAME-16)(BSP), L2
	SUB	L2, L1
	// Restore frame pointer
	MOVD	112(L1), L3
	ADD	$STACK_BIAS, L3
	MOVD	L3, BFP
	// Restore stack pointer
	MOVD	L1, BSP

	MOVW	O0, ret+16(FP)
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),NOSPLIT,$32-24
	MOVD	$fn+0(FP), I3
	MOVD	I3, (FIXED_FRAME+0)(BSP)
	MOVD	frame+8(FP), I3
	MOVD	I3, (FIXED_FRAME+8)(BSP)
	MOVD	framesize+16(FP), I3
	MOVD	I3, (FIXED_FRAME+16)(BSP)
	MOVD	$runtime·cgocallback_gofunc(SB), I3
	CALL	I3
	RET

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.go for more details.
TEXT ·cgocallback_gofunc(SB),NOSPLIT,$32-24
	NO_LOCAL_POINTERS

	// Load m and g from thread-local storage.
	MOVB	runtime·iscgo(SB), I1
	CMP	I1, ZR
	BED	nocgo
	CALL	runtime·load_g(SB)
nocgo:

	// If g is nil, Go did not create the current thread.
	// Call needm to obtain one for temporary use.
	// In this case, we're running on the thread stack, so there's
	// lots of space, but the linker doesn't know. Hide the call from
	// the linker analysis by using an indirect call.
	CMP	g, ZR
	BED	needm

	MOVD	g_m(g), O0
	MOVD	O0, savedm-8(SP)
	JMP	havem

needm:
	MOVD	g, savedm-8(SP) // g is zero, so is m.
	MOVD	$runtime·needm(SB), RT1
	CALL	(RT1)

	// Set m->sched.sp = SP, so that if a panic happens
	// during the function we are about to execute, it will
	// have a valid SP to run on the g0 stack.
	// The next few lines (after the havem label)
	// will save this SP onto the stack and then write
	// the same SP back to m->sched.sp. That seems redundant,
	// but if an unrecovered panic happens, unwindm will
	// restore the g->sched.sp from the stack location
	// and then systemstack will try to use it. If we don't set it here,
	// that restored SP will be uninitialized (typically 0) and
	// will not be usable.
	MOVD	g_m(g), O0
	MOVD	m_g0(O0), I1
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(I1)

havem:
	// Now there's a valid m, and we're running on its m->g0.
	// Save current m->g0->sched.sp on stack and then set it to SP.
	// Save current sp in m->g0->sched.sp in preparation for
	// switch back to m->curg stack.
	// NOTE: unwindm knows that the saved g->sched.sp is at 8(I3) aka savedsp-16(SP).
	MOVD	m_g0(O0), I1
	MOVD	(g_sched+gobuf_sp)(I1), I4
	MOVD	I4, savedsp-16(SP)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(I1)

	// Switch to m->curg stack and call runtime.cgocallbackg.
	// Because we are taking over the execution of m->curg
	// but *not* resuming what had been running, we need to
	// save that information (m->curg->sched) so we can restore it.
	// We can restore m->curg->sched.sp easily, because calling
	// runtime.cgocallbackg leaves SP unchanged upon return.
	// To save m->curg->sched.pc, we push it onto the stack.
	// This has the added benefit that it looks to the traceback
	// routine like cgocallbackg is going to return to that
	// PC (because the frame we allocate below has the same
	// size as cgocallback_gofunc's frame declared above)
	// so that the traceback will seamlessly trace back into
	// the earlier calls.
	//
	// In the new goroutine, -16(SP) and -8(SP) are unused.
	MOVD	m_curg(O0), g
	CALL	runtime·save_g(SB)

	MOVD	(g_sched+gobuf_sp)(g), I4 // prepare stack as I4
	MOVD	(g_sched+gobuf_pc)(g), L6
	MOVD	L6, -(FIXED_FRAME+16)(I4)
	MOVD	$-(FIXED_FRAME+16)(I4), TMP
	MOVD	TMP, BSP
	CALL	runtime·cgocallbackg(SB)

	// Restore g->sched (== m->curg->sched) from saved values.
	MOVD	0(BSP), L6
	MOVD	L6, (g_sched+gobuf_pc)(g)
	MOVD	$(FIXED_FRAME+16)(BSP), I4
	MOVD	I4, (g_sched+gobuf_sp)(g)

	// Switch back to m->g0's stack and restore m->g0->sched.sp.
	// (Unlike m->curg, the g0 goroutine never uses sched.pc,
	// so we do not have to restore it.)
	MOVD	g_m(g), O0
	MOVD	m_g0(O0), g
	CALL	runtime·save_g(SB)

	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP
	MOVD	savedsp-16(SP), I4
	MOVD	I4, (g_sched+gobuf_sp)(g)

	// If the m on entry was nil, we called needm above to borrow an m
	// for the duration of the call. Since the call is over, return it with dropm.
	MOVD	savedm-8(SP), O1
	CMP	O1, ZR
	BNED	droppedm
	MOVD	$runtime·dropm(SB), RT1
	CALL	(RT1)
droppedm:

	// Done!
	RET

// Called from cgo wrappers, this function returns g->m->curg.stack.hi.
// Must obey the gcc calling convention.
TEXT _cgo_topofstack(SB),NOSPLIT,$32
	// g and RT1 might be clobbered by load_g. They
	// are callee-save in the gcc calling convention, so save them.
	MOVD	RT1, savedRT1-8(SP)
	MOVD	g, saveG-16(SP)

	CALL	runtime·load_g(SB)
	MOVD	g_m(g), I3
	MOVD	m_curg(I3), I3
	MOVD	(g_stack+stack_hi)(I3), I3

	MOVD	saveG-16(SP), g
	MOVD	savedRT1-8(SP), RT1
	RET

// void setg(G*); set g. for use by needm.
TEXT runtime·setg(SB), NOSPLIT, $0-8
	MOVD	gg+0(FP), g
	// This only happens if iscgo, so jump straight to save_g
	CALL	runtime·save_g(SB)
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$16
	MOVD	O0, g
	MOVD	RT1, savedRT1-8(SP)
	CALL	runtime·save_g(SB)
	MOVD	savedRT1-8(SP), RT1
	RET

// check that SP is in range [g->stack.lo, g->stack.hi)
TEXT runtime·stackcheck(SB), NOSPLIT, $0
	MOVD	BSP, I4
	MOVD	(g_stack+stack_hi)(g), I3
	CMP	I4, I3
	BGD	2(PC);
	UNDEF

	MOVD	(g_stack+stack_lo)(g), I3
	CMP	I3, I4
	BGD	2(PC);
	UNDEF

	RET

TEXT runtime·getcallerpc(SB),NOSPLIT,$16-16
	MOVD	(8*15+FIXED_FRAME+16)(BSP), I1		// LR saved by caller
	MOVD	runtime·stackBarrierPC(SB), I4
	CMP	I1, I4
	BNED	nobar
	// Get original return PC.
	CALL	runtime·nextBarrierPC(SB)
	MOVD	FIXED_FRAME+0(BSP), I1
nobar:
	MOVD	I1, ret+8(FP)
	RET

TEXT runtime·setcallerpc(SB),NOSPLIT,$16-16
	MOVD	pc+8(FP), I1
	MOVD	(8*15+FIXED_FRAME+16)(BSP), I4
	MOVD	runtime·stackBarrierPC(SB), L6
	CMP	I4, L6
	BED	setbar
	MOVD	I1, (8*15+FIXED_FRAME+16)(BSP)		// set LR in caller
	RET
setbar:
	// Set the stack barrier return PC.
	MOVD	I1, FIXED_FRAME+0(BSP)
	CALL	runtime·setNextBarrierPC(SB)
	RET

TEXT runtime·getcallersp(SB),NOSPLIT,$0-16
	MOVD	argp+0(FP), I3
	SUB	$FIXED_FRAME, I3
	MOVD	I3, ret+8(FP)
	RET

TEXT runtime·abort(SB),NOSPLIT|NOFRAME,$0-0
	JMPL	ZR, ZR
	UNDEF

// func cputicks() int64
TEXT runtime·cputicks(SB),NOSPLIT,$0-0
	RD	TICK, I3
	MOVD	I3, ret+0(FP)
	RET

// memhash_varlen(p unsafe.Pointer, h seed) uintptr
// redirects to memhash(p, h, size) using the size
// stored in the closure.
TEXT runtime·memhash_varlen(SB),NOSPLIT,$48-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD	p+0(FP), I1
	MOVD	h+8(FP), I4
	MOVD	8(CTXT), L6
	MOVD	I1, FIXED_FRAME+0(BSP)
	MOVD	I4, FIXED_FRAME+8(BSP)
	MOVD	L6, FIXED_FRAME+16(BSP)
	CALL	runtime·memhash(SB)
	MOVD	FIXED_FRAME+24(BSP), I1
	MOVD	I1, ret+16(FP)
	RET

// memequal(p, q unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT|NOFRAME,$0-25
	MOVD	a+0(FP), I3
	MOVD	b+8(FP), I5
	MOVD	size+16(FP), I1
	ADD	I3, I1, O1
	MOVD	$1, TMP
	MOVB	TMP, ret+24(FP)
	CMP	I3, I5
	BED	done
loop:
	CMP	I3, O1
	BED	done
	MOVUB	(I3), I4
	ADD	$1, I3
	MOVUB	(I5), L6
	ADD	$1, I5
	CMP	I4, L6
	BED	loop

	MOVB	ZR, ret+24(FP)
done:
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$48-17
	MOVD	a+0(FP), I1
	MOVD	b+8(FP), I4
	CMP	I1, I4
	BED	eq
	MOVD	8(CTXT), L6    // compiler stores size at offset 8 in the closure
	MOVD	I1, FIXED_FRAME+0(BSP)
	MOVD	I4, FIXED_FRAME+8(BSP)
	MOVD	L6, FIXED_FRAME+16(BSP)
	CALL	runtime·memequal(SB)
	MOVD	$FIXED_FRAME+24(BSP), I1
	MOVUB	(I1), I1
	MOVB	I1, ret+16(FP)
	RET
eq:
	MOVD	$1, I1
	MOVB	I1, ret+16(FP)
	RET

//
// functions for other packages
//
TEXT bytes·IndexByte(SB),NOSPLIT,$0-40
	MOVD	s+0(FP), I1
	MOVD	s_len+8(FP), I4
	MOVUB	c+24(FP), L6	// byte to find
	MOVD	I1, O1		// store base for later
	SUB	$1, I1
	ADD	I1, I4		// end-1

loop:
	CMP	I1, I4
	BED	notfound
	ADD	$1, I1
	MOVUB	(I1), O0
	CMP	L6, O0
	BNEW	loop

	SUB	O1, I1		// remove base
	MOVD	I1, ret+32(FP)
	RET

notfound:
	MOVD	$-1, I1
	MOVD	I1, ret+32(FP)
	RET

TEXT strings·IndexByte(SB),NOSPLIT,$0-32
	MOVD	p+0(FP), I1
	MOVD	b_len+8(FP), I4
	MOVUB	c+16(FP), L6	// byte to find
	MOVD	I1, O1		// store base for later
	SUB	$1, I1
	ADD	I1, I4		// end-1

loop:
	CMP	I1, I4
	BED	notfound
	ADD	$1, I1
	MOVUB	(I1), O0
	CMP	L6, O0
	BNEW	loop

	SUB	O1, I1		// remove base
	MOVD	I1, ret+24(FP)
	RET

notfound:
	MOVD	$-1, I1
	MOVD	I1, ret+24(FP)
	RET

TEXT runtime·fastrand1(SB),NOSPLIT|NOFRAME,$0-4
	MOVD	g_m(g), I4
	MOVUW	m_fastrand(I4), I1
	ADD	I1, I1
	CMP	ZR, I1
	BGEW	2(PC)
	XOR	$0x88888eef, I1
	MOVW	I1, m_fastrand(I4)
	MOVW	I1, ret+0(FP)
	RET

TEXT runtime·return0(SB), NOSPLIT, $0
	MOVW	ZR, RT1
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	I3, I3	// NOP
	CALL	runtime·goexit1(SB)	// does not return

// TODO(aram):
TEXT runtime·prefetcht0(SB),NOSPLIT|NOFRAME,$0-8
	RET

TEXT runtime·prefetcht1(SB),NOSPLIT|NOFRAME,$0-8
	RET

TEXT runtime·prefetcht2(SB),NOSPLIT|NOFRAME,$0-8
	RET

TEXT runtime·prefetchnta(SB),NOSPLIT|NOFRAME,$0-8
	RET

TEXT runtime·sigreturn(SB),NOSPLIT|NOFRAME,$0-8
	RET

// This is called from .init_array and follows the platform, not Go, ABI.
TEXT runtime·addmoduledata(SB),NOSPLIT,$0-0
	MOVD	runtime·lastmoduledatap(SB), I3
	MOVD	O0, moduledata_next(I3)
	MOVD	O0, runtime·lastmoduledatap(SB)
	RET

TEXT ·checkASM(SB),NOSPLIT,$0-1
	OR	$1, ZR, I1
	MOVB	I1, ret+0(FP)
	RET

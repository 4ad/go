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
	// BSP = stack; R9 = argc; R8 = argv

	// initialize essential registers
	CALL	runtime·reginit(SB)

	SUB	$(FIXED_FRAME+16), BSP
	MOVD	$(FIXED_FRAME+0)(BSP), RT1
	MOVW	R9, (RT1) // argc
	MOVD	R8, FIXED_FRAME+8(BSP) // argv

	// create istack out of the given (operating system) stack.
	// _cgo_init may update stackguard.
	MOVD	$runtime·g0(SB), g
	MOVD BSP, RT1
	MOVD	$(-64*1024)(RT1), RT2
	MOVD	RT2, g_stackguard0(g)
	MOVD	RT2, g_stackguard1(g)
	MOVD	RT2, (g_stack+stack_lo)(g)
	MOVD	RT1, (g_stack+stack_hi)(g)

	// if there is a _cgo_init, call it using the gcc ABI.
	MOVD	_cgo_init(SB), R12
	CMP	ZR, R12
	BED	nocgo

	MOVD	TLS, O3			// arg 3: TLS base pointer
	MOVD	$runtime·tls_g(SB), O2 	// arg 2: &tls_g
	MOVD	$setg_gcc<>(SB), O1	// arg 1: setg
	MOVD	g, O0			// arg 0: G
	// C functions expect FIXED_FRAME bytes of space on caller stack frame.
	MOVD	BSP, L1
	SUB	$FIXED_FRAME, BSP
	CALL	(R12)
	MOVD	L1, BSP
	
	MOVD	_cgo_init(SB), R12
	CMP	ZR, R12
	BED	nocgo

nocgo:
	// update stackguard after _cgo_init
	MOVD	(g_stack+stack_lo)(g), R3
	ADD	$const__StackGuard, R3
	MOVD	R3, g_stackguard0(g)
	MOVD	R3, g_stackguard1(g)

	// set the per-goroutine and per-mach "registers"
	MOVD	$runtime·m0(SB), R3

	// save m->g0 = g0
	MOVD	g, m_g0(R3)
	// save m0 to g0->m
	MOVD	R3, g_m(g)

	CALL	runtime·check(SB)

	MOVD	BSP, RT1
	MOVW	8(RT1), R3	// copy argc
	MOVW	R3, -8(RT1)
	MOVD	16(RT1), R3		// copy argv
	MOVD	R3, 0(RT1)
	CALL	runtime·args(SB)
	CALL	runtime·osinit(SB)
	CALL	runtime·schedinit(SB)

	// create a new goroutine to start program
	MOVD	$runtime·mainPC(SB), RT1		// entry
	SUB	$(32+FIXED_FRAME), BSP
	MOVD	RT1, FIXED_FRAME+0(BSP)
	MOVD	ZR, FIXED_FRAME+8(BSP)
	MOVD	ZR, FIXED_FRAME+16(BSP)
	MOVD	ZR, FIXED_FRAME+24(BSP)
	CALL	runtime·newproc(SB)
	ADD	$(32+FIXED_FRAME), BSP

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
	// TODO(aram): do we need to initialize FP registers?
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), R3
	MOVD	BSP, R1
	MOVD	R1, gobuf_sp(R3)
	MOVD	LR, gobuf_pc(R3)
	MOVD	g, gobuf_g(R3)
	MOVD	ZR, gobuf_lr(R3)
	MOVD	ZR, gobuf_ret(R3)
	MOVD	ZR, gobuf_ctxt(R3)
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), R5
	MOVD	gobuf_g(R5), g
	CALL	runtime·save_g(SB)

	MOVD	0(g), R4	// make sure g is not nil
	MOVD	gobuf_sp(R5), R1
	MOVD	R1, BSP
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
TEXT runtime·mcall(SB), NOSPLIT|NOFRAME, $0-8
	// Save caller state in g->sched
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
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
	MOVD	TMP, BSP	// sp = m->g0->sched.sp
	SUB	$16, BSP
	MOVD	R3, (176+0)(BSP)
	MOVD	$0, (176+8)(BSP)
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
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
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
	MOVD	R3, BSP

	// call target function
	MOVD	0(CTXT), R3	// code pointer
	CALL	(R3)

	// switch back to g
	MOVD	g_m(g), R3
	MOVD	m_curg(R3), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP
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
TEXT runtime·morestack(SB),NOSPLIT|NOFRAME,$0-0
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
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	LR, (g_sched+gobuf_pc)(g)
	MOVD	R3, (g_sched+gobuf_lr)(g)

	// Called from f.
	// Set m->morebuf to f's callers.
	MOVD	R3, (m_morebuf+gobuf_pc)(R8)	// f's caller's PC
	MOVD	BSP, TMP
	MOVD	TMP, (m_morebuf+gobuf_sp)(R8)	// f's caller's BSP
	MOVD	g, (m_morebuf+gobuf_g)(R8)

	// Call newstack on m->g0's stack.
	MOVD	m_g0(R8), g
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
	MOVD	$runtime·badreflectcall(SB), R1
	JMPL	R1, ZR

#define CALLFN(NAME,MAXSIZE)			\
TEXT NAME(SB), WRAPPER, $MAXSIZE-24;		\
	NO_LOCAL_POINTERS;			\
	/* copy arguments to stack */		\
	MOVD	arg+16(FP), R3;			\
	MOVUW	argsize+24(FP), R4;			\
	MOVD	BSP, R5;				\
	ADD	$(FIXED_FRAME-1), R5;			\
	SUB	$1, R3;				\
	ADD	R5, R4;				\
	CMP	R5, R4;				\
	BED	6(PC);				\
	MOVUB	(R3), R6;			\
	ADD	$1, R3;				\
	MOVUB	R6, (R5);			\
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
	MOVD	BSP, R5;				\
	ADD	R6, R5; 			\
	ADD	R6, R3;				\
	SUB	R6, R4;				\
	ADD	$(FIXED_FRAME-1), R5;			\
	SUB	$1, R3;				\
	ADD	R5, R4;				\
loop:						\
	CMP	R5, R4;				\
	BED	end;				\
	MOVUB	(R5), R6;			\
	ADD	$1, R5;				\
	MOVUB	R6, (R3);			\
	ADD	$1, R3;			\
	JMP	loop;				\
end:						\
	/* execute write barrier updates */	\
	MOVD	argtype+0(FP), R8;		\
	MOVD	arg+16(FP), R3;			\
	MOVUW	n+24(FP), R4;			\
	MOVUW	retoffset+28(FP), R6;		\
	MOVD	R8, (FIXED_FRAME+0)(BSP);			\
	MOVD	R3, (FIXED_FRAME+8)(BSP);			\
	MOVD	R4, (FIXED_FRAME+16)(BSP);			\
	MOVD	R6, (FIXED_FRAME+24)(BSP);			\
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
TEXT runtime·aeshash(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R1
TEXT runtime·aeshash32(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R1
TEXT runtime·aeshash64(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R1
TEXT runtime·aeshashstr(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R1
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	RD	CCR, R2
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT|NOFRAME, $0-16
	MOVD	(8*15)(BSP), R1
	SUB	$4, R1
	MOVD	R1, LR

	MOVD	fv+0(FP), CTXT
	MOVD	argp+8(FP), TMP
	MOVD	TMP, BSP
	SUB	$FIXED_FRAME, BSP
	MOVD	0(CTXT), R3
	JMPL	R3, ZR

// Save state of caller into g->sched.
TEXT gosave<>(SB),NOSPLIT|NOFRAME,$0
	MOVD	LR, (g_sched+gobuf_pc)(g)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	$0, (g_sched+gobuf_ret)(g)
	MOVD	$0, (g_sched+gobuf_ctxt)(g)
	RET

// func asmcgocall(fn, arg unsafe.Pointer) int32
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
TEXT ·asmcgocall(SB),NOSPLIT|NOFRAME,$0-20
	MOVD	fn+0(FP), R3
	MOVD	arg+8(FP), R4

	MOVD	BSP, R10		// save original stack pointer
	MOVD	g, R5

	// Figure out if we need to switch to m->g0 stack.
	// We get called to create new OS threads too, and those
	// come in on the m->g0 stack already.
	MOVD	g_m(g), R6
	MOVD	m_g0(R6), R6
	CMP	R6, g
	BED	g0
	CALL	gosave<>(SB)
	MOVD	R6, g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP

	// Now on a scheduling stack (a pthread-created stack).
g0:
	// Save room for two of our pointers.
	SUB	$16, BSP
	MOVD	R5, -8(BFP)	// save old g on stack
	MOVD	(g_stack+stack_hi)(R5), R5
	SUB	R10, R5
	MOVD	R5, -16(BFP)	// save depth in old g stack (can't just save SP, as stack might be copied during a callback)
	CALL	(R3)
	MOVD	R8, R9

	// Restore g, stack pointer.
	// R8 is errno, so don't touch it
	MOVD	-8(BFP), g
	MOVD    (g_stack+stack_hi)(g), R5
	MOVD    -16(BFP), R6
	SUB     R6, R5
	MOVD    24(R5), R2
	CALL	runtime·save_g(SB)
	MOVD	(g_stack+stack_hi)(g), R5
	MOVD	-16(BFP), R6
	SUB	R6, R5
	MOVD	R5, BSP

	MOVW	R8, ret+16(FP)
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),NOSPLIT,$24-24
	MOVD	$fn+0(FP), R1
	MOVD	R1, (FIXED_FRAME+0)(BSP)
	MOVD	frame+8(FP), R1
	MOVD	R1, (FIXED_FRAME+8)(BSP)
	MOVD	framesize+16(FP), R1
	MOVD	R1, (FIXED_FRAME+16)(BSP)
	MOVD	$runtime·cgocallback_gofunc(SB), R1
	CALL	R1
	RET

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.go for more details.
TEXT ·cgocallback_gofunc(SB),NOSPLIT,$24-24
	NO_LOCAL_POINTERS

	// Load m and g from thread-local storage.
	MOVB	runtime·iscgo(SB), R3
	CMP	R3, ZR
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

	MOVD	g_m(g), R8
	MOVD	R8, savedm-8(SP)
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
	MOVD	g_m(g), R8
	MOVD	m_g0(R8), R3
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(R3)

havem:
	// Now there's a valid m, and we're running on its m->g0.
	// Save current m->g0->sched.sp on stack and then set it to SP.
	// Save current sp in m->g0->sched.sp in preparation for
	// switch back to m->curg stack.
	// NOTE: unwindm knows that the saved g->sched.sp is at 8(R1) aka savedsp-16(SP).
	MOVD	m_g0(R8), R3
	MOVD	(g_sched+gobuf_sp)(R3), R4
	MOVD	R4, savedsp-16(SP)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(R3)

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
	MOVD	m_curg(R8), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), R4 // prepare stack as R4
	MOVD	(g_sched+gobuf_pc)(g), R5
	MOVD	R5, -(FIXED_FRAME+16)(R4)
	MOVD	$-(FIXED_FRAME+16)(R4), TMP
	MOVD	TMP, BSP
	CALL	runtime·cgocallbackg(SB)

	// Restore g->sched (== m->curg->sched) from saved values.
	MOVD	0(BSP), R5
	MOVD	R5, (g_sched+gobuf_pc)(g)
	MOVD	$(FIXED_FRAME+16)(BSP), R4
	MOVD	R4, (g_sched+gobuf_sp)(g)

	// Switch back to m->g0's stack and restore m->g0->sched.sp.
	// (Unlike m->curg, the g0 goroutine never uses sched.pc,
	// so we do not have to restore it.)
	MOVD	g_m(g), R8
	MOVD	m_g0(R8), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP
	MOVD	savedsp-16(SP), R4
	MOVD	R4, (g_sched+gobuf_sp)(g)

	// If the m on entry was nil, we called needm above to borrow an m
	// for the duration of the call. Since the call is over, return it with dropm.
	MOVD	savedm-8(SP), R6
	CMP	R6, ZR
	BNED	droppedm
	MOVD	$runtime·dropm(SB), RT1
	CALL	(RT1)
droppedm:

	// Done!
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
	MOVD	FIXED_FRAME+8*15(BFP), R3		// LR saved by caller
	MOVD	runtime·stackBarrierPC(SB), R4
	CMP	R4, R3
	BNED	nobar
	// Get original return PC.
	CALL	runtime·nextBarrierPC(SB)
	MOVD	FIXED_FRAME+0(R1), R3
nobar:
	MOVD	R3, ret+8(FP)
	RET

TEXT runtime·setcallerpc(SB),NOSPLIT,$8-16
	MOVD	pc+8(FP), R3
	MOVD	FIXED_FRAME+8(BSP), R4
	MOVD	runtime·stackBarrierPC(SB), R5
	CMP	R4, R5
	BED	setbar
	MOVD	R3, FIXED_FRAME+8*15(BFP)		// set LR in caller
	RET
setbar:
	// Set the stack barrier return PC.
	MOVD	R3, FIXED_FRAME+0(R1)
	CALL	runtime·setNextBarrierPC(SB)
	RET

TEXT runtime·getcallersp(SB),NOSPLIT,$0-16
	MOVD	argp+0(FP), R1
	SUB	$FIXED_FRAME, R1
	MOVD	R1, ret+8(FP)
	RET

TEXT runtime·abort(SB),NOSPLIT|NOFRAME,$0-0
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
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD	p+0(FP), R3
	MOVD	h+8(FP), R4
	MOVD	8(CTXT), R5
	MOVD	R3, FIXED_FRAME+0(R1)
	MOVD	R4, FIXED_FRAME+8(R1)
	MOVD	R5, FIXED_FRAME+16(R1)
	CALL	runtime·memhash(SB)
	MOVD	FIXED_FRAME+24(R1), R3
	MOVD	R3, ret+16(FP)
	RET

// memequal(p, q unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT|NOFRAME,$0-25
	MOVD	a+0(FP), R1
	MOVD	b+8(FP), R2
	MOVD	size+16(FP), R3
	ADD	R1, R3, R6
	MOVD	$1, TMP
	MOVB	TMP, ret+24(FP)
	CMP	R1, R2
	BED	done
loop:
	CMP	R1, R6
	BED	done
	MOVUB	(R1), R4
	ADD	$1, R1
	MOVUB	(R2), R5
	ADD $1, R2
	CMP	R4, R5
	BED	loop

	MOVB	ZR, ret+24(FP)
done:
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$40-17
	MOVD	a+0(FP), R3
	MOVD	b+8(FP), R4
	CMP	R3, R4
	BED	eq
	MOVD	8(CTXT), R5    // compiler stores size at offset 8 in the closure
	MOVD	R3, FIXED_FRAME+0(R1)
	MOVD	R4, FIXED_FRAME+8(R1)
	MOVD	R5, FIXED_FRAME+16(R1)
	CALL	runtime·memequal(SB)
	MOVD	$FIXED_FRAME+24(BSP), R3
	MOVUB	(R3), R3
	MOVB	R3, ret+16(FP)
	RET
eq:
	MOVD	$1, R3
	MOVB	R3, ret+16(FP)
	RET

// eqstring tests whether two strings are equal.
// The compiler guarantees that strings passed
// to eqstring have equal length.
// See runtime_test.go:eqstring_generic for
// equivalent Go code.
TEXT runtime·eqstring(SB),NOSPLIT,$0-33
	MOVD	s1str+0(FP), R1
	MOVD	s1len+8(FP), R2
	MOVD	s2str+16(FP), R3
	ADD	R1, R2		// end
loop:
	CMP	R1, R2
	BED	equal		// reaches the end
	MOVUB	(R1), R4
	ADD	$1, R1
	MOVUB	(R3), R5
	ADD	$1, R3
	CMP	R4, R5
	BED	loop
notequal:
	MOVB	ZR, ret+32(FP)
	RET
equal:
	MOVD	$1, R1
	MOVB	R1, ret+32(FP)
	RET

//
// functions for other packages
//
TEXT bytes·IndexByte(SB),NOSPLIT,$0-40
	MOVD	s+0(FP), R3
	MOVD	s_len+8(FP), R4
	MOVUB	c+24(FP), R5	// byte to find
	MOVD	R3, R6		// store base for later
	SUB	$1, R3
	ADD	R3, R4		// end-1

loop:
	CMP	R3, R4
	BED	notfound
	MOVUB	(R3), R8
	ADD	$1, R3
	CMP	R5, R8
	BNEW	loop

	SUB	R6, R3		// remove base
	MOVD	R3, ret+32(FP)
	RET

notfound:
	MOVD	$-1, R3
	MOVD	R3, ret+32(FP)
	RET

TEXT strings·IndexByte(SB),NOSPLIT,$0-32
	MOVD	p+0(FP), R3
	MOVD	b_len+8(FP), R4
	MOVUB	c+16(FP), R5	// byte to find
	MOVD	R3, R6		// store base for later
	SUB	$1, R3
	ADD	R3, R4		// end-1

loop:
	CMP	R3, R4
	BED	notfound
	MOVUB	(R3), R8
	ADD	$1, R3
	CMP	R5, R8
	BNEW	loop

	SUB	R6, R3		// remove base
	MOVD	R3, ret+24(FP)
	RET

notfound:
	MOVD	$-1, R3
	MOVD	R3, ret+24(FP)
	RET

// TODO: share code with memequal?
TEXT bytes·Equal(SB),NOSPLIT,$0-49
	MOVD	a_len+8(FP), R3
	MOVD	b_len+32(FP), R4

	CMP	R3, R4		// unequal lengths are not equal
	BNED	noteq

	MOVD	a+0(FP), R5
	MOVD	b+24(FP), R6
	SUB	$1, R5
	SUB	$1, R6
	ADD	R5, R3		// end-1

loop:
	CMP	R5, R3
	BED	equal		// reached the end
	MOVUB	(R5), R4
	ADD	$1, R5
	MOVUB	(R6), R8
	ADD	$1, R6
	CMP	R4, R8
	BEW	loop

noteq:
	MOVB	ZR, ret+48(FP)
	RET

equal:
	MOVD	$1, R3
	MOVB	R3, ret+48(FP)
	RET

TEXT runtime·fastrand1(SB),NOSPLIT|NOFRAME,$0-4
	MOVD	g_m(g), R4
	MOVUW	m_fastrand(R4), R3
	ADD	R3, R3
	CMP	ZR, R3
	BGEW	2(PC)
	XOR	$0x88888eef, R3
	MOVW	R3, m_fastrand(R4)
	MOVW	R3, ret+0(FP)
	RET

TEXT runtime·return0(SB), NOSPLIT, $0
	MOVW	ZR, R8
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT|NOFRAME,$0-0
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
	MOVD	runtime·lastmoduledatap(SB), R1
	MOVD	R8, moduledata_next(R1)
	MOVD	R8, runtime·lastmoduledatap(SB)
	RET

TEXT ·checkASM(SB),NOSPLIT,$0-1
	OR	$1, ZR, R3
	MOVB	R3, ret+0(FP)
	RET

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"
#include "asm_sparc64.h"

DATA dbgbuf(SB)/8, $"\n\n"
GLOBL dbgbuf(SB), $8

TEXT runtime·rt0_go(SB),NOSPLIT,$16-0
	// BSP = stack; I0 = argc; I1 = argv

	// initialize essential registers
	CALL	runtime·reginit(SB)

	MOVW	I0, L1	// argc
	MOVD	I1, L2	// argv

	// create istack out of the given (operating system) stack.
	// _cgo_init may update stackguard.
	MOVD	$runtime·g0(SB), g
	MOVD BSP, RT1
	MOVD	$(-64*1024)(BSP), RT2
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
	MOVD	g, L3
	CALL	(R12)
	MOVD	L3, g

nocgo:
	// update stackguard after _cgo_init
	MOVD	(g_stack+stack_lo)(g), R25
	ADD	$const__StackGuard, R25
	MOVD	R25, g_stackguard0(g)
	MOVD	R25, g_stackguard1(g)

	// set the per-goroutine and per-mach "registers"
	MOVD	$runtime·m0(SB), R25

	// save m->g0 = g0
	MOVD	g, m_g0(R25)
	// save m0 to g0->m
	MOVD	R25, g_m(g)

	CALL	runtime·check(SB)

	MOVD	L1, FIXED_FRAME+0(BSP)	// copy argc
	MOVD	L2, FIXED_FRAME+8(BSP)	// copy argv
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
	// TODO(aram): do we need to initialize FP registers?
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), R25
	MOVD	BSP, R27
	MOVD	R27, gobuf_sp(R25)
	MOVD	OLR, gobuf_pc(R25)
	MOVD	g, gobuf_g(R25)
	MOVD	ZR, gobuf_lr(R25)
	MOVD	ZR, gobuf_ret(R25)
	MOVD	ZR, gobuf_ctxt(R25)
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT|NOFRAME, $0-8
	MOVD	buf+0(FP), R22
	MOVD	gobuf_g(R22), g
	CALL	runtime·save_g(SB)

	MOVD	0(g), R28	// make sure g is not nil
	MOVD	gobuf_sp(R22), R27
	MOVD	R27, BSP
	MOVD	gobuf_lr(R22), OLR
	MOVD	gobuf_ret(R22), R27
	MOVD	gobuf_ctxt(R22), CTXT
	MOVD	ZR, gobuf_sp(R22)
	MOVD	ZR, gobuf_ret(R22)
	MOVD	ZR, gobuf_lr(R22)
	MOVD	ZR, gobuf_ctxt(R22)
	CMP	ZR, ZR // set condition codes for == test, needed by stack split
	MOVD	gobuf_pc(R22), R8
	JMPL	R8, ZR

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

	// Switch to m->g0 & its stack, call fn.
	MOVD	g, R25
	MOVD	g_m(g), R8
	MOVD	m_g0(R8), g
	CALL	runtime·save_g(SB)
	CMP	g, R25
	BNED	ok
	JMP	runtime·badmcall(SB)
ok:
	MOVD	fn+0(FP), CTXT			// context
	MOVD	0(CTXT), R28			// code pointer
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP	// sp = m->g0->sched.sp
	SUB	$16, BSP
	MOVD	R25, (176+0)(BSP)
	MOVD	$0, (176+8)(BSP)
	CALL	(R28)
	JMP	runtime·badmcall2(SB)

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
	UNDEF
	CALL	(ILR)	// make sure this function is not leaf
	RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
	MOVD	fn+0(FP), R25	// R25 = fn
	MOVD	R25, CTXT		// context
	MOVD	g_m(g), R28	// R28 = m

	MOVD	m_gsignal(R28), R22	// R22 = gsignal
	CMP	g, R22
	BED	noswitch

	MOVD	m_g0(R28), R22	// R22 = g0
	CMP	g, R22
	BED	noswitch

	MOVD	m_curg(R28), R8
	CMP	g, R8
	BED	switch

	// Bad: g is not gsignal, not g0, not curg. What is it?
	// Hide call from linker nosplit analysis.
	MOVD	$runtime·badsystemstack(SB), R25
	CALL	(R25)

switch:
	// save our state in g->sched. Pretend to
	// be systemstack_switch if the G stack is scanned.
	MOVD	$runtime·systemstack_switch(SB), R8
	ADD	$8, R8	// get past prologue
	MOVD	R8, (g_sched+gobuf_pc)(g)
	MOVD	BSP, TMP
	MOVD	TMP, (g_sched+gobuf_sp)(g)
	MOVD	$0, (g_sched+gobuf_lr)(g)
	MOVD	g, (g_sched+gobuf_g)(g)

	// switch to g0
	MOVD	R22, g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), R25
	// make it look like mstart called systemstack on g0, to stop traceback
	SUB	$16, R25
	AND	$~15, R25
	MOVD	$runtime·mstart(SB), R28
	MOVD	R28, 0(R25)
	MOVD	R25, BSP

	// call target function
	MOVD	0(CTXT), R25	// code pointer
	CALL	(R25)

	// switch back to g
	MOVD	g_m(g), R25
	MOVD	m_curg(R25), g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP
	MOVD	$0, (g_sched+gobuf_sp)(g)
	RET

noswitch:
	// already on m stack, just call directly
	MOVD	0(CTXT), R25	// code pointer
	CALL	(R25)
	RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
// Caller has already loaded:
// R25 prolog's LR
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB),NOSPLIT|NOFRAME,$0-0
	UNDEF

TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	ZR, CTXT
	JMP	runtime·morestack(SB)

TEXT runtime·stackBarrier(SB),NOSPLIT,$0
	RET

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
	MOVD	$runtime·badreflectcall(SB), R27
	JMPL	R27, ZR

#define CALLFN(NAME,MAXSIZE)			\
TEXT NAME(SB), WRAPPER, $MAXSIZE-24;		\
	NO_LOCAL_POINTERS;			\
	/* copy arguments to stack */		\
	MOVD	arg+16(FP), R25;			\
	MOVUW	argsize+24(FP), R28;			\
	MOVD	BSP, R22;				\
	ADD	$(FIXED_FRAME-1), R22;			\
	SUB	$1, R25;				\
	ADD	R22, R28;				\
	CMP	R22, R28;				\
	BED	6(PC);				\
	MOVUB	(R25), R9;			\
	ADD	$1, R25;				\
	MOVUB	R9, (R22);			\
	ADD	$1, R22;				\
	JMP	-6(PC);				\
	/* call function */			\
	MOVD	f+8(FP), CTXT;			\
	MOVD	(CTXT), R27;			\
	PCDATA  $PCDATA_StackMapIndex, $0;	\
	CALL	(R27);				\
	/* copy return values back */		\
	MOVD	arg+16(FP), R25;			\
	MOVUW	n+24(FP), R28;			\
	MOVUW	retoffset+28(FP), R9;		\
	MOVD	BSP, R22;				\
	ADD	R9, R22; 			\
	ADD	R9, R25;				\
	SUB	R9, R28;				\
	ADD	$(FIXED_FRAME-1), R22;			\
	SUB	$1, R25;				\
	ADD	R22, R28;				\
loop:						\
	CMP	R22, R28;				\
	BED	end;				\
	MOVUB	(R22), R9;			\
	ADD	$1, R22;				\
	MOVUB	R9, (R25);			\
	ADD	$1, R25;			\
	JMP	loop;				\
end:						\
	/* execute write barrier updates */	\
	MOVD	argtype+0(FP), R8;		\
	MOVD	arg+16(FP), R25;			\
	MOVUW	n+24(FP), R28;			\
	MOVUW	retoffset+28(FP), R9;		\
	MOVD	R8, (FIXED_FRAME+0)(BSP);			\
	MOVD	R25, (FIXED_FRAME+8)(BSP);			\
	MOVD	R28, (FIXED_FRAME+16)(BSP);			\
	MOVD	R9, (FIXED_FRAME+24)(BSP);			\
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
	MOVW	(ZR), R27
TEXT runtime·aeshash32(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R27
TEXT runtime·aeshash64(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R27
TEXT runtime·aeshashstr(SB),NOSPLIT|NOFRAME,$0-0
	MOVW	(ZR), R27
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	RD	CCR, R29
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT|NOFRAME, $0-16
	MOVD	(8*15)(BSP), R27
	SUB	$4, R27
	MOVD	R27, OLR

	MOVD	fv+0(FP), CTXT
	MOVD	argp+8(FP), TMP
	MOVD	TMP, BSP
	SUB	$FIXED_FRAME, BSP
	MOVD	0(CTXT), R25
	JMPL	R25, ZR

// Save state of caller into g->sched.
TEXT gosave<>(SB),NOSPLIT|NOFRAME,$0
	MOVD	OLR, (g_sched+gobuf_pc)(g)
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
	MOVD	OLR, R21	// save LR
	MOVD	fn+0(FP), R25
	MOVD	arg+8(FP), O0
	MOVD	O0, FIXED_FRAME(BSP)

	MOVD	BSP, R10		// save original stack pointer
	MOVD	g, R22

	// Figure out if we need to switch to m->g0 stack.
	// We get called to create new OS threads too, and those
	// come in on the m->g0 stack already.
	MOVD	g_m(g), R9
	MOVD	m_g0(R9), R9
	CMP	R9, g
	BED	g0
	CALL	gosave<>(SB)
	MOVD	R9, g
	CALL	runtime·save_g(SB)
	MOVD	(g_sched+gobuf_sp)(g), TMP
	MOVD	TMP, BSP

	// Now on a scheduling stack (a pthread-created stack).
g0:
	MOVD	R22, R19	// save old g
	MOVD	(g_stack+stack_hi)(R22), R22
	SUB	R10, R22
	MOVD	R22, R20	// save depth in old g stack (can't just save SP, as stack might be copied during a callback)
	CALL	(R25)
	MOVD	R8, R9

	// Restore g, stack pointer.
	// R8 is errno, so don't touch it
	MOVD	R19, g
	MOVD    (g_stack+stack_hi)(g), R22
	SUB     R20, R22
	MOVD    24(R22), R29
	CALL	runtime·save_g(SB)
	MOVD    (g_stack+stack_hi)(g), R22
	SUB     R20, R22
	MOVD	R22, BSP

	MOVD	R21, OLR
	MOVW	R8, ret+16(FP)
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),NOSPLIT,$32-24
	UNDEF

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.go for more details.
TEXT ·cgocallback_gofunc(SB),NOSPLIT,$32-24
	UNDEF

// Called from cgo wrappers, this function returns g->m->curg.stack.hi.
// Must obey the gcc calling convention.
TEXT _cgo_topofstack(SB),NOSPLIT,$32
	// g and TMP might be clobbered by load_g. They
	// are callee-save in the gcc calling convention, so save them.
	MOVD	TMP, savedTMP-8(SP)
	MOVD	g, saveG-16(SP)

	CALL	runtime·load_g(SB)
	MOVD	g_m(g), R27
	MOVD	m_curg(R27), R27
	MOVD	(g_stack+stack_hi)(R27), R27

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
TEXT setg_gcc<>(SB),NOSPLIT,$16
	MOVD	R8, g
	MOVD	TMP, savedTMP-8(SP)
	CALL	runtime·save_g(SB)
	MOVD	savedTMP-8(SP), TMP
	RET

TEXT runtime·getcallerpc(SB),NOSPLIT,$16-16
	MOVD	FIXED_FRAME+8*15(BFP), R25		// LR saved by caller
	MOVD	runtime·stackBarrierPC(SB), R28
	CMP	R28, R25
	BNED	nobar
	// Get original return PC.
	CALL	runtime·nextBarrierPC(SB)
	MOVD	FIXED_FRAME+0(R27), R25
nobar:
	MOVD	R25, ret+8(FP)
	RET

TEXT runtime·setcallerpc(SB),NOSPLIT,$16-16
	MOVD	pc+8(FP), R25
	MOVD	FIXED_FRAME+8(BSP), R28
	MOVD	runtime·stackBarrierPC(SB), R22
	CMP	R28, R22
	BED	setbar
	MOVD	R25, FIXED_FRAME+8*15(BFP)		// set LR in caller
	RET
setbar:
	// Set the stack barrier return PC.
	MOVD	R25, FIXED_FRAME+0(R27)
	CALL	runtime·setNextBarrierPC(SB)
	RET

TEXT runtime·getcallersp(SB),NOSPLIT,$0-16
	MOVD	argp+0(FP), R27
	SUB	$FIXED_FRAME, R27
	MOVD	R27, ret+8(FP)
	RET

TEXT runtime·abort(SB),NOSPLIT|NOFRAME,$0-0
	JMPL	ZR, ZR
	UNDEF

// func cputicks() int64
TEXT runtime·cputicks(SB),NOSPLIT,$0-0
	RD	TICK, R27
	MOVD	R27, ret+0(FP)
	RET

// memhash_varlen(p unsafe.Pointer, h seed) uintptr
// redirects to memhash(p, h, size) using the size
// stored in the closure.
TEXT runtime·memhash_varlen(SB),NOSPLIT,$48-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD	p+0(FP), R25
	MOVD	h+8(FP), R28
	MOVD	8(CTXT), R22
	MOVD	R25, FIXED_FRAME+0(R27)
	MOVD	R28, FIXED_FRAME+8(R27)
	MOVD	R22, FIXED_FRAME+16(R27)
	CALL	runtime·memhash(SB)
	MOVD	FIXED_FRAME+24(R27), R25
	MOVD	R25, ret+16(FP)
	RET

// memequal(p, q unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT|NOFRAME,$0-25
	MOVD	a+0(FP), R27
	MOVD	b+8(FP), R29
	MOVD	size+16(FP), R25
	ADD	R27, R25, R9
	MOVD	$1, TMP
	MOVB	TMP, ret+24(FP)
	CMP	R27, R29
	BED	done
loop:
	CMP	R27, R9
	BED	done
	MOVUB	(R27), R28
	ADD	$1, R27
	MOVUB	(R29), R22
	ADD $1, R29
	CMP	R28, R22
	BED	loop

	MOVB	ZR, ret+24(FP)
done:
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$48-17
	MOVD	a+0(FP), R25
	MOVD	b+8(FP), R28
	CMP	R25, R28
	BED	eq
	MOVD	8(CTXT), R22    // compiler stores size at offset 8 in the closure
	MOVD	R25, FIXED_FRAME+0(R27)
	MOVD	R28, FIXED_FRAME+8(R27)
	MOVD	R22, FIXED_FRAME+16(R27)
	CALL	runtime·memequal(SB)
	MOVD	$FIXED_FRAME+24(BSP), R25
	MOVUB	(R25), R25
	MOVB	R25, ret+16(FP)
	RET
eq:
	MOVD	$1, R25
	MOVB	R25, ret+16(FP)
	RET

// eqstring tests whether two strings are equal.
// The compiler guarantees that strings passed
// to eqstring have equal length.
// See runtime_test.go:eqstring_generic for
// equivalent Go code.
TEXT runtime·eqstring(SB),NOSPLIT,$0-33
	MOVD	s1str+0(FP), R27
	MOVD	s1len+8(FP), R29
	MOVD	s2str+16(FP), R25
	ADD	R27, R29		// end
loop:
	CMP	R27, R29
	BED	equal		// reaches the end
	MOVUB	(R27), R28
	ADD	$1, R27
	MOVUB	(R25), R22
	ADD	$1, R25
	CMP	R28, R22
	BED	loop
notequal:
	MOVB	ZR, ret+32(FP)
	RET
equal:
	MOVD	$1, R27
	MOVB	R27, ret+32(FP)
	RET

//
// functions for other packages
//
TEXT bytes·IndexByte(SB),NOSPLIT,$0-40
	MOVD	s+0(FP), R25
	MOVD	s_len+8(FP), R28
	MOVUB	c+24(FP), R22	// byte to find
	MOVD	R25, R9		// store base for later
	SUB	$1, R25
	ADD	R25, R28		// end-1

loop:
	CMP	R25, R28
	BED	notfound
	MOVUB	(R25), R8
	ADD	$1, R25
	CMP	R22, R8
	BNEW	loop

	SUB	R9, R25		// remove base
	MOVD	R25, ret+32(FP)
	RET

notfound:
	MOVD	$-1, R25
	MOVD	R25, ret+32(FP)
	RET

TEXT strings·IndexByte(SB),NOSPLIT,$0-32
	MOVD	p+0(FP), R25
	MOVD	b_len+8(FP), R28
	MOVUB	c+16(FP), R22	// byte to find
	MOVD	R25, R9		// store base for later
	SUB	$1, R25
	ADD	R25, R28		// end-1

loop:
	CMP	R25, R28
	BED	notfound
	MOVUB	(R25), R8
	ADD	$1, R25
	CMP	R22, R8
	BNEW	loop

	SUB	R9, R25		// remove base
	MOVD	R25, ret+24(FP)
	RET

notfound:
	MOVD	$-1, R25
	MOVD	R25, ret+24(FP)
	RET

// TODO: share code with memequal?
TEXT bytes·Equal(SB),NOSPLIT,$0-49
	MOVD	a_len+8(FP), R25
	MOVD	b_len+32(FP), R28

	CMP	R25, R28		// unequal lengths are not equal
	BNED	noteq

	MOVD	a+0(FP), R22
	MOVD	b+24(FP), R9
	SUB	$1, R22
	SUB	$1, R9
	ADD	R22, R25		// end-1

loop:
	CMP	R22, R25
	BED	equal		// reached the end
	MOVUB	(R22), R28
	ADD	$1, R22
	MOVUB	(R9), R8
	ADD	$1, R9
	CMP	R28, R8
	BEW	loop

noteq:
	MOVB	ZR, ret+48(FP)
	RET

equal:
	MOVD	$1, R25
	MOVB	R25, ret+48(FP)
	RET

TEXT runtime·fastrand1(SB),NOSPLIT|NOFRAME,$0-4
	MOVD	g_m(g), R28
	MOVUW	m_fastrand(R28), R25
	ADD	R25, R25
	CMP	ZR, R25
	BGEW	2(PC)
	XOR	$0x88888eef, R25
	MOVW	R25, m_fastrand(R28)
	MOVW	R25, ret+0(FP)
	RET

TEXT runtime·return0(SB), NOSPLIT, $0
	MOVW	ZR, R8
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	R27, R27	// NOP
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

TEXT ·checkASM(SB),NOSPLIT,$0-1
	OR	$1, ZR, R25
	MOVB	R25, ret+0(FP)
	RET

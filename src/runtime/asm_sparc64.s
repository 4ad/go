// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "funcdata.h"
#include "textflag.h"

DATA dbgbuf(SB)/2, $"\n\n"
GLOBL dbgbuf(SB), 16, $2

TEXT runtime·rt0_go(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$1, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

DATA	runtime·mainPC+0(SB)/8,$runtime·main(SB)
GLOBL	runtime·mainPC(SB),RODATA,$8

TEXT runtime·breakpoint(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$2, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·asminit(SB),NOSPLIT,$-8-0
	RET

TEXT runtime·reginit(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$3, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), NOSPLIT, $-8-8
	// TODO(aram):
	MOVD	$4, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT, $-8-8
	// TODO(aram):
	MOVD	$5, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// void mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT, $-8-8
	// TODO(aram):
	MOVD	$6, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
	// TODO(aram):
	MOVD	$7, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
	// TODO(aram):
	MOVD	$8, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
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
	MOVD	$9, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·morestack_noctxt(SB),NOSPLIT,$-4-0
	// TODO(aram):
	MOVD	$10, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·stackBarrier(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$11, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// reflectcall: call a function with the given argument list
// func call(argtype *_type, f *FuncVal, arg *byte, argsize, retoffset uint32).
// we don't have variable-sized frames, so we use a small number
// of constant-sized-frame functions to encode a few bits of size in the pc.
// Caution: ugly multiline assembly macros in your future!

TEXT reflect·call(SB), NOSPLIT, $0-0
	B	·reflectcall(SB)

TEXT ·reflectcall(SB), NOSPLIT, $-8-32
	// TODO(aram):
	MOVD	$12, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// AES hashing not implemented for SPARC64.
TEXT runtime·aeshash(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$13, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET
TEXT runtime·aeshash32(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$14, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET
TEXT runtime·aeshash64(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$15, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET
TEXT runtime·aeshashstr(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$16, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET
	
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	// TODO(aram):
	MOVD	$17, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// void jmpdefer(fv, sp);
// called from deferreturn.
// 1. grab stored LR for caller
// 2. sub 4 bytes to get back to BL deferreturn
// 3. BR to fn
TEXT runtime·jmpdefer(SB), NOSPLIT, $-8-16
	// TODO(aram):
	MOVD	$18, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// Save state of caller into g->sched. Smashes R0.
TEXT gosave<>(SB),NOSPLIT,$-8
	// TODO(aram):
	MOVD	$19, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// func asmcgocall(fn, arg unsafe.Pointer) int32
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
TEXT ·asmcgocall(SB),NOSPLIT,$0-20
	// TODO(aram):
	MOVD	$20, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),NOSPLIT,$24-24
	// TODO(aram):
	MOVD	$21, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.go for more details.
TEXT ·cgocallback_gofunc(SB),NOSPLIT,$24-24
	// TODO(aram):
	MOVD	$22, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// Called from cgo wrappers, this function returns g->m->curg.stack.hi.
// Must obey the gcc calling convention.
TEXT _cgo_topofstack(SB),NOSPLIT,$24
	// TODO(aram):
	MOVD	$23, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// void setg(G*); set g. for use by needm.
TEXT runtime·setg(SB), NOSPLIT, $0-8
	// TODO(aram):
	MOVD	$24, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$8
	// TODO(aram):
	MOVD	$25, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·getcallerpc(SB),NOSPLIT,$8-16
	// TODO(aram):
	MOVD	$26, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·setcallerpc(SB),NOSPLIT,$8-16
	// TODO(aram):
	MOVD	$27, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·getcallersp(SB),NOSPLIT,$0-16
	// TODO(aram):
	MOVD	$28, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·abort(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$29, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

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
	MOVD	$30, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// memequal(p, q unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT,$-8-25
	// TODO(aram):
	MOVD	$31, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$40-17
	// TODO(aram):
	MOVD	$32, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·cmpstring(SB),NOSPLIT,$-4-40
	// TODO(aram):
	MOVD	$33, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT bytes·Compare(SB),NOSPLIT,$-4-56
	// TODO(aram):
	MOVD	$34, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
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
TEXT runtime·cmpbody<>(SB),NOSPLIT,$-4-0
	// TODO(aram):
	MOVD	$35, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// eqstring tests whether two strings are equal.
// The compiler guarantees that strings passed
// to eqstring have equal length.
// See runtime_test.go:eqstring_generic for
// equivalent Go code.
TEXT runtime·eqstring(SB),NOSPLIT,$0-33
	// TODO(aram):
	MOVD	$36, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

//
// functions for other packages
//
TEXT bytes·IndexByte(SB),NOSPLIT,$0-40
	// TODO(aram):
	MOVD	$37, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT strings·IndexByte(SB),NOSPLIT,$0-32
	// TODO(aram):
	MOVD	$38, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// TODO: share code with memequal?
TEXT bytes·Equal(SB),NOSPLIT,$0-49
	// TODO(aram):
	MOVD	$39, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·fastrand1(SB),NOSPLIT,$-8-4
	// TODO(aram):
	MOVD	$40, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT runtime·return0(SB), NOSPLIT, $0
	// TODO(aram):
	MOVD	$41, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT,$-8-0
	// TODO(aram):
	MOVD	$42, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
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
	MOVD	$43, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

TEXT ·checkASM(SB),NOSPLIT,$0-1
	// TODO(aram):
	MOVD	$44, TMP
	ADD	$'!', TMP, TMP
	MOVD	TMP, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf, R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), TMP
	CALL	TMP
	UNDEF
	RET

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#define STACK_BIAS 0x7ff

/*
 * Apple still insists on underscore prefixes for C function names.
 */
#if defined(__APPLE__)
#define EXT(s) _##s
#else
#define EXT(s) s
#endif

// Apple's ld64 wants 4-byte alignment for ARM code sections.
// .align in both Apple as and GNU as treat n as aligning to 2**n bytes.
.align	2

/*
 * void crosscall1(void (*fn)(void), void (*setg_gcc)(void *g), void *g)
 *
 * Calling into the gc tool chain, where all registers are caller save.
 * Called from standard SPARC ABI, where r16-r29 are callee-save, so they
 * must be saved explicitly.
 */
.globl EXT(crosscall1)
EXT(crosscall1):
	save  %sp, -304, %sp

	stx %l0, [%fp-STACK_BIAS-0*8]
	stx %l1, [%fp-STACK_BIAS-1*8]
	stx %l2, [%fp-STACK_BIAS-2*8]
	stx %l3, [%fp-STACK_BIAS-3*8]
	stx %l4, [%fp-STACK_BIAS-4*8]
	stx %l5, [%fp-STACK_BIAS-5*8]
	stx %l6, [%fp-STACK_BIAS-6*8]
	stx %l7, [%fp-STACK_BIAS-7*8]
	stx %i0, [%fp-STACK_BIAS-64-0*8]
	stx %i1, [%fp-STACK_BIAS-64-1*8]
	stx %i2, [%fp-STACK_BIAS-64-2*8]
	stx %i3, [%fp-STACK_BIAS-64-3*8]
	stx %i4, [%fp-STACK_BIAS-64-4*8]
	stx %i5, [%fp-STACK_BIAS-64-5*8]

	mov	%i0, %l0
	mov	%i1, %l1
	mov	%i2, %i0

	call %l1
	nop
	call %l0
	nop

	stx %i5, [%fp-STACK_BIAS-64-5*8]
	stx %i4, [%fp-STACK_BIAS-64-4*8]
	stx %i3, [%fp-STACK_BIAS-64-3*8]
	stx %i2, [%fp-STACK_BIAS-64-2*8]
	stx %i1, [%fp-STACK_BIAS-64-1*8]
	stx %i0, [%fp-STACK_BIAS-64-0*8]
	stx %l7, [%fp-STACK_BIAS-7*8]
	stx %l6, [%fp-STACK_BIAS-6*8]
	stx %l5, [%fp-STACK_BIAS-5*8]
	stx %l4, [%fp-STACK_BIAS-4*8]
	stx %l3, [%fp-STACK_BIAS-3*8]
	stx %l2, [%fp-STACK_BIAS-2*8]
	stx %l1, [%fp-STACK_BIAS-1*8]
	stx %l0, [%fp-STACK_BIAS-0*8]

	rett  %i7 + 8
	nop

#ifdef __ELF__
.section .note.GNU-stack,"",%progbits
#endif

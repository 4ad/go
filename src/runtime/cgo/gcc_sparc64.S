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
 * Called from standard SPARC ABI, where r16-r31 are callee-save, so they
 * must be saved explicitly.
 */
.globl EXT(crosscall1)
EXT(crosscall1):
	stx %l0, [%06+STACK_BIAS+0*8]
	stx %l1, [%06+STACK_BIAS+1*8]
	stx %l2, [%06+STACK_BIAS+2*8]
	stx %l3, [%06+STACK_BIAS+3*8]
	stx %l4, [%06+STACK_BIAS+4*8]
	stx %l5, [%06+STACK_BIAS+5*8]
	stx %l6, [%06+STACK_BIAS+6*8]
	stx %l7, [%06+STACK_BIAS+7*8]
	stx %i0, [%06+STACK_BIAS+64+0*8]
	stx %i1, [%06+STACK_BIAS+64+1*8]
	stx %i2, [%06+STACK_BIAS+64+2*8]
	stx %i3, [%06+STACK_BIAS+64+3*8]
	stx %i4, [%06+STACK_BIAS+64+4*8]
	stx %i5, [%06+STACK_BIAS+64+5*8]
	stx %i6, [%06+STACK_BIAS+64+6*8]
	stx %i7, [%06+STACK_BIAS+64+7*8]

	mov	%i0, %l0
	mov	%i1, %l1
	mov	%i2, %i0

	call %l1
	call %l0

	ldx [%06+STACK_BIAS+64+7*8], %i7
	ldx [%06+STACK_BIAS+64+6*8], %i6
	ldx [%06+STACK_BIAS+64+5*8], %i5
	ldx [%06+STACK_BIAS+64+4*8], %i4
	ldx [%06+STACK_BIAS+64+3*8], %i3
	ldx [%06+STACK_BIAS+64+2*8], %i2
	ldx [%06+STACK_BIAS+64+1*8], %i1
	ldx [%06+STACK_BIAS+64+0*8], %i0
	ldx [%06+STACK_BIAS+7*8], %l7
	ldx [%06+STACK_BIAS+6*8], %l6
	ldx [%06+STACK_BIAS+5*8], %l5
	ldx [%06+STACK_BIAS+4*8], %l4
	ldx [%06+STACK_BIAS+3*8], %l3
	ldx [%06+STACK_BIAS+2*8], %l2
	ldx [%06+STACK_BIAS+1*8], %l1
	ldx [%06+STACK_BIAS+0*8], %l0
	ret


#ifdef __ELF__
.section .note.GNU-stack,"",%progbits
#endif

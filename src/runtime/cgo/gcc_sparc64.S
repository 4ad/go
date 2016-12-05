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

// Solaris wants 4-byte alignment for SPARC code sections.
.align	4

/*
 * void crosscall1(void (*fn)(void), void (*setg_gcc)(void *g), void *g)
 *
 * Calling into the gc tool chain, where all registers are caller save.
 * Called from standard SPARC ABI, where r16-r29 are callee-save, so they
 * must be saved explicitly.
 */
.globl EXT(crosscall1)
EXT(crosscall1):
	save	%sp, -304, %sp
	flushw

	mov	%i0, %l0
	mov	%i1, %l1
	mov	%i2, %o0

	call %l1
	nop
	call %l0
	nop

	ret
	restore
#ifdef __ELF__
.size	crosscall1,.-crosscall1
.type	crosscall1,#function
.section .note.GNU-stack,"",%progbits
#endif

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// void runtime·memmove(void*, void*, uintptr)
TEXT runtime·memmove(SB), NOSPLIT, $-8-24
	MOVD	to+0(FP), R3
	MOVD	from+8(FP), R4
	MOVD	n+16(FP), R5
	CMP	ZR, R5
	BNED	check
	RET

check:
	AND $7, R5, R6

	CMP	R3, R4
	BLD	backward

	// Copying forward proceeds by copying R10/8 words then copying R6 bytes.
	// R3 and R4 are advanced as we copy.

	CMP	ZR, R10		// Do we need to do any word-by-word copying?
	BED	noforwardlarge

	ADD	R3, R10, R9	// R9 points just past where we copy by word

forwardlargeloop:
	MOVD	(R4), R8	// R8 is just a scratch register
	ADD	$8, R4
	MOVD	R8, (R3)
	ADD	$8, R3
	CMP	R3, R9
	BNED	forwardlargeloop

noforwardlarge:
	CMP	ZR, R6		// Do we need to do any byte-by-byte copying?
	BNED	forwardtail
	RET

forwardtail:
	ADD	R3, R6, R9	// R9 points just past the destination memory

forwardtailloop:
	MOVUB (R4), R8
	ADD	$1, R4
	MOVUB	R8, (R3)
	ADD	$1, R3
	CMP	R3, R9
	BNED	forwardtailloop
	RET

backward:
	// Copying backwards proceeds by copying R6 bytes then copying R10/8 words.
	// R3 and R4 are advanced to the end of the destination/source buffers
	// respectively and moved back as we copy.

	ADD	R4, R5, R4	// R4 points just past the last source byte
	ADD	R3, R5, R3	// R3 points just past the last destination byte

	CMP	ZR, R6		// Do we need to do any byte-by-byte copying?
	BED	nobackwardtail

	SUB	R6, R3, R9	// R9 points at the lowest destination byte that should be copied by byte.
backwardtailloop:
	ADD	$-1, R4
	MOVUB	(R4), R8
	ADD	$-1, R3
	MOVUB	R8, -1(R3)
	CMP	R9, R3
	BNED	backwardtailloop

nobackwardtail:
	CMP     ZR, R10		// Do we need to do any word-by-word copying?
	BNED	backwardlarge
	RET

backwardlarge:
        SUB	R10, R3, R9      // R9 points at the lowest destination byte

backwardlargeloop:
	ADD	$-8, R4
	MOVD	(R4), R8
	ADD	$-8, R3
	MOVD	R8, (R3)
	CMP	R9, R3
	BNED	backwardlargeloop
	RET

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT runtime·memclr(SB),NOSPLIT|NOFRAME,$0-16
	MOVD	ptr+0(FP), R3
	MOVD	n+8(FP), R4
	AND $7, R4, R6

	CMP	ZR, R5
	BED	nowords

	ADD	R3, R5, R5

wordloop: // TODO: Optimize for unaligned ptr.
	MOVD	ZR, (R3)
	ADD	$8, R3
	CMP	R3, R5
	BNED	wordloop
nowords:
        CMP	$0, R6
        BED	done

	ADD	R3, R6, R6

byteloop:
	MOVUB	ZR, 1(R3)
	ADD	$1, R3
	CMP	R3, R6
	BNED	byteloop
done:
	RET

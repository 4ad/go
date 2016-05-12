// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT runtime·memclr(SB),NOSPLIT|NOFRAME,$0-16
	MOVD	ptr+0(FP), R3
	MOVD	n+8(FP), R4
	AND $7, R4, R9

	CMP	ZR, R22
	BED	nowords

	ADD	R3, R22, R22

wordloop: // TODO: Optimize for unaligned ptr.
	MOVD	ZR, (R3)
	ADD	$8, R3
	CMP	R3, R22
	BNED	wordloop
nowords:
        CMP	$0, R9
        BED	done

	ADD	R3, R9, R9

byteloop:
	MOVUB	ZR, (R3)
	ADD	$1, R3
	CMP	R3, R9
	BNED	byteloop
done:
	RET

// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// funcdata for functions with no local variables in frame.
// Define two zero-length bitmaps, because the same index is used
// for the local variables as for the argument frame, and assembly
// frames have two argument bitmaps, one without results and one with results.
DATA runtime·no_pointers_stackmap+0x00(SB)/4, $2
DATA runtime·no_pointers_stackmap+0x04(SB)/4, $0
GLOBL runtime·no_pointers_stackmap(SB),RODATA, $8

TEXT runtime·nop(SB),NOSPLIT,$0-0
	RET

//GLOBL runtime·mheap_(SB), NOPTR, $0
//GLOBL runtime·memstats(SB), NOPTR, $0

// Linker has a bug, and we need non-zero length symbols in
// these sections.

DATA type·runtime·moduledata(SB)/8, $224 // must match module size
GLOBL type·runtime·moduledata(SB), 0, $224

DATA data(SB)/4, $2
GLOBL data(SB), 0, $4

DATA rodata(SB)/4, $1
GLOBL rodata(SB), RODATA, $4

// .noptrdata
DATA noptrdata(SB)/4, $3
GLOBL noptrdata(SB), NOPTR, $4

// .bss
GLOBL bss(SB), 0, $4

// .noptrbss
GLOBL noptrbss(SB), NOPTR, $4

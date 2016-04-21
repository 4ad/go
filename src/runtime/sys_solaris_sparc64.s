// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for SPARC64, Solaris.
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"


// void libc_miniterrno(void *(*___errno)(void));
//
// Set the TLS errno pointer in M.
//
// Called using runtime·asmcgocall from os_solaris.c:/minit.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·miniterrno(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$70, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// int64 runtime·nanotime1(void);
//
// clock_gettime(3c) wrapper because Timespec is too large for
// runtime·nanotime stack.
//
// Called using runtime·sysvicall6 from os_solaris.c:/nanotime.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·nanotime1(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$71, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// pipe(3c) wrapper that returns fds in AX, DX.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·pipe1(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$72, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// Call a library function with SysV calling conventions.
// The called function can take a maximum of 6 INTEGER class arguments,
// see 
//   Michael Matz, Jan Hubicka, Andreas Jaeger, and Mark Mitchell
//   System V Application Binary Interface 
//   AMD64 Architecture Processor Supplement
// section 3.2.3.
//
// Called by runtime·asmcgocall or runtime·cgocall.
// NOT USING GO CALLING CONVENTION.
TEXT runtime·asmsysvicall6(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$73, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// uint32 tstart_sysvicall(M *newm);
TEXT runtime·tstart_sysvicall(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$74, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// Careful, this is called by __sighndlr, a libc function. We must preserve
// registers as per AMD 64 ABI.
TEXT runtime·sigtramp(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$75, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// Called from runtime·usleep (Go). Can be called on Go stack, on OS stack,
// can also be called in cgo callback path without a g->m.
TEXT runtime·usleep1(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$76, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// Runs on OS stack. duration (in µs units) is in DI.
TEXT runtime·usleep2(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$77, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// Runs on OS stack, called from runtime·osyield.
TEXT runtime·osyield1(SB),NOSPLIT,$0
	// TODO(aram):
	MOVD	$78, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

// func now() (sec int64, nsec int32)
TEXT time·now(SB),NOSPLIT,$8-12
	// TODO(aram):
	MOVD	$79, R1
	ADD	$'!', R1, R1
	MOVD	R1, dbgbuf(SB)
	MOVD	$2, R8
	MOVD	$dbgbuf(SB), R9
	MOVD	$2, R10
	MOVD	$libc_exit(SB), R1
	CALL	R1
	UNDEF
	RET

TEXT	_rt0_sparc64_solaris(SB),7,$0-0
	OR	ZR, ZR, R24
	OR	ZR, ZR, R25
	OR	ZR, ZR, TMP
	OR	ZR, ZR, R27
	OR	ZR, ZR, R28
	OR	ZR, ZR, CTXT
	CALL	main(SB)
	// sys_exit(1)
	MOVD	$1, R8
	MOVD	$1, TMP
	TA	$0x40
	RET

TEXT	hello(SB),7,$0-0
	MOVD	$1, R8
	MOVD	$hellomsg(SB), R9
	MOVD	$7, R10
	MOVD	$libc_write(SB), R1
	CALL	R1
	RET

DATA hellomsg(SB)/8, $"hello!\n"
GLOBL hellomsg(SB), 16, $8

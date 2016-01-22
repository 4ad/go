TEXT	_rt0_sparc64_solaris(SB),7,$0-0
	OR	ZR, ZR, R24
	OR	ZR, ZR, R25
	OR	ZR, ZR, TMP
	OR	ZR, ZR, R27
	OR	ZR, ZR, R28
	OR	ZR, ZR, R29
	CALL	main(SB)
	// sys_exit(1)
	MOVD	$1, R8
	MOVD	$1, TMP
	TA	$0x40
	RET

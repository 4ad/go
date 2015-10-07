TEXT	_rt0_sparc64_solaris(SB),7,$0-0
	MOVD	$0, R24
	MOVD	$0, R25
	MOVD	$0, R26
	MOVD	$0, R27
	MOVD	$0, R28
	MOVD	$0, R29
	CALL	main(SB)
	// sys_exit(1)
	MOVD	$1, R8
	MOVD	$1, TMP
	TA	$0x40
	RET

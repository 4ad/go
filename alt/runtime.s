TEXT	_rt0_sparc64_solaris(SB),7,$-8
	CALL	main(SB)
	// sys_exit(1)
	MOVD	$1, R8
	MOVD	$1, TMP
	TA	$0x40
	RET

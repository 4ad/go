TEXT	main(SB),7,$-8
	MOVD	$1, R8
	MOVD	$msg(SB), R9
	MOVD	$6, R10
	MOVD	$4, TMP	// SYS_WRITE
	TA	$0x40
	RET

DATA msg(SB)/8, $"hello\n"
GLOBL msg(SB), SNOPTRDATA, $8

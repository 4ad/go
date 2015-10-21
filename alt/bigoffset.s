TEXT	main(SB),7,$-8
	MOVD	$64(R2), R3
	MOVD	$1, R8
	MOVD	$msg(SB), R9
	MOVD	$4, R10
	MOVD	$4, TMP	// SYS_WRITE
	TA	$0x40
	RET

DATA msg(SB)/8, $"big\n"
GLOBL msg(SB), SNOPTRDATA, $8

TEXT	main(SB),512|7,$0
	MOVD	$64(L2), L3
	MOVD	$1, R8
	MOVD	$msg(SB), R9
	MOVD	$4, R10
	MOVD	$4, TMP	// SYS_WRITE
	TA	$0x40
	RET

DATA msg(SB)/8, $"big\n"
GLOBL msg(SB), 16, $8

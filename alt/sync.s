TEXT	main(SB),7,$-8
	JMP	cas(SB)
	RET

TEXT cas(SB), 7, $0-17
	MOVD	ptr+0(FP), R3
	MOVUW	old+8(FP), R1
	MOVUW	new+12(FP), R2
	MEMBAR	$15
	CASW	(R3), R1, R2
	XOR	R1, R2, R1
	SUBCC	R1, ZR, ZR
	SUBC	$-1, ZR, R1
	MEMBAR	$15
	AND	$0xff, R1, R1
	SRAW	$0, R1, R1
	MOVB	R1, ret+16(FP)
	RET

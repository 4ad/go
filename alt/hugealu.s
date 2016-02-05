TEXT	main(SB),7,$-8
	ADD	$0xf00abcd, R1
	AND	$0xf00abcd, R1, R2
	ADD	$0x1f00abcd, R3
	MOVD	R4, 0xf00ddd(R5)
	MOVD	R6, -16521(R7)
	RET

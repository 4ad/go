TEXT	main(SB),7,$0-0
	CALL	foo0(SB)
	RET

TEXT	foo0(SB),7,$0-0
	CALL	foo1(SB)
	RET

TEXT	foo1(SB),7,$0-0
	CALL	foo2(SB)
	RET

TEXT	foo2(SB),7,$0-0
	CALL	foo3(SB)
	RET

TEXT	foo3(SB),7,$0-0
	CALL	foo4(SB)
	RET

TEXT	foo4(SB),7,$0-0
	CALL	foo5(SB)
	RET

TEXT	foo5(SB),7,$0-0
	// sys_write(1, "", _)
	MOVD	$1, R8
	MOVD	$msg(SB), R9
	MOVD	$6, R10
	MOVD	$4, TMP
	TA	$0x40
	RET

DATA msg(SB)/8, $"hello\n"
GLOBL msg(SB), SNOPTRDATA, $8

TEXT	main(SB),7,$32-0
	MOVD	$1, fd-24(SP)
	MOVD	$msg(SB), RT1
	MOVD	RT1, buf-16(SP)
	CALL	foo0(SB)
	RET

TEXT	foo0(SB),7,$32-16
	MOVD	fd+0(FP), L1
	MOVD	buf+8(FP), L2

	MOVD	L1, fd-32(SP)
	MOVD	L2, buf-24(SP)
	MOVD	$6, len-16(SP)
	MOVD	$4, trap-8(SP)
	CALL	foo1(SB)
	RET

TEXT	foo1(SB),7,$0-32
	MOVD	fd+0(FP), R8
	MOVD	buf+8(FP), R9
	MOVD	len+16(FP), R10
	MOVD	trap+24(FP), TMP
	TA	$0x40
	RET

DATA msg(SB)/8, $"hello\n"
GLOBL msg(SB), 16, $8

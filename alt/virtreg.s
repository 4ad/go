TEXT	main(SB),7,$32-0
	MOVD	$1, fd-32(SP)
	MOVD	$msg(SB), buf-24(SP)
	MOVD	$6, len-16(SP)
	MOVD	$4, trap-8(SP)
	CALL	foo0(SB)
	RET

TEXT	foo0(SB),7,$0-16
	MOVD	fd+0(FP), R8
	MOVD	buf+8(FP), R9
	MOVD	len+16(FP), R10
	MOVD	trap+24(FP), TMP
	TA	$0x40
	RET

DATA msg(SB)/8, $"hello\n"
GLOBL msg(SB), SNOPTRDATA, $8

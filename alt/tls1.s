#define TLSBSS	256

TEXT	main(SB),512|7,$0
	MOVD	$runtime·tls_g+0(SB), R1
	RET

GLOBL runtime·tls_g+0(SB), TLSBSS, $8

TEXT	main(SB),7,$-8
	CMP	R2, R1
	SUBCC	R2, R1, ZR
	BLE	ICC, label
	MOVD	$1, R1
	RET
label:
	MOVD	$2, R1
	RET

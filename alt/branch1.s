TEXT	main(SB),7,$-8
	BNE	XCC, l1
	BNED	l1
	BLE	ICC, l2
	BLEW	l2
	MOVD	$1, R1
l1:
l2:
	RET

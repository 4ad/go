#!/usr/bin/awk -f

! /DIE/ {
	printf("%s\n", $0)
}

/DIE/ {
	code++
	printf("	// TODO(aram):\n");
	printf("	MOVD	$%d, TMP\n", code);
	printf("	ADD	$'!', TMP, TMP\n");
	printf("	MOVD	TMP, dbgbuf(SB)\n");
	printf("	MOVD	$2, R8\n");
	printf("	MOVD	$dbgbuf(SB), R9\n");
	printf("	MOVD	$2, R10\n");
	printf("	MOVD	$libc_exit(SB), TMP\n");
	printf("	CALL	TMP\n")
	printf("	UNDEF\n")
}

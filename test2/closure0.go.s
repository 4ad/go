"".newfunc t=1 size=64 value=0 args=0x8 locals=0x0 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:13)	TEXT	"".newfunc(SB), $0-8
	0x0000 00000 (/Users/aram/go/test2/closure0.go:13)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:13)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:13)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:13)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:13)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	FUNCDATA	$0, gclocals·0fb5f740dc3899c17d2f00dd94c805d6(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	MOVD	$0, R8
	0x0018 00024 (/Users/aram/go/test2/closure0.go:13)	MOVD	$"".newfunc.func1·f(SB), R8
	0x0030 00048 (/Users/aram/go/test2/closure0.go:13)	MOVD	R8, "".~r0(FP)
	0x0034 00052 (/Users/aram/go/test2/closure0.go:13)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 90 10 20 00 0b 00 00 00 8a 11 60 00  .s.w.. .......`.
	0x0020 8b 29 70 20 11 00 00 00 90 12 20 00 90 11 40 08  .)p ...... ...@.
	0x0030 d0 77 a8 af 81 c7 e0 08 81 e8 20 00 00 00 00 00  .w........ .....
	rel 24+8 t=6 "".newfunc.func1·f+0
	rel 36+8 t=5 "".newfunc.func1·f+0
"".newfunc2 t=1 size=128 value=0 args=0x10 locals=0x10
	0x0000 00000 (/Users/aram/go/test2/closure0.go:14)	TEXT	"".newfunc2(SB), $16-16
	0x0000 00000 (/Users/aram/go/test2/closure0.go:14)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:14)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:14)	MOVD	$-192, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:14)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:14)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	FUNCDATA	$0, gclocals·aecfa9ecce04c513ee6b217848214030(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	MOVD	$0, R8
	0x0018 00024 (/Users/aram/go/test2/closure0.go:14)	MOVD	$type.struct { F uintptr; x int }(SB), R8
	0x0030 00048 (/Users/aram/go/test2/closure0.go:14)	MOVD	R8, 0+176+2047(RSP)
	0x0034 00052 (/Users/aram/go/test2/closure0.go:14)	PCDATA	$0, $0
	0x0034 00052 (/Users/aram/go/test2/closure0.go:14)	CALL	runtime.newobject(SB)
	0x0038 00056 (/Users/aram/go/test2/closure0.go:14)	RNOP
	0x003c 00060 (/Users/aram/go/test2/closure0.go:14)	MOVD	8+176+2047(RSP), R8
	0x0040 00064 (/Users/aram/go/test2/closure0.go:14)	MOVD	R8, R10
	0x0044 00068 (/Users/aram/go/test2/closure0.go:14)	MOVD	R10, R8
	0x0048 00072 (/Users/aram/go/test2/closure0.go:14)	MOVD	$"".newfunc2.func1(SB), R9
	0x0060 00096 (/Users/aram/go/test2/closure0.go:14)	MOVD	R9, (R8)
	0x0064 00100 (/Users/aram/go/test2/closure0.go:14)	MOVD	R10, R8
	0x0068 00104 (/Users/aram/go/test2/closure0.go:14)	MOVD	"".x(FP), R9
	0x006c 00108 (/Users/aram/go/test2/closure0.go:14)	MOVD	R9, 8(R8)
	0x0070 00112 (/Users/aram/go/test2/closure0.go:14)	MOVD	R10, R8
	0x0074 00116 (/Users/aram/go/test2/closure0.go:14)	MOVD	R8, "".~r1+8(FP)
	0x0078 00120 (/Users/aram/go/test2/closure0.go:14)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 40 9d e3 80 01  ..........?@....
	0x0010 fe 73 a8 77 90 10 20 00 0b 00 00 00 8a 11 60 00  .s.w.. .......`.
	0x0020 8b 29 70 20 11 00 00 00 90 12 20 00 90 11 40 08  .)p ...... ...@.
	0x0030 d0 73 a8 af 40 00 00 00 01 00 00 00 d0 5b a8 b7  .s..@........[..
	0x0040 94 10 00 08 90 10 00 0a 0b 00 00 00 8a 11 60 00  ..............`.
	0x0050 8b 29 70 20 13 00 00 00 92 12 60 00 92 11 40 09  .)p ......`...@.
	0x0060 d2 72 20 00 90 10 00 0a d2 5f a8 af d2 72 20 08  .r ......_...r .
	0x0070 90 10 00 0a d0 77 a8 b7 81 c7 e0 08 81 e8 20 00  .....w........ .
	rel 24+8 t=6 type.struct { F uintptr; x int }+0
	rel 36+8 t=5 type.struct { F uintptr; x int }+0
	rel 52+4 t=14 runtime.newobject+0
	rel 72+8 t=6 "".newfunc2.func1+0
	rel 84+8 t=5 "".newfunc2.func1+0
"".main t=1 size=784 value=0 args=0x0 locals=0x50
	0x0000 00000 (/Users/aram/go/test2/closure0.go:16)	TEXT	"".main(SB), $80-0
	0x0000 00000 (/Users/aram/go/test2/closure0.go:16)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:16)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:16)	MOVD	$-256, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:16)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:16)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:16)	FUNCDATA	$0, gclocals·4329624ce4271de83fc7c43fc9c7e126(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:16)	FUNCDATA	$1, gclocals·68bfa0232bc220e249cf20baf5784e82(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:17)	PCDATA	$0, $0
	0x0014 00020 (/Users/aram/go/test2/closure0.go:17)	CALL	"".newfunc(SB)
	0x0018 00024 (/Users/aram/go/test2/closure0.go:17)	RNOP
	0x001c 00028 (/Users/aram/go/test2/closure0.go:17)	MOVD	0+176+2047(RSP), R8
	0x0020 00032 (/Users/aram/go/test2/closure0.go:17)	MOVD	R8, "".autotmp_0001-24(SP)
	0x0024 00036 (/Users/aram/go/test2/closure0.go:17)	PCDATA	$0, $1
	0x0024 00036 (/Users/aram/go/test2/closure0.go:17)	CALL	"".newfunc(SB)
	0x0028 00040 (/Users/aram/go/test2/closure0.go:17)	RNOP
	0x002c 00044 (/Users/aram/go/test2/closure0.go:17)	MOVD	0+176+2047(RSP), R8
	0x0030 00048 (/Users/aram/go/test2/closure0.go:17)	MOVD	R8, R10
	0x0034 00052 (/Users/aram/go/test2/closure0.go:17)	MOVD	"".autotmp_0001-24(SP), R8
	0x0038 00056 (/Users/aram/go/test2/closure0.go:17)	MOVD	R8, R9
	0x003c 00060 (/Users/aram/go/test2/closure0.go:17)	MOVD	R10, R8
	0x0040 00064 (/Users/aram/go/test2/closure0.go:17)	MOVD	R8, "".y-32(SP)
	0x0044 00068 (/Users/aram/go/test2/closure0.go:18)	MOVD	$1, R8
	0x0048 00072 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, 0+176+2047(RSP)
	0x004c 00076 (/Users/aram/go/test2/closure0.go:18)	MOVD	R9, R8
	0x0050 00080 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, CTXT
	0x0054 00084 (/Users/aram/go/test2/closure0.go:18)	MOVD	(CTXT), RT1
	0x0058 00088 (/Users/aram/go/test2/closure0.go:18)	PCDATA	$0, $2
	0x0058 00088 (/Users/aram/go/test2/closure0.go:18)	CALL	CTXT, RT1
	0x005c 00092 (/Users/aram/go/test2/closure0.go:18)	RNOP
	0x0060 00096 (/Users/aram/go/test2/closure0.go:18)	MOVD	8+176+2047(RSP), R8
	0x0064 00100 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, R9
	0x0068 00104 (/Users/aram/go/test2/closure0.go:18)	MOVD	R9, R8
	0x006c 00108 (/Users/aram/go/test2/closure0.go:18)	CMP	$1, R8
	0x0070 00112 (/Users/aram/go/test2/closure0.go:18)	BNED	664
	0x0074 00116 (/Users/aram/go/test2/closure0.go:18)	RNOP
	0x0078 00120 (/Users/aram/go/test2/closure0.go:18)	MOVD	$2, R8
	0x007c 00124 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, 0+176+2047(RSP)
	0x0080 00128 (/Users/aram/go/test2/closure0.go:18)	MOVD	"".y-32(SP), R8
	0x0084 00132 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, CTXT
	0x0088 00136 (/Users/aram/go/test2/closure0.go:18)	MOVD	(CTXT), RT1
	0x008c 00140 (/Users/aram/go/test2/closure0.go:18)	PCDATA	$0, $0
	0x008c 00140 (/Users/aram/go/test2/closure0.go:18)	CALL	CTXT, RT1
	0x0090 00144 (/Users/aram/go/test2/closure0.go:18)	RNOP
	0x0094 00148 (/Users/aram/go/test2/closure0.go:18)	MOVD	8+176+2047(RSP), R8
	0x0098 00152 (/Users/aram/go/test2/closure0.go:18)	MOVD	R8, R9
	0x009c 00156 (/Users/aram/go/test2/closure0.go:18)	MOVD	R9, R8
	0x00a0 00160 (/Users/aram/go/test2/closure0.go:18)	CMP	$2, R8
	0x00a4 00164 (/Users/aram/go/test2/closure0.go:18)	BNED	664
	0x00a8 00168 (/Users/aram/go/test2/closure0.go:18)	RNOP
	0x00ac 00172 (/Users/aram/go/test2/closure0.go:22)	MOVD	$2, R8
	0x00b0 00176 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, 0+176+2047(RSP)
	0x00b4 00180 (/Users/aram/go/test2/closure0.go:22)	PCDATA	$0, $0
	0x00b4 00180 (/Users/aram/go/test2/closure0.go:22)	CALL	"".newfunc2(SB)
	0x00b8 00184 (/Users/aram/go/test2/closure0.go:22)	RNOP
	0x00bc 00188 (/Users/aram/go/test2/closure0.go:22)	MOVD	8+176+2047(RSP), R8
	0x00c0 00192 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, "".autotmp_0005-24(SP)
	0x00c4 00196 (/Users/aram/go/test2/closure0.go:22)	MOVD	$1, R8
	0x00c8 00200 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, 0+176+2047(RSP)
	0x00cc 00204 (/Users/aram/go/test2/closure0.go:22)	PCDATA	$0, $1
	0x00cc 00204 (/Users/aram/go/test2/closure0.go:22)	CALL	"".newfunc2(SB)
	0x00d0 00208 (/Users/aram/go/test2/closure0.go:22)	RNOP
	0x00d4 00212 (/Users/aram/go/test2/closure0.go:22)	MOVD	8+176+2047(RSP), R8
	0x00d8 00216 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, R10
	0x00dc 00220 (/Users/aram/go/test2/closure0.go:22)	MOVD	"".autotmp_0005-24(SP), R8
	0x00e0 00224 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, R9
	0x00e4 00228 (/Users/aram/go/test2/closure0.go:22)	MOVD	R10, R8
	0x00e8 00232 (/Users/aram/go/test2/closure0.go:22)	MOVD	R8, "".y-32(SP)
	0x00ec 00236 (/Users/aram/go/test2/closure0.go:23)	MOVD	$1, R8
	0x00f0 00240 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, 0+176+2047(RSP)
	0x00f4 00244 (/Users/aram/go/test2/closure0.go:23)	MOVD	R9, R8
	0x00f8 00248 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, CTXT
	0x00fc 00252 (/Users/aram/go/test2/closure0.go:23)	MOVD	(CTXT), RT1
	0x0100 00256 (/Users/aram/go/test2/closure0.go:23)	PCDATA	$0, $2
	0x0100 00256 (/Users/aram/go/test2/closure0.go:23)	CALL	CTXT, RT1
	0x0104 00260 (/Users/aram/go/test2/closure0.go:23)	RNOP
	0x0108 00264 (/Users/aram/go/test2/closure0.go:23)	MOVD	8+176+2047(RSP), R8
	0x010c 00268 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, R9
	0x0110 00272 (/Users/aram/go/test2/closure0.go:23)	MOVD	R9, R8
	0x0114 00276 (/Users/aram/go/test2/closure0.go:23)	CMP	$2, R8
	0x0118 00280 (/Users/aram/go/test2/closure0.go:23)	BNED	552
	0x011c 00284 (/Users/aram/go/test2/closure0.go:23)	RNOP
	0x0120 00288 (/Users/aram/go/test2/closure0.go:23)	MOVD	$2, R8
	0x0124 00292 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, 0+176+2047(RSP)
	0x0128 00296 (/Users/aram/go/test2/closure0.go:23)	MOVD	"".y-32(SP), R8
	0x012c 00300 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, CTXT
	0x0130 00304 (/Users/aram/go/test2/closure0.go:23)	MOVD	(CTXT), RT1
	0x0134 00308 (/Users/aram/go/test2/closure0.go:23)	PCDATA	$0, $0
	0x0134 00308 (/Users/aram/go/test2/closure0.go:23)	CALL	CTXT, RT1
	0x0138 00312 (/Users/aram/go/test2/closure0.go:23)	RNOP
	0x013c 00316 (/Users/aram/go/test2/closure0.go:23)	MOVD	8+176+2047(RSP), R8
	0x0140 00320 (/Users/aram/go/test2/closure0.go:23)	MOVD	R8, R9
	0x0144 00324 (/Users/aram/go/test2/closure0.go:23)	MOVD	R9, R8
	0x0148 00328 (/Users/aram/go/test2/closure0.go:23)	CMP	$1, R8
	0x014c 00332 (/Users/aram/go/test2/closure0.go:23)	BNED	552
	0x0150 00336 (/Users/aram/go/test2/closure0.go:23)	RNOP
	0x0154 00340 (/Users/aram/go/test2/closure0.go:28)	MOVUB	"".fail(SB), R8
	0x0170 00368 (/Users/aram/go/test2/closure0.go:28)	CMP	$0, R8
	0x0174 00372 (/Users/aram/go/test2/closure0.go:28)	BEW	544
	0x0178 00376 (/Users/aram/go/test2/closure0.go:28)	RNOP
	0x017c 00380 (/Users/aram/go/test2/closure0.go:29)	MOVD	$go.string."fail"(SB), R8
	0x0194 00404 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, "".autotmp_0009-16(SP)
	0x0198 00408 (/Users/aram/go/test2/closure0.go:29)	MOVD	$4, R8
	0x019c 00412 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, "".autotmp_0009-8(SP)
	0x01a0 00416 (/Users/aram/go/test2/closure0.go:29)	MOVD	$type.string(SB), R8
	0x01b8 00440 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, 0+176+2047(RSP)
	0x01bc 00444 (/Users/aram/go/test2/closure0.go:29)	MOVD	$"".autotmp_0009-16(SP), R8
	0x01c0 00448 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, 8+176+2047(RSP)
	0x01c4 00452 (/Users/aram/go/test2/closure0.go:29)	MOVD	$0, R8
	0x01c8 00456 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, 16+176+2047(RSP)
	0x01cc 00460 (/Users/aram/go/test2/closure0.go:29)	PCDATA	$0, $3
	0x01cc 00460 (/Users/aram/go/test2/closure0.go:29)	CALL	runtime.convT2E(SB)
	0x01d0 00464 (/Users/aram/go/test2/closure0.go:29)	RNOP
	0x01d4 00468 (/Users/aram/go/test2/closure0.go:29)	MOVD	$24+176+2047(RSP), R8
	0x01d8 00472 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, R8
	0x01dc 00476 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, R9
	0x01e0 00480 (/Users/aram/go/test2/closure0.go:29)	MOVD	$0+176+2047(RSP), R8
	0x01e4 00484 (/Users/aram/go/test2/closure0.go:29)	MOVD	R8, R8
	0x01e8 00488 (/Users/aram/go/test2/closure0.go:29)	ADD	$-8, R9
	0x01ec 00492 (/Users/aram/go/test2/closure0.go:29)	ADD	$-8, R8
	0x01f0 00496 (/Users/aram/go/test2/closure0.go:29)	MOVD	$8, R11
	0x01f4 00500 (/Users/aram/go/test2/closure0.go:29)	MOVD	(R9)(R11*1), R10
	0x01f8 00504 (/Users/aram/go/test2/closure0.go:29)	ADD	R11, R9
	0x01fc 00508 (/Users/aram/go/test2/closure0.go:29)	MOVD	R10, (R8)(R11*1)
	0x0200 00512 (/Users/aram/go/test2/closure0.go:29)	ADD	R11, R8
	0x0204 00516 (/Users/aram/go/test2/closure0.go:29)	MOVD	(R9)(R11*1), R10
	0x0208 00520 (/Users/aram/go/test2/closure0.go:29)	ADD	R11, R9
	0x020c 00524 (/Users/aram/go/test2/closure0.go:29)	MOVD	R10, (R8)(R11*1)
	0x0210 00528 (/Users/aram/go/test2/closure0.go:29)	ADD	R11, R8
	0x0214 00532 (/Users/aram/go/test2/closure0.go:29)	PCDATA	$0, $3
	0x0214 00532 (/Users/aram/go/test2/closure0.go:29)	CALL	runtime.gopanic(SB)
	0x0218 00536 (/Users/aram/go/test2/closure0.go:29)	RNOP
	0x021c 00540 (/Users/aram/go/test2/closure0.go:29)	UNDEF
	0x0220 00544 (/Users/aram/go/test2/closure0.go:31)	RETRESTORE
	0x0228 00552 (/Users/aram/go/test2/closure0.go:24)	PCDATA	$0, $0
	0x0228 00552 (/Users/aram/go/test2/closure0.go:24)	CALL	runtime.printlock(SB)
	0x022c 00556 (/Users/aram/go/test2/closure0.go:24)	RNOP
	0x0230 00560 (/Users/aram/go/test2/closure0.go:24)	MOVD	$go.string."newfunc2 returned broken funcs"(SB), R8
	0x0248 00584 (/Users/aram/go/test2/closure0.go:24)	MOVD	R8, 0+176+2047(RSP)
	0x024c 00588 (/Users/aram/go/test2/closure0.go:24)	MOVD	$30, R8
	0x0250 00592 (/Users/aram/go/test2/closure0.go:24)	MOVD	R8, 8+176+2047(RSP)
	0x0254 00596 (/Users/aram/go/test2/closure0.go:24)	PCDATA	$0, $0
	0x0254 00596 (/Users/aram/go/test2/closure0.go:24)	CALL	runtime.printstring(SB)
	0x0258 00600 (/Users/aram/go/test2/closure0.go:24)	RNOP
	0x025c 00604 (/Users/aram/go/test2/closure0.go:24)	PCDATA	$0, $0
	0x025c 00604 (/Users/aram/go/test2/closure0.go:24)	CALL	runtime.printnl(SB)
	0x0260 00608 (/Users/aram/go/test2/closure0.go:24)	RNOP
	0x0264 00612 (/Users/aram/go/test2/closure0.go:24)	PCDATA	$0, $0
	0x0264 00612 (/Users/aram/go/test2/closure0.go:24)	CALL	runtime.printunlock(SB)
	0x0268 00616 (/Users/aram/go/test2/closure0.go:24)	RNOP
	0x026c 00620 (/Users/aram/go/test2/closure0.go:25)	MOVD	$1, R8
	0x0270 00624 (/Users/aram/go/test2/closure0.go:25)	MOVUB	R8, R8
	0x0274 00628 (/Users/aram/go/test2/closure0.go:25)	MOVUB	R8, "".fail(SB)
	0x0290 00656 (/Users/aram/go/test2/closure0.go:28)	JMP	340
	0x0294 00660 (/Users/aram/go/test2/closure0.go:28)	RNOP
	0x0298 00664 (/Users/aram/go/test2/closure0.go:19)	PCDATA	$0, $0
	0x0298 00664 (/Users/aram/go/test2/closure0.go:19)	CALL	runtime.printlock(SB)
	0x029c 00668 (/Users/aram/go/test2/closure0.go:19)	RNOP
	0x02a0 00672 (/Users/aram/go/test2/closure0.go:19)	MOVD	$go.string."newfunc returned broken funcs"(SB), R8
	0x02b8 00696 (/Users/aram/go/test2/closure0.go:19)	MOVD	R8, 0+176+2047(RSP)
	0x02bc 00700 (/Users/aram/go/test2/closure0.go:19)	MOVD	$29, R8
	0x02c0 00704 (/Users/aram/go/test2/closure0.go:19)	MOVD	R8, 8+176+2047(RSP)
	0x02c4 00708 (/Users/aram/go/test2/closure0.go:19)	PCDATA	$0, $0
	0x02c4 00708 (/Users/aram/go/test2/closure0.go:19)	CALL	runtime.printstring(SB)
	0x02c8 00712 (/Users/aram/go/test2/closure0.go:19)	RNOP
	0x02cc 00716 (/Users/aram/go/test2/closure0.go:19)	PCDATA	$0, $0
	0x02cc 00716 (/Users/aram/go/test2/closure0.go:19)	CALL	runtime.printnl(SB)
	0x02d0 00720 (/Users/aram/go/test2/closure0.go:19)	RNOP
	0x02d4 00724 (/Users/aram/go/test2/closure0.go:19)	PCDATA	$0, $0
	0x02d4 00724 (/Users/aram/go/test2/closure0.go:19)	CALL	runtime.printunlock(SB)
	0x02d8 00728 (/Users/aram/go/test2/closure0.go:19)	RNOP
	0x02dc 00732 (/Users/aram/go/test2/closure0.go:20)	MOVD	$1, R8
	0x02e0 00736 (/Users/aram/go/test2/closure0.go:20)	MOVUB	R8, R8
	0x02e4 00740 (/Users/aram/go/test2/closure0.go:20)	MOVUB	R8, "".fail(SB)
	0x0300 00768 (/Users/aram/go/test2/closure0.go:22)	JMP	172
	0x0304 00772 (/Users/aram/go/test2/closure0.go:22)	RNOP
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 00 9d e3 80 01  ..........?.....
	0x0010 fe 73 a8 77 40 00 00 00 01 00 00 00 d0 5b a8 af  .s.w@........[..
	0x0020 d0 77 a7 e7 40 00 00 00 01 00 00 00 d0 5b a8 af  .w..@........[..
	0x0030 94 10 00 08 d0 5f a7 e7 92 10 00 08 90 10 00 0a  ....._..........
	0x0040 d0 77 a7 df 90 10 20 01 d0 73 a8 af 90 10 00 09  .w.... ..s......
	0x0050 84 10 00 08 c2 58 a0 00 9f c0 60 00 01 00 00 00  .....X....`.....
	0x0060 d0 5b a8 b7 92 10 00 08 90 10 00 09 80 a2 20 01  .[............ .
	0x0070 12 68 00 8a 01 00 00 00 90 10 20 02 d0 73 a8 af  .h........ ..s..
	0x0080 d0 5f a7 df 84 10 00 08 c2 58 a0 00 9f c0 60 00  ._.......X....`.
	0x0090 01 00 00 00 d0 5b a8 b7 92 10 00 08 90 10 00 09  .....[..........
	0x00a0 80 a2 20 02 12 68 00 7d 01 00 00 00 90 10 20 02  .. ..h.}...... .
	0x00b0 d0 73 a8 af 40 00 00 00 01 00 00 00 d0 5b a8 b7  .s..@........[..
	0x00c0 d0 77 a7 e7 90 10 20 01 d0 73 a8 af 40 00 00 00  .w.... ..s..@...
	0x00d0 01 00 00 00 d0 5b a8 b7 94 10 00 08 d0 5f a7 e7  .....[......._..
	0x00e0 92 10 00 08 90 10 00 0a d0 77 a7 df 90 10 20 01  .........w.... .
	0x00f0 d0 73 a8 af 90 10 00 09 84 10 00 08 c2 58 a0 00  .s...........X..
	0x0100 9f c0 60 00 01 00 00 00 d0 5b a8 b7 92 10 00 08  ..`......[......
	0x0110 90 10 00 09 80 a2 20 02 12 68 00 44 01 00 00 00  ...... ..h.D....
	0x0120 90 10 20 02 d0 73 a8 af d0 5f a7 df 84 10 00 08  .. ..s..._......
	0x0130 c2 58 a0 00 9f c0 60 00 01 00 00 00 d0 5b a8 b7  .X....`......[..
	0x0140 92 10 00 08 90 10 00 09 80 a2 20 01 12 68 00 37  .......... ..h.7
	0x0150 01 00 00 00 0b 00 00 00 8a 11 60 00 8b 29 70 20  ..........`..)p 
	0x0160 11 00 00 00 90 12 20 00 90 11 40 08 d0 0a 20 00  ...... ...@... .
	0x0170 80 a2 20 00 02 48 00 2b 01 00 00 00 0b 00 00 00  .. ..H.+........
	0x0180 8a 11 60 00 8b 29 70 20 11 00 00 00 90 12 20 00  ..`..)p ...... .
	0x0190 90 11 40 08 d0 77 a7 ef 90 10 20 04 d0 77 a7 f7  ..@..w.... ..w..
	0x01a0 0b 00 00 00 8a 11 60 00 8b 29 70 20 11 00 00 00  ......`..)p ....
	0x01b0 90 12 20 00 90 11 40 08 d0 73 a8 af 90 07 a7 ef  .. ...@..s......
	0x01c0 d0 73 a8 b7 90 10 20 00 d0 73 a8 bf 40 00 00 00  .s.... ..s..@...
	0x01d0 01 00 00 00 90 03 a8 c7 90 10 00 08 92 10 00 08  ................
	0x01e0 90 03 a8 af 90 10 00 08 92 02 7f f8 90 02 3f f8  ..............?.
	0x01f0 96 10 20 08 d4 5a 40 0b 92 02 40 0b d4 72 00 0b  .. ..Z@...@..r..
	0x0200 90 02 00 0b d4 5a 40 0b 92 02 40 0b d4 72 00 0b  .....Z@...@..r..
	0x0210 90 02 00 0b 40 00 00 00 01 00 00 00 00 0d ea d0  ....@...........
	0x0220 81 c7 e0 08 81 e8 20 00 40 00 00 00 01 00 00 00  ...... .@.......
	0x0230 0b 00 00 00 8a 11 60 00 8b 29 70 20 11 00 00 00  ......`..)p ....
	0x0240 90 12 20 00 90 11 40 08 d0 73 a8 af 90 10 20 1e  .. ...@..s.... .
	0x0250 d0 73 a8 b7 40 00 00 00 01 00 00 00 40 00 00 00  .s..@.......@...
	0x0260 01 00 00 00 40 00 00 00 01 00 00 00 90 10 20 01  ....@......... .
	0x0270 90 0a 20 ff 0b 00 00 00 8a 11 60 00 8b 29 70 20  .. .......`..)p 
	0x0280 21 00 00 00 a0 14 20 00 a0 11 40 10 d0 2c 20 00  !..... ...@.., .
	0x0290 10 4f ff b1 01 00 00 00 40 00 00 00 01 00 00 00  .O......@.......
	0x02a0 0b 00 00 00 8a 11 60 00 8b 29 70 20 11 00 00 00  ......`..)p ....
	0x02b0 90 12 20 00 90 11 40 08 d0 73 a8 af 90 10 20 1d  .. ...@..s.... .
	0x02c0 d0 73 a8 b7 40 00 00 00 01 00 00 00 40 00 00 00  .s..@.......@...
	0x02d0 01 00 00 00 40 00 00 00 01 00 00 00 90 10 20 01  ....@......... .
	0x02e0 90 0a 20 ff 0b 00 00 00 8a 11 60 00 8b 29 70 20  .. .......`..)p 
	0x02f0 21 00 00 00 a0 14 20 00 a0 11 40 10 d0 2c 20 00  !..... ...@.., .
	0x0300 10 4f ff 6b 01 00 00 00 00 00 00 00 00 00 00 00  .O.k............
	rel 20+4 t=14 "".newfunc+0
	rel 36+4 t=14 "".newfunc+0
	rel 180+4 t=14 "".newfunc2+0
	rel 204+4 t=14 "".newfunc2+0
	rel 340+8 t=6 "".fail+0
	rel 352+8 t=5 "".fail+0
	rel 380+8 t=6 go.string."fail"+0
	rel 392+8 t=5 go.string."fail"+0
	rel 416+8 t=6 type.string+0
	rel 428+8 t=5 type.string+0
	rel 460+4 t=14 runtime.convT2E+0
	rel 532+4 t=14 runtime.gopanic+0
	rel 552+4 t=14 runtime.printlock+0
	rel 560+8 t=6 go.string."newfunc2 returned broken funcs"+0
	rel 572+8 t=5 go.string."newfunc2 returned broken funcs"+0
	rel 596+4 t=14 runtime.printstring+0
	rel 604+4 t=14 runtime.printnl+0
	rel 612+4 t=14 runtime.printunlock+0
	rel 628+8 t=6 "".fail+0
	rel 640+8 t=5 "".fail+0
	rel 664+4 t=14 runtime.printlock+0
	rel 672+8 t=6 go.string."newfunc returned broken funcs"+0
	rel 684+8 t=5 go.string."newfunc returned broken funcs"+0
	rel 708+4 t=14 runtime.printstring+0
	rel 716+4 t=14 runtime.printnl+0
	rel 724+4 t=14 runtime.printunlock+0
	rel 740+8 t=6 "".fail+0
	rel 752+8 t=5 "".fail+0
"".ff t=1 size=96 value=0 args=0x8 locals=0x10 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:33)	TEXT	"".ff(SB), $16-8
	0x0000 00000 (/Users/aram/go/test2/closure0.go:33)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:33)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:33)	MOVD	$-192, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:33)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:33)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:33)	FUNCDATA	$0, gclocals·0fb5f740dc3899c17d2f00dd94c805d6(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:33)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:34)	MOVD	ZR, "".autotmp_0012-16(SP)
	0x0018 00024 (/Users/aram/go/test2/closure0.go:34)	MOVD	ZR, "".autotmp_0012-8(SP)
	0x001c 00028 (/Users/aram/go/test2/closure0.go:34)	MOVD	$"".autotmp_0012-16(SP), R8
	0x0020 00032 (/Users/aram/go/test2/closure0.go:34)	MOVD	R8, R10
	0x0024 00036 (/Users/aram/go/test2/closure0.go:34)	MOVD	R10, R8
	0x0028 00040 (/Users/aram/go/test2/closure0.go:34)	MOVD	$"".ff.func1(SB), R9
	0x0040 00064 (/Users/aram/go/test2/closure0.go:34)	MOVD	R9, (R8)
	0x0044 00068 (/Users/aram/go/test2/closure0.go:34)	MOVD	R10, R8
	0x0048 00072 (/Users/aram/go/test2/closure0.go:34)	MOVD	"".x(FP), R9
	0x004c 00076 (/Users/aram/go/test2/closure0.go:34)	MOVD	R9, 8(R8)
	0x0050 00080 (/Users/aram/go/test2/closure0.go:36)	MOVD	R10, R8
	0x0054 00084 (/Users/aram/go/test2/closure0.go:36)	MOVD	R8, R9
	0x0058 00088 (/Users/aram/go/test2/closure0.go:37)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 40 9d e3 80 01  ..........?@....
	0x0010 fe 73 a8 77 c0 77 a7 ef c0 77 a7 f7 90 07 a7 ef  .s.w.w...w......
	0x0020 94 10 00 08 90 10 00 0a 0b 00 00 00 8a 11 60 00  ..............`.
	0x0030 8b 29 70 20 13 00 00 00 92 12 60 00 92 11 40 09  .)p ......`...@.
	0x0040 d2 72 20 00 90 10 00 0a d2 5f a8 af d2 72 20 08  .r ......_...r .
	0x0050 90 10 00 0a 92 10 00 08 81 c7 e0 08 81 e8 20 00  .............. .
	rel 40+8 t=6 "".ff.func1+0
	rel 52+8 t=5 "".ff.func1+0
"".call t=1 size=32 value=0 args=0x8 locals=0x0 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:39)	TEXT	"".call(SB), $0-8
	0x0000 00000 (/Users/aram/go/test2/closure0.go:39)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:39)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:39)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:39)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:39)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:39)	FUNCDATA	$0, gclocals·c2cb4d487f2f43f67e975313b3bca002(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:39)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:40)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 81 c7 e0 08 81 e8 20 00 00 00 00 00  .s.w...... .....
"".newfunc.func1 t=1 size=48 value=0 args=0x10 locals=0x0 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:13)	TEXT	"".newfunc.func1(SB), $0-16
	0x0000 00000 (/Users/aram/go/test2/closure0.go:13)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:13)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:13)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:13)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:13)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	FUNCDATA	$0, gclocals·aecfa9ecce04c513ee6b217848214030(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:13)	MOVD	$0, R8
	0x0018 00024 (/Users/aram/go/test2/closure0.go:13)	MOVD	"".x(FP), R8
	0x001c 00028 (/Users/aram/go/test2/closure0.go:13)	MOVD	R8, "".~r1+8(FP)
	0x0020 00032 (/Users/aram/go/test2/closure0.go:13)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 90 10 20 00 d0 5f a8 af d0 77 a8 b7  .s.w.. .._...w..
	0x0020 81 c7 e0 08 81 e8 20 00 00 00 00 00 00 00 00 00  ...... .........
"".newfunc2.func1 t=1 size=48 value=0 args=0x10 locals=0x0 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:14)	TEXT	"".newfunc2.func1(SB), $0-16
	0x0000 00000 (/Users/aram/go/test2/closure0.go:14)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:14)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:14)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:14)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:14)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	FUNCDATA	$0, gclocals·aecfa9ecce04c513ee6b217848214030(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:14)	MOVD	8(CTXT), R8
	0x0018 00024 (/Users/aram/go/test2/closure0.go:14)	MOVD	R8, R9
	0x001c 00028 (/Users/aram/go/test2/closure0.go:14)	MOVD	$0, R8
	0x0020 00032 (/Users/aram/go/test2/closure0.go:14)	MOVD	R9, R8
	0x0024 00036 (/Users/aram/go/test2/closure0.go:14)	MOVD	R8, "".~r1+8(FP)
	0x0028 00040 (/Users/aram/go/test2/closure0.go:14)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 d0 58 a0 08 92 10 00 08 90 10 20 00  .s.w.X........ .
	0x0020 90 10 00 09 d0 77 a8 b7 81 c7 e0 08 81 e8 20 00  .....w........ .
"".ff.func1 t=1 size=48 value=0 args=0x0 locals=0x0 leaf
	0x0000 00000 (/Users/aram/go/test2/closure0.go:34)	TEXT	"".ff.func1(SB), $0-0
	0x0000 00000 (/Users/aram/go/test2/closure0.go:34)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:34)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:34)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:34)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:34)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:34)	FUNCDATA	$0, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:34)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:34)	MOVD	8(CTXT), R8
	0x0018 00024 (/Users/aram/go/test2/closure0.go:34)	MOVD	R8, R9
	0x001c 00028 (/Users/aram/go/test2/closure0.go:36)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 d0 58 a0 08 92 10 00 08 81 c7 e0 08  .s.w.X..........
	0x0020 81 e8 20 00 00 00 00 00 00 00 00 00 00 00 00 00  .. .............
"".init t=1 size=208 value=0 args=0x0 locals=0x0
	0x0000 00000 (/Users/aram/go/test2/closure0.go:41)	TEXT	"".init(SB), $0-0
	0x0000 00000 (/Users/aram/go/test2/closure0.go:41)	RNOP
	0x0004 00004 (/Users/aram/go/test2/closure0.go:41)	RNOP
	0x0008 00008 (/Users/aram/go/test2/closure0.go:41)	MOVD	$-176, RT1
	0x000c 00012 (/Users/aram/go/test2/closure0.go:41)	SAVE	RT1, RSP, RSP
	0x0010 00016 (/Users/aram/go/test2/closure0.go:41)	MOVD	ILR, -56+176+2047(RSP)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:41)	FUNCDATA	$0, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:41)	FUNCDATA	$1, gclocals·2002e13acf59079a1a5782c918894579(SB)
	0x0014 00020 (/Users/aram/go/test2/closure0.go:41)	MOVUB	"".initdone·(SB), R8
	0x0030 00048 (/Users/aram/go/test2/closure0.go:41)	CMP	$1, R8
	0x0034 00052 (/Users/aram/go/test2/closure0.go:41)	BLEUW	68
	0x0038 00056 (/Users/aram/go/test2/closure0.go:41)	RNOP
	0x003c 00060 (/Users/aram/go/test2/closure0.go:41)	RETRESTORE
	0x0044 00068 (/Users/aram/go/test2/closure0.go:41)	MOVUB	"".initdone·(SB), R8
	0x0060 00096 (/Users/aram/go/test2/closure0.go:41)	CMP	$1, R8
	0x0064 00100 (/Users/aram/go/test2/closure0.go:41)	BNEW	120
	0x0068 00104 (/Users/aram/go/test2/closure0.go:41)	RNOP
	0x006c 00108 (/Users/aram/go/test2/closure0.go:41)	PCDATA	$0, $0
	0x006c 00108 (/Users/aram/go/test2/closure0.go:41)	CALL	runtime.throwinit(SB)
	0x0070 00112 (/Users/aram/go/test2/closure0.go:41)	RNOP
	0x0074 00116 (/Users/aram/go/test2/closure0.go:41)	UNDEF
	0x0078 00120 (/Users/aram/go/test2/closure0.go:41)	MOVD	$1, R8
	0x007c 00124 (/Users/aram/go/test2/closure0.go:41)	MOVUB	R8, R8
	0x0080 00128 (/Users/aram/go/test2/closure0.go:41)	MOVUB	R8, "".initdone·(SB)
	0x009c 00156 (/Users/aram/go/test2/closure0.go:41)	MOVD	$2, R8
	0x00a0 00160 (/Users/aram/go/test2/closure0.go:41)	MOVUB	R8, R8
	0x00a4 00164 (/Users/aram/go/test2/closure0.go:41)	MOVUB	R8, "".initdone·(SB)
	0x00c0 00192 (/Users/aram/go/test2/closure0.go:41)	RETRESTORE
	0x0000 01 00 00 00 01 00 00 00 82 10 3f 50 9d e3 80 01  ..........?P....
	0x0010 fe 73 a8 77 0b 00 00 00 8a 11 60 00 8b 29 70 20  .s.w......`..)p 
	0x0020 11 00 00 00 90 12 20 00 90 11 40 08 d0 0a 20 00  ...... ...@... .
	0x0030 80 a2 20 01 08 48 00 04 01 00 00 00 81 c7 e0 08  .. ..H..........
	0x0040 81 e8 20 00 0b 00 00 00 8a 11 60 00 8b 29 70 20  .. .......`..)p 
	0x0050 11 00 00 00 90 12 20 00 90 11 40 08 d0 0a 20 00  ...... ...@... .
	0x0060 80 a2 20 01 12 48 00 05 01 00 00 00 40 00 00 00  .. ..H......@...
	0x0070 01 00 00 00 00 0d ea d0 90 10 20 01 90 0a 20 ff  .......... ... .
	0x0080 0b 00 00 00 8a 11 60 00 8b 29 70 20 21 00 00 00  ......`..)p !...
	0x0090 a0 14 20 00 a0 11 40 10 d0 2c 20 00 90 10 20 02  .. ...@.., ... .
	0x00a0 90 0a 20 ff 0b 00 00 00 8a 11 60 00 8b 29 70 20  .. .......`..)p 
	0x00b0 21 00 00 00 a0 14 20 00 a0 11 40 10 d0 2c 20 00  !..... ...@.., .
	0x00c0 81 c7 e0 08 81 e8 20 00 00 00 00 00 00 00 00 00  ...... .........
	rel 20+8 t=6 "".initdone·+0
	rel 32+8 t=5 "".initdone·+0
	rel 68+8 t=6 "".initdone·+0
	rel 80+8 t=5 "".initdone·+0
	rel 108+4 t=14 runtime.throwinit+0
	rel 128+8 t=6 "".initdone·+0
	rel 140+8 t=5 "".initdone·+0
	rel 164+8 t=6 "".initdone·+0
	rel 176+8 t=5 "".initdone·+0
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·0fb5f740dc3899c17d2f00dd94c805d6 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 01 00 00 00 00              ............
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·aecfa9ecce04c513ee6b217848214030 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 02 00 00 00 00              ............
go.string.hdr."newfunc returned broken funcs" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 1d  ................
	rel 0+8 t=1 go.string."newfunc returned broken funcs"+0
go.string."newfunc returned broken funcs" t=8 dupok size=30 value=0
	0x0000 6e 65 77 66 75 6e 63 20 72 65 74 75 72 6e 65 64  newfunc returned
	0x0010 20 62 72 6f 6b 65 6e 20 66 75 6e 63 73 00         broken funcs.
go.string.hdr."newfunc2 returned broken funcs" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 1e  ................
	rel 0+8 t=1 go.string."newfunc2 returned broken funcs"+0
go.string."newfunc2 returned broken funcs" t=8 dupok size=31 value=0
	0x0000 6e 65 77 66 75 6e 63 32 20 72 65 74 75 72 6e 65  newfunc2 returne
	0x0010 64 20 62 72 6f 6b 65 6e 20 66 75 6e 63 73 00     d broken funcs.
go.string.hdr."fail" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 04  ................
	rel 0+8 t=1 go.string."fail"+0
go.string."fail" t=8 dupok size=5 value=0
	0x0000 66 61 69 6c 00                                   fail.
gclocals·68bfa0232bc220e249cf20baf5784e82 t=8 dupok size=24 value=0
	0x0000 00 00 00 04 00 00 00 04 00 00 00 00 02 00 00 00  ................
	0x0010 01 00 00 00 04 00 00 00                          ........
gclocals·4329624ce4271de83fc7c43fc9c7e126 t=8 dupok size=8 value=0
	0x0000 00 00 00 04 00 00 00 00                          ........
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·0fb5f740dc3899c17d2f00dd94c805d6 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 01 00 00 00 00              ............
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·c2cb4d487f2f43f67e975313b3bca002 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 01 01 00 00 00              ............
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·aecfa9ecce04c513ee6b217848214030 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 02 00 00 00 00              ............
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·aecfa9ecce04c513ee6b217848214030 t=8 dupok size=12 value=0
	0x0000 00 00 00 01 00 00 00 02 00 00 00 00              ............
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
gclocals·2002e13acf59079a1a5782c918894579 t=8 dupok size=8 value=0
	0x0000 00 00 00 01 00 00 00 00                          ........
"".fail t=31 size=1 value=0
"".initdone· t=31 size=1 value=0
"".newfunc·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".newfunc+0
"".newfunc2·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".newfunc2+0
"".main·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".main+0
"".ff·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".ff+0
"".call·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".call+0
"".newfunc.func1·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".newfunc.func1+0
"".newfunc2.func1·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".newfunc2.func1+0
"".ff.func1·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".ff.func1+0
"".init·f t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 "".init+0
runtime.gcbits.01 t=8 dupok size=1 value=0
	0x0000 01                                               .
go.string.hdr."func(int) int" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 0d  ................
	rel 0+8 t=1 go.string."func(int) int"+0
go.string."func(int) int" t=8 dupok size=14 value=0
	0x0000 66 75 6e 63 28 69 6e 74 29 20 69 6e 74 00        func(int) int.
type.func(int) int t=8 dupok size=136 value=0
	0x0000 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00 08  ................
	0x0010 87 32 3c 98 00 08 08 33 00 00 00 00 00 00 00 00  .2<....3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 0d 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 01  ................
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01  ................
	0x0070 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.algarray+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+8 t=1 go.string."func(int) int"+0
	rel 72+8 t=1 type.func(int) int+120
	rel 96+8 t=1 type.func(int) int+128
	rel 120+8 t=1 type.int+0
	rel 128+8 t=1 type.int+0
go.typelink.func(int) int	func(int) int t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 type.func(int) int+0
runtime.gcbits. t=8 dupok size=0 value=0
go.string.hdr."struct { F uintptr; x int }" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 1b  ................
	rel 0+8 t=1 go.string."struct { F uintptr; x int }"+0
go.string."struct { F uintptr; x int }" t=8 dupok size=28 value=0
	0x0000 73 74 72 75 63 74 20 7b 20 46 20 75 69 6e 74 70  struct { F uintp
	0x0010 74 72 3b 20 78 20 69 6e 74 20 7d 00              tr; x int }.
go.string.hdr.".F" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02  ................
	rel 0+8 t=1 go.string.".F"+0
go.string.".F" t=8 dupok size=3 value=0
	0x0000 2e 46 00                                         .F.
go.string.hdr."x" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01  ................
	rel 0+8 t=1 go.string."x"+0
go.string."x" t=8 dupok size=2 value=0
	0x0000 78 00                                            x.
type.struct { F uintptr; x int } t=8 dupok size=168 value=0
	0x0000 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00 00  ................
	0x0010 1f 44 df f0 00 08 08 99 00 00 00 00 00 00 00 00  .D..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 1b 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02  ................
	0x0050 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0090 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00a0 00 00 00 00 00 00 00 08                          ........
	rel 24+8 t=1 runtime.algarray+96
	rel 32+8 t=1 runtime.gcbits.+0
	rel 40+8 t=1 go.string."struct { F uintptr; x int }"+0
	rel 64+8 t=1 type.struct { F uintptr; x int }+88
	rel 88+8 t=1 go.string.hdr.".F"+0
	rel 96+8 t=1 go.importpath."".+0
	rel 104+8 t=1 type.uintptr+0
	rel 128+8 t=1 go.string.hdr."x"+0
	rel 136+8 t=1 go.importpath."".+0
	rel 144+8 t=1 type.int+0
go.string.hdr."*struct { F uintptr; x int }" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 1c  ................
	rel 0+8 t=1 go.string."*struct { F uintptr; x int }"+0
go.string."*struct { F uintptr; x int }" t=8 dupok size=29 value=0
	0x0000 2a 73 74 72 75 63 74 20 7b 20 46 20 75 69 6e 74  *struct { F uint
	0x0010 70 74 72 3b 20 78 20 69 6e 74 20 7d 00           ptr; x int }.
type.*struct { F uintptr; x int } t=8 dupok size=72 value=0
	0x0000 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00 08  ................
	0x0010 c4 99 ea 82 00 08 08 36 00 00 00 00 00 00 00 00  .......6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 1c 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.algarray+80
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+8 t=1 go.string."*struct { F uintptr; x int }"+0
	rel 64+8 t=1 type.struct { F uintptr; x int }+0
go.typelink.*struct { F uintptr; x int }	*struct { F uintptr; x int } t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 type.*struct { F uintptr; x int }+0
go.string.hdr."func()" t=8 dupok size=16 value=0
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 06  ................
	rel 0+8 t=1 go.string."func()"+0
go.string."func()" t=8 dupok size=7 value=0
	0x0000 66 75 6e 63 28 29 00                             func().
type.func() t=8 dupok size=120 value=0
	0x0000 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00 08  ................
	0x0010 f6 82 bc f6 00 08 08 33 00 00 00 00 00 00 00 00  .......3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 06 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.algarray+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+8 t=1 go.string."func()"+0
	rel 72+8 t=1 type.func()+120
	rel 96+8 t=1 type.func()+120
go.typelink.func()	func() t=8 dupok size=8 value=0
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 type.func()+0

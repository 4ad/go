// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import "cmd/internal/obj"

// General purpose registers, kept in the low bits of Prog.Reg.
const (
	// integer
	REG_R0 = obj.RBaseSPARC64 + iota
	REG_R1
	REG_R2
	REG_R3
	REG_R4
	REG_R5
	REG_R6
	REG_R7
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15
	REG_R16
	REG_R17
	REG_R18
	REG_R19
	REG_R20
	REG_R21
	REG_R22
	REG_R23
	REG_R24
	REG_R25
	REG_R26
	REG_R27
	REG_R28
	REG_R29
	REG_R30
	REG_R31

	// single-precision floating point
	REG_F0
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_F8
	REG_F9
	REG_F10
	REG_F11
	REG_F12
	REG_F13
	REG_F14
	REG_F15
	REG_F16
	REG_F17
	REG_F18
	REG_F19
	REG_F20
	REG_F21
	REG_F22
	REG_F23
	REG_F24
	REG_F25
	REG_F26
	REG_F27
	REG_F28
	REG_F29
	REG_F30
	REG_F31

	// double-precision floating point; the first half is aliased to
	// single-precision registers, that is: Dn is aliased to Fn, Fn+1,
	// where n ≤ 30.
	REG_D0
	REG_D32
	REG_D2
	REG_D34
	REG_D4
	REG_D36
	REG_D6
	REG_D38
	REG_D8
	REG_D40
	REG_D10
	REG_D42
	REG_D12
	REG_D44
	REG_D14
	REG_D46
	REG_D16
	REG_D48
	REG_D18
	REG_D50
	REG_D20
	REG_D52
	REG_D22
	REG_D54
	REG_D24
	REG_D56
	REG_D26
	REG_D58
	REG_D28
	REG_D60
	REG_D30
	REG_D62
)

const (
	// floating-point condition-code registers
	REG_FCC0 = REG_R0 + 256 + iota
	REG_FCC1
	REG_FCC2
	REG_FCC3
)

const (
	// integer condition-code flags
	REG_ICC = REG_R0 + 384
	REG_XCC = REG_R0 + 384 + 2
)

const (
	REG_SPECIAL = REG_R0 + 512

	REG_CCR  = REG_SPECIAL + 2
	REG_TICK = REG_SPECIAL + 4
	REG_RPC  = REG_SPECIAL + 5

	REG_LAST = REG_R0 + 1024
)

// Register assignments:
const (
	REG_ZR   = REG_R0
	REG_TMP  = REG_R1
	REG_RSP  = REG_R14
	REG_LR   = REG_R15
	REG_GOT  = REG_R23
	REG_RFP  = REG_R30
	REG_FTMP = REG_F0
	REG_DTMP = REG_D0
)

const (
	StackAlign            = 8      // ABI says 16
	StackBias             = 0x7ff  // craziness
	WindowSaveAreaSize    = 16 * 8 // only slots for RFP and PLR used
	ArgumentsSaveAreaSize = 0      // unused, normally 6 * 8
	MinStackFrameSize     = WindowSaveAreaSize + ArgumentsSaveAreaSize
)

// Prog.mark
const (
	FOLL = 1 << iota
	LABEL
	LEAF
)

const (
	ClassUnknown = iota

	ClassReg        // R1..R31
	ClassFloatReg   // F0..F31
	ClassDoubleReg  // D0..D62
	ClassCond       // ICC, XCC
	ClassFloatCond  // FCC0..FCC3
	ClassSpecialReg // TICK, CCR, etc

	ClassZero       // $0 or ZR
	ClassConst5     // unsigned 5-bit constant
	ClassConst6     // unsigned 6-bit constant
	ClassConst13    // signed 13-bit constant
	ClassConst31_   // signed 32-bit constant, negative
	ClassConst31    // signed 32-bit constant, positive or zero
	ClassConst32    // 32-bit constant
	ClassConst      // 64-bit constant
	ClassFloatConst // floating-point constant

	ClassRegReg     // $(Rn+Rm) or $(Rn)(Rm*1)
	ClassRegConst13 // $n(R), n is 13-bit signed
	ClassRegConst   // $n(R), n large

	ClassIndirRegReg // (Rn+Rm) or (Rn)(Rm*1)
	ClassIndir0      // (R)
	ClassIndir13     // n(R), n is 13-bit signed
	ClassIndir       // n(R), n large

	ClassShortBranch // n(PC) branch target, n is 21-bit signed, mod 4

	ClassAddr // $sym(SB)
	ClassMem  // sym(SB)

	ClassTextSize
	ClassNone
)

var cnames = []string{
	ClassUnknown:     "ClassUnknown",
	ClassReg:         "ClassReg",
	ClassFloatReg:    "ClassFloatReg",
	ClassDoubleReg:   "ClassDoubleReg",
	ClassCond:        "ClassCond",
	ClassFloatCond:   "ClassFloatCond",
	ClassSpecialReg:  "ClassSpecialReg",
	ClassZero:        "ClassZero",
	ClassConst5:      "ClassConst5",
	ClassConst6:      "ClassConst6",
	ClassConst13:     "ClassConst13",
	ClassConst31_:    "ClassConst31-",
	ClassConst31:     "ClassConst31+",
	ClassConst32:     "ClassConst32",
	ClassConst:       "ClassConst",
	ClassFloatConst:  "ClassFloatConst",
	ClassRegReg:      "ClassRegReg",
	ClassRegConst13:  "ClassRegConst13",
	ClassRegConst:    "ClassRegConst",
	ClassIndirRegReg: "ClassIndirRegReg",
	ClassIndir0:      "ClassIndir0",
	ClassIndir13:     "ClassIndir13",
	ClassIndir:       "ClassIndir",
	ClassShortBranch: "ClassShortBranch",
	ClassAddr:        "ClassAddr",
	ClassMem:         "ClassMem",
	ClassTextSize:    "ClassTextSize",
	ClassNone:        "ClassNone",
}

//go:generate go run ../stringer.go -i $GOFILE -o anames.go -p sparc64

const (
	AADD = obj.ABaseSPARC64 + obj.A_ARCHSPECIFIC + iota
	AADDCC
	AADDC
	AADDCCC
	AAND
	AANDCC
	AANDN
	AANDNCC

	// These are the two-operand SPARCv9 32-, and 64-bit, branch
	// on integer condition codes with prediction (BPcc), not the
	// single-operand SPARCv8 32-bit branch on integer condition
	// codes (Bicc).
	ABN
	ABNE
	ABE
	ABG
	ABLE
	ABGE
	ABL
	ABGU
	ABLEU
	ABCC
	ABCS
	ABPOS
	ABNEG
	ABVC
	ABVS

	ABRZ
	ABRLEZ
	ABRLZ
	ABRNZ
	ABRGZ
	ABRGEZ
	ACASW
	ACASD
	AFABSS
	AFABSD
	AFADDS
	AFADDD
	AFBA
	AFBN
	AFBU
	AFBG
	AFBUG
	AFBL
	AFBUL
	AFBLG
	AFBNE
	AFBE
	AFBUE
	AFBGE
	AFBUGE
	AFBLE
	AFBULE
	AFBO
	AFCMPS
	AFCMPD
	AFDIVS
	AFDIVD
	AFITOS
	AFITOD
	AFLUSH
	AFMOVS // the SPARC64 instruction, and alias for loads and stores
	AFMOVD // the SPARC64 instruction, and alias for loads and stores
	AFMULS
	AFMULD
	AFSMULD
	AFNEGS
	AFNEGD
	AFSQRTS
	AFSQRTD
	AFSTOX
	AFDTOX
	AFSTOI
	AFDTOI
	AFSTOD
	AFDTOS
	AFSUBS
	AFSUBD
	AFXTOS
	AFXTOD
	AJMPL
	ALDSB
	ALDSH
	ALDSW
	ALDUB
	ALDUH
	ALDUW
	ALDD
	ALDSF
	ALDDF
	AMEMBAR
	AMULD
	ASDIVD
	AUDIVD
	AOR
	AORCC
	AORN
	AORNCC
	ARD
	ASETHI
	ASLLW
	ASRLW
	ASRAW
	ASLLD
	ASRLD
	ASRAD
	ASTB
	ASTH
	ASTW
	ASTD
	ASTSF
	ASTDF
	ASUB
	ASUBCC
	ASUBC
	ASUBCCC
	ATA
	AXOR
	AXORCC
	AXNOR
	AXNORCC

	// Pseudo-instructions, aliases to SPARC64 instructions and
	// synthetic instructions.
	ANEG
	AMOVUB
	AMOVB
	AMOVUH
	AMOVH
	AMOVUW
	AMOVW
	AMOVD // also the SPARC64 synthetic instruction
	ARNOP // SETHI $0, ZR

	AWORD
	ADWORD
)

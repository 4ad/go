// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import "cmd/internal/obj"

// General purpose registers, kept in the low bits of Prog.Reg.
const (
	// integer
	REG_G0 = obj.RBaseSPARC64 + iota
	REG_G1
	REG_G2
	REG_G3
	REG_G4
	REG_G5
	REG_G6
	REG_G7
	REG_O0
	REG_O1
	REG_O2
	REG_O3
	REG_O4
	REG_O5
	REG_O6
	REG_O7
	REG_L0
	REG_L1
	REG_L2
	REG_L3
	REG_L4
	REG_L5
	REG_L6
	REG_L7
	REG_I0
	REG_I1
	REG_I2
	REG_I3
	REG_I4
	REG_I5
	REG_I6
	REG_I7

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

	// common single/double-precision virtualized registers.
	// Yn is aliased to F2n, F2n+1, D2n.
	REG_Y0
	REG_Y1
	REG_Y2
	REG_Y3
	REG_Y4
	REG_Y5
	REG_Y6
	REG_Y7
	REG_Y8
	REG_Y9
	REG_Y10
	REG_Y11
	REG_Y12
	REG_Y13
	REG_Y14
	REG_Y15
)

const (
	// floating-point condition-code registers
	REG_FCC0 = REG_G0 + 256 + iota
	REG_FCC1
	REG_FCC2
	REG_FCC3
)

const (
	// integer condition-code flags
	REG_ICC = REG_G0 + 384
	REG_XCC = REG_G0 + 384 + 2
)

const (
	REG_SPECIAL = REG_G0 + 512

	REG_CCR  = REG_SPECIAL + 2
	REG_TICK = REG_SPECIAL + 4
	REG_RPC  = REG_SPECIAL + 5

	REG_BSP = REG_RSP + 256
	REG_BFP = REG_RFP + 256

	REG_LAST = REG_G0 + 1024
)

// Register assignments:
const (
	REG_ZR   = REG_G0
	REG_RT1  = REG_G1
	REG_CTXT = REG_G2
	REG_G    = REG_G3
	REG_RT2  = REG_G4
	REG_TMP  = REG_G5
	REG_TLS  = REG_G7
	REG_RSP  = REG_O6
	REG_OLR  = REG_O7
	REG_TMP2 = REG_L0
	REG_RFP  = REG_I6
	REG_ILR  = REG_I7
	REG_FTMP = REG_F30
	REG_DTMP = REG_D30
	REG_YTMP = REG_Y15
	REG_YTWO = REG_Y14
)

const (
	REG_MIN = REG_G0
	REG_MAX = REG_I5
)

const (
	StackBias             = 0x7ff  // craziness
	WindowSaveAreaSize    = 16 * 8 // only slots for RFP and PLR used
	ArgumentsSaveAreaSize = 6 * 8  // unused
	MinStackFrameSize     = WindowSaveAreaSize + ArgumentsSaveAreaSize
)

const (
	BIG = 1<<12 - 1 // magnitude of smallest negative immediate
)

// Prog.mark
const (
	FOLL = 1 << iota
	LABEL
	LEAF
)

const (
	ClassUnknown = iota

	ClassReg    // R1..R31
	ClassFReg   // F0..F31
	ClassDReg   // D0..D62
	ClassCond   // ICC, XCC
	ClassFCond  // FCC0..FCC3
	ClassSpcReg // TICK, CCR, etc

	ClassZero     // $0 or ZR
	ClassConst5   // unsigned 5-bit constant
	ClassConst6   // unsigned 6-bit constant
	ClassConst10  // signed 10-bit constant
	ClassConst11  // signed 11-bit constant
	ClassConst13  // signed 13-bit constant
	ClassConst31_ // signed 32-bit constant, negative
	ClassConst31  // signed 32-bit constant, positive or zero
	ClassConst32  // 32-bit constant
	ClassConst    // 64-bit constant
	ClassFConst   // floating-point constant

	ClassRegReg     // $(Rn+Rm) or $(Rn)(Rm*1)
	ClassRegConst13 // $n(R), n is 13-bit signed
	ClassRegConst   // $n(R), n large

	ClassIndirRegReg // (Rn+Rm) or (Rn)(Rm*1)
	ClassIndir0      // (R)
	ClassIndir13     // n(R), n is 13-bit signed
	ClassIndir       // n(R), n large

	ClassBranch      // n(PC) branch target, n is 21-bit signed, mod 4
	ClassLargeBranch // n(PC) branch target, n is 32-bit signed, mod 4

	ClassAddr    // $sym(SB)
	ClassMem     // sym(SB)
	ClassTLSAddr // $tlssym(SB)
	ClassTLSMem  // tlssym(SB)

	ClassTextSize
	ClassNone

	ClassBias = 64 // BFP or BSP present in Addr, bitwise OR with classes above
)

var cnames = []string{
	ClassUnknown:     "ClassUnknown",
	ClassReg:         "ClassReg",
	ClassFReg:        "ClassFReg",
	ClassDReg:        "ClassDReg",
	ClassCond:        "ClassCond",
	ClassFCond:       "ClassFCond",
	ClassSpcReg:      "ClassSpcReg",
	ClassZero:        "ClassZero",
	ClassConst5:      "ClassConst5",
	ClassConst6:      "ClassConst6",
	ClassConst10:     "ClassConst10",
	ClassConst11:     "ClassConst11",
	ClassConst13:     "ClassConst13",
	ClassConst31_:    "ClassConst31-",
	ClassConst31:     "ClassConst31+",
	ClassConst32:     "ClassConst32",
	ClassConst:       "ClassConst",
	ClassFConst:      "ClassFConst",
	ClassRegReg:      "ClassRegReg",
	ClassRegConst13:  "ClassRegConst13",
	ClassRegConst:    "ClassRegConst",
	ClassIndirRegReg: "ClassIndirRegReg",
	ClassIndir0:      "ClassIndir0",
	ClassIndir13:     "ClassIndir13",
	ClassIndir:       "ClassIndir",
	ClassBranch:      "ClassBranch",
	ClassLargeBranch: "ClassLargeBranch",
	ClassAddr:        "ClassAddr",
	ClassMem:         "ClassMem",
	ClassTLSAddr:     "ClassTLSAddr",
	ClassTLSMem:      "ClassTLSMem",
	ClassTextSize:    "ClassTextSize",
	ClassNone:        "ClassNone",
	ClassBias:        "ClassBias",
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
	AFLUSHW
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
	ALDD
	ALDDF
	ALDSF
	ALDUH
	ALDUW
	AMEMBAR
	AMOVA
	AMOVCC
	AMOVCS
	AMOVE
	AMOVG
	AMOVGE
	AMOVGU
	AMOVL
	AMOVLE
	AMOVLEU
	AMOVN
	AMOVNE
	AMOVNEG
	AMOVPOS
	AMOVRGEZ
	AMOVRGZ
	AMOVRLEZ
	AMOVRLZ
	AMOVRNZ
	AMOVRZ
	AMOVVC
	AMOVVS
	AMULD
	AOR
	AORCC
	AORN
	AORNCC
	ARD
	ARESTORE // not used under normal circumstances
	ASAVE    // not used under normal circumstances
	ASDIVD
	ASETHI
	AUDIVD
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
	ACMP // SUBCC R1, R2, ZR
	ANEG
	AMOVUB
	AMOVB
	AMOVUH
	AMOVH
	AMOVUW
	AMOVW
	AMOVD // also the SPARC64 synthetic instruction
	ARNOP // SETHI $0, ZR

	// These are aliases to two-operand SPARCv9 32-, and 64-bit,
	// branch on integer condition codes with prediction (BPcc),
	// with ICC implied.
	ABNW
	ABNEW
	ABEW
	ABGW
	ABLEW
	ABGEW
	ABLW
	ABGUW
	ABLEUW
	ABCCW
	ABCSW
	ABPOSW
	ABNEGW
	ABVCW
	ABVSW

	// These are aliases to two-operand SPARCv9 32-, and 64-bit,
	// branch on integer condition codes with prediction (BPcc),
	// with XCC implied.
	ABND
	ABNED
	ABED
	ABGD
	ABLED
	ABGED
	ABLD
	ABGUD
	ABLEUD
	ABCCD
	ABCSD
	ABPOSD
	ABNEGD
	ABVCD
	ABVSD

	AWORD
	ADWORD

	// JMPL $8(ILR), ZR
	//      RESTORE $0, ZR, ZR
	ARETRESTORE

	ALAST
)

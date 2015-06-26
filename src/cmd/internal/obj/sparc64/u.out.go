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

	// floating point
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

	REG_BSP = REG_R14 + 64
	REG_BFP = REG_R30 + 64
)

// Register assignments:
const (
	RegZero = REG_R0
	RegRSP  = REG_R14
	RegLink = REG_R15
	RegFP   = REG_R30
)

const (
	ClassUnknown   = iota
	ClassReg       // R1..R31
	ClassFloatReg  // F0..F31
	ClassBiased    // BSP or BFP
	ClassPairComma // (Rn, Rn+1)
	ClassPairPlus  // (Rn+Rm)

	ClassZero       // $0 or ZR
	ClassConst13    // signed 13-bit constant
	ClassConst      // 64-bit constant
	ClassFloatConst // floting-point constant

	ClassEffectiveAddr13 // $n(R), n is 13-bit signed
	ClassEffectiveAddr   // $n(R), n large

	ClassIndir0  // (R)
	ClassIndir13 // n(R), n is 13-bit signed
	ClassIndir   // n(R), n large

	ClassAddr // $sym(SB)
	ClassMem  // sym(SB)

	ClassTextSize
	ClassNone
)

var cnames = []string{
	ClassUnknown:         "ClassUnknown",
	ClassReg:             "ClassReg",
	ClassFloatReg:        "ClassFloatReg",
	ClassBiased:          "ClassBiased",
	ClassPairComma:       "ClassPairComma",
	ClassPairPlus:        "ClassPairPlus",
	ClassZero:            "ClassZero",
	ClassConst13:         "ClassConst13",
	ClassConst:           "ClassConst",
	ClassFloatConst:      "ClassFloatConst",
	ClassEffectiveAddr13: "ClassEffectiveAddr13",
	ClassEffectiveAddr:   "ClassEffectiveAddr",
	ClassIndir0:          "ClassIndir0",
	ClassIndir13:         "ClassIndir13",
	ClassIndir:           "ClassIndir",
	ClassAddr:            "ClassAddr",
	ClassMem:             "ClassMem",
	ClassTextSize:        "ClassTextSize",
	ClassNone:            "ClassNone",
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
	AFWTOS
	AFWTOD
	AFLUSH
	AFMOVS
	AFMOVD
	AFMULS
	AFMULD
	AFSMULD
	AFNEGS
	AFNEGD
	AFSQRTS
	AFSQRTD
	AFSTOXD
	AFDTOXD
	AFSTOXW
	AFDTOXW
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
	ARDCCR
	ARDTICK
	ARDPC
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
	AXOR
	AXORCC
	AXNOR
	AXNORCC
)

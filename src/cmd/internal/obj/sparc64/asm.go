// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"errors"
	"fmt"
)

type Optab struct {
	as int16
	a1 int8
	a2 int8
	a3 int8
}

var optab = map[Optab]int{
	Optab{obj.ATEXT, ClassAddr, ClassNone, ClassTextSize}: 0,

	Optab{AADD, ClassReg, ClassNone, ClassReg}:  1,
	Optab{AAND, ClassReg, ClassNone, ClassReg}:  1,
	Optab{AMULD, ClassReg, ClassNone, ClassReg}: 1,
	Optab{AADD, ClassReg, ClassReg, ClassReg}:   1,
	Optab{AAND, ClassReg, ClassReg, ClassReg}:   1,
	Optab{AMULD, ClassReg, ClassReg, ClassReg}:  1,
	Optab{ASLLD, ClassReg, ClassReg, ClassReg}:  1,
	Optab{ASLLW, ClassReg, ClassReg, ClassReg}:  1,

	Optab{AFADDD, ClassDoubleReg, ClassNone, ClassDoubleReg}:      1,
	Optab{AFADDD, ClassDoubleReg, ClassDoubleReg, ClassDoubleReg}: 1,
	Optab{AFSMULD, ClassFloatReg, ClassFloatReg, ClassDoubleReg}:  1,

	Optab{AMOVD, ClassReg, ClassNone, ClassReg}: 2,

	Optab{AADD, ClassReg, ClassConst13, ClassReg}:  3,
	Optab{AAND, ClassReg, ClassConst13, ClassReg}:  3,
	Optab{AMULD, ClassReg, ClassConst13, ClassReg}: 3,
	Optab{ASLLD, ClassReg, ClassConst6, ClassReg}:  3,
	Optab{ASLLW, ClassReg, ClassConst5, ClassReg}:  3,

	Optab{AMOVD, ClassConst13, ClassNone, ClassReg}: 4,

	Optab{ALDD, ClassPairPlus, ClassNone, ClassReg}:        5,
	Optab{ASTD, ClassReg, ClassNone, ClassPairPlus}:        6,
	Optab{ALDDF, ClassPairPlus, ClassNone, ClassDoubleReg}: 5,
	Optab{ASTDF, ClassDoubleReg, ClassNone, ClassPairPlus}: 6,

	Optab{ALDD, ClassIndir13, ClassNone, ClassReg}:        7,
	Optab{ASTD, ClassReg, ClassNone, ClassIndir13}:        8,
	Optab{ALDDF, ClassIndir13, ClassNone, ClassDoubleReg}: 7,
	Optab{ASTDF, ClassDoubleReg, ClassNone, ClassIndir13}: 8,

	Optab{ARD, ClassSpecialReg, ClassNone, ClassReg}: 9,

	Optab{ACASD, ClassIndir0, ClassReg, ClassReg}: 10,

	Optab{AFSTOD, ClassFloatReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFDTOS, ClassDoubleReg, ClassNone, ClassFloatReg}: 11,

	Optab{AFMOVD, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,

	Optab{AFXTOD, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFITOD, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFXTOS, ClassDoubleReg, ClassNone, ClassFloatReg}:  11,
	Optab{AFITOS, ClassFloatReg, ClassNone, ClassFloatReg}:   11,

	Optab{AFSTOX, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFDTOX, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFDTOI, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,
	Optab{AFSTOI, ClassFloatReg, ClassNone, ClassFloatReg}:   11,

	Optab{AFABSD, ClassDoubleReg, ClassNone, ClassDoubleReg}: 11,

	Optab{ASETHI, ClassConst32, ClassNone, ClassReg}: 12,

	Optab{AMEMBAR, ClassConst, ClassNone, ClassNone}: 13,

	Optab{AFCMPD, ClassDoubleReg, ClassDoubleReg, ClassFloatCondReg}: 14,
	Optab{AFCMPD, ClassDoubleReg, ClassDoubleReg, ClassNone}:         14,

	Optab{AMOVD, ClassConst32, ClassNone, ClassReg}:  15,
	Optab{AMOVD, ClassConst31_, ClassNone, ClassReg}: 16,
}

// Compatible classes, if something accepts a $hugeconst, it
// can also accept $smallconst, $0 and ZR. Something that accepts a
// register, can also accept $0, etc.
var cc = map[int8][]int8{
	ClassReg:           {ClassZero},
	ClassConst13:       {ClassConst6, ClassConst5, ClassZero},
	ClassConst31:       {ClassConst6, ClassConst5, ClassZero},
	ClassConst32:       {ClassConst31_, ClassConst31, ClassConst13, ClassConst6, ClassConst5, ClassZero},
	ClassConst:         {ClassConst32, ClassConst31_, ClassConst31, ClassConst13, ClassConst6, ClassConst5, ClassZero},
	ClassEffectiveAddr: {ClassEffectiveAddr13},
	ClassIndir13:       {ClassIndir0},
	ClassIndir:         {ClassIndir13, ClassIndir0},
}

var isInstDouble = map[int16]bool{
	AFADDD:  true,
	AFSUBD:  true,
	AFABSD:  true,
	AFCMPD:  true,
	AFDIVD:  true,
	AFMOVD:  true,
	AFMULD:  true,
	AFNEGD:  true,
	AFSQRTD: true,
	ALDDF:   true,
	ASTDF:   true,
}

var isInstFloat = map[int16]bool{
	AFADDS:  true,
	AFSUBS:  true,
	AFABSS:  true,
	AFCMPS:  true,
	AFDIVS:  true,
	AFMOVS:  true,
	AFMULS:  true,
	AFNEGS:  true,
	AFSQRTS: true,
	ALDSF:   true,
	ASTSF:   true,
}

// Compatible instructions, if an asm* function accepts AADD,
// it accepts ASUBCCC too.
var ci = map[int16][]int16{
	AADD:     {AADDCC, AADDC, AADDCCC, ASUB, ASUBCC, ASUBC, ASUBCCC},
	AAND:     {AANDCC, AANDN, AANDNCC, AOR, AORCC, AORN, AORNCC, AXOR, AXORCC, AXNOR, AXNORCC},
	obj.AJMP: {ABN, ABNE, ABE, ABG, ABLE, ABGE, ABL, ABGU, ABLEU, ABCC, ABCS, ABPOS, ABNEG, ABVC, ABVS},
	ABRZ:     {ABRLEZ, ABRLZ, ABRNZ, ABRGZ, ABRGEZ},
	ACASD:    {ACASW},
	AFABSD:   {AFABSS, AFNEGD, AFNEGS, AFSQRTD, AFNEGS},
	AFADDD:   {AFADDS, AFSUBS, AFSUBD, AFMULD, AFMULS, AFSMULD, AFDIVD, AFDIVS},
	AFBA:     {AFBN, AFBU, AFBG, AFBUG, AFBL, AFBUL, AFBLG, AFBNE, AFBE, AFBUE, AFBGE, AFBUGE, AFBLE, AFBULE, AFBO},
	AFCMPD:   {AFCMPS},
	AFITOD:   {AFITOS},
	AFMOVD:   {AFMOVS},
	AFSTOD:   {AFDTOS},
	AFXTOD:   {AFXTOS},
	ALDD:     {ALDSB, ALDSH, ALDSW, ALDUB, ALDUH, ALDUW, AMOVSB, AMOVSH, AMOVSW, AMOVB, AMOVH, AMOVW, AMOVD},
	ALDDF:    {ALDSF, AFMOVD, AFMOVS},
	AMULD:    {ASDIVD, AUDIVD},
	ARD:      {AMOVD},
	ASLLD:    {ASRLD, ASRAD},
	ASLLW:    {ASLLW, ASRLW, ASRAW},
	ASTD:     {ASTB, ASTH, ASTW, AMOVB, AMOVH, AMOVW, AMOVD},
	ASTDF:    {ASTSF, AFMOVD, AFMOVS},
}

func init() {
	// For each line in optab, duplicate it so that we'll also
	// have a line that will accept compatible instructions, but
	// only if there isn't an already existent line with the same
	// key. Also change operand type, if the instruction is a double.
	for o, v := range optab {
		for _, c := range ci[o.as] {
			do := o
			do.as = c
			if isInstDouble[o.as] && isInstFloat[do.as] {
				if do.a1 == ClassDoubleReg {
					do.a1 = ClassFloatReg
				}
				if do.a2 == ClassDoubleReg {
					do.a2 = ClassFloatReg
				}
				if do.a3 == ClassDoubleReg {
					do.a3 = ClassFloatReg
				}
			}
			_, ok := optab[do]
			if !ok {
				optab[do] = v
			}
		}
	}
	// For each line in optab that accepts a large-class operand,
	// duplicate it so that we'll also have a line that accepts a
	// small-class operand, but do it only if there isn't an already
	// existent line with the same key.
	for o, v := range optab {
		for _, c := range cc[o.a1] {
			do := o
			do.a1 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = v
			}
		}
	}
	for o, v := range optab {
		for _, c := range cc[o.a2] {
			do := o
			do.a2 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = v
			}
		}
	}
	for o, v := range optab {
		for _, c := range cc[o.a3] {
			do := o
			do.a3 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = v
			}
		}
	}
}

func oplook(p *obj.Prog) (int, error) {
	o := Optab{as: p.As, a1: p.From.Class, a2: ClassNone, a3: p.To.Class}
	if p.From3 != nil {
		o.a2 = p.From3.Class
	}
	v, ok := optab[o]
	if !ok {
		return 0, fmt.Errorf("illegal combination %v %v %v %v, %d %d", p, DRconv(o.a1), DRconv(o.a2), DRconv(o.a3), p.From.Type, p.To.Type)
	}
	return v, nil
}

func ir(imm22 uint32, rd int16) uint32 {
	return uint32(rd)&31<<25 | uint32(imm22&(1<<23-1))
}

func d22(a, disp22 int) uint32 {
	return uint32(a&1<<29 | disp22&(1<<23-1))
}

func d19(a, cc1, cc0, p, disp19 int) uint32 {
	return uint32(a&1<<29 | cc1&1<<21 | cc0&1<<20 | p&1<<19 | disp19&(1<<20-1))
}

func d30(disp30 int) uint32 {
	return uint32(disp30 & (1<<31 - 1))
}

func rrr(rs1, imm_asi, rs2, rd int16) uint32 {
	return uint32(uint32(rd)&31<<25 | uint32(rs1)&31<<14 | uint32(imm_asi)&255<<5 | uint32(rs2)&31)
}

func rsr(rs1 int16, simm13 int64, rd int16) uint32 {
	return uint32(int(rd)&31<<25 | int(rs1)&31<<14 | 1<<13 | int(simm13)&(1<<14-1))
}

func rd(r int16) uint32 {
	return uint32(int(r) & 31 << 25)
}

func op(op int) uint32 {
	return uint32(op << 30)
}

func op3(op, op3 int) uint32 {
	return uint32(op<<30 | op3<<19)
}

func op2(op2 int) uint32 {
	return uint32(op2 << 22)
}

func cond(cond int) uint32 {
	return uint32(cond << 25)
}

func opf(opf int) uint32 {
	return uint32(opf << 5)
}

func opload(a int16) uint32 {
	switch a {
	// Load integer.
	case ALDSB, AMOVSB:
		return op3(3, 9)
	case ALDSH, AMOVSH:
		return op3(3, 10)
	case ALDSW, AMOVSW:
		return op3(3, 8)
	case ALDUB, AMOVB:
		return op3(3, 1)
	case ALDUH, AMOVH:
		return op3(3, 2)
	case ALDUW, AMOVW:
		return op3(3, 0)
	case ALDD, AMOVD:
		return op3(3, 11)

	// Load floating-point register.
	case ALDSF, AFMOVS:
		return op3(3, 0x20)
	case ALDDF, AFMOVD:
		return op3(3, 0x23)

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
}

func opstore(a int16) uint32 {
	switch a {
	// Store Integer.
	case ASTB, AMOVB:
		return op3(3, 5)
	case ASTH, AMOVH:
		return op3(3, 6)
	case ASTW, AMOVW:
		return op3(3, 4)
	case ASTD, AMOVD:
		return op3(3, 14)

	// Store floating-point.
	case ASTSF, AFMOVS:
		return op3(3, 0x24)
	case ASTDF, AFMOVD:
		return op3(3, 0x27)

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
}

func oprd(a int16) uint32 {
	switch a {
	// Read ancillary state register.
	case ARD, AMOVD:
		return op3(2, 0x28)

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
}

func opalu(a int16) uint32 {
	switch a {
	// Add.
	case AADD:
		return op3(2, 0)
	case AADDCC:
		return op3(2, 16)
	case AADDC:
		return op3(2, 8)
	case AADDCCC:
		return op3(2, 24)

	// AND logical operation.
	case AAND:
		return op3(2, 1)
	case AANDCC:
		return op3(2, 17)
	case AANDN:
		return op3(2, 5)
	case AANDNCC:
		return op3(2, 21)

	// Multiply and divide.
	case AMULD:
		return op3(2, 9)
	case ASDIVD:
		return op3(2, 0x2D)
	case AUDIVD:
		return op3(2, 0xD)

	// OR logical operation.
	case AOR, AMOVD:
		return op3(2, 2)
	case AORCC:
		return op3(2, 18)
	case AORN:
		return op3(2, 6)
	case AORNCC:
		return op3(2, 22)

	// Subtract.
	case ASUB:
		return op3(2, 4)
	case ASUBCC:
		return op3(2, 20)
	case ASUBC:
		return op3(2, 12)
	case ASUBCCC:
		return op3(2, 28)

	// XOR logical operation.
	case AXOR:
		return op3(2, 3)
	case AXORCC:
		return op3(2, 19)
	case AXNOR:
		return op3(2, 7)
	case AXNORCC:
		return op3(2, 23)

	// Floating-Point Add
	case AFADDS:
		return op3(2, 0x34) | opf(0x41)
	case AFADDD:
		return op3(2, 0x34) | opf(0x42)

	// Floating-point subtract.
	case AFSUBS:
		return op3(2, 0x34) | opf(0x45)
	case AFSUBD:
		return op3(2, 0x34) | opf(0x46)

	// Floating-point divide.
	case AFDIVS:
		return op3(2, 0x34) | opf(0x4D)
	case AFDIVD:
		return op3(2, 0x34) | opf(0x4E)

	// Floating-point multiply.
	case AFMULS:
		return op3(2, 0x34) | opf(0x49)
	case AFMULD:
		return op3(2, 0x34) | opf(0x4A)
	case AFSMULD:
		return op3(2, 0x34) | opf(0x69)

	// Shift.
	case ASLLW:
		return op3(2, 0x25)
	case ASRLW:
		return op3(2, 0x26)
	case ASRAW:
		return op3(2, 0x27)
	case ASLLD:
		return op3(2, 0x25) | 1<<12
	case ASRLD:
		return op3(2, 0x26) | 1<<12
	case ASRAD:
		return op3(2, 0x27) | 1<<12

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
}

func opcode(a int16) uint32 {
	switch a {
	// Branch on integer condition codes with prediction (BPcc).
	case obj.AJMP:
		return cond(8) | op2(1)
	case ABN:
		return cond(0) | op2(1)
	case ABNE:
		return cond(9) | op2(1)
	case ABE:
		return cond(1) | op2(1)
	case ABG:
		return cond(10) | op2(1)
	case ABLE:
		return cond(2) | op2(1)
	case ABGE:
		return cond(11) | op2(1)
	case ABL:
		return cond(3) | op2(1)
	case ABGU:
		return cond(12) | op2(1)
	case ABLEU:
		return cond(4) | op2(1)
	case ABCC:
		return cond(13) | op2(1)
	case ABCS:
		return cond(5) | op2(1)
	case ABPOS:
		return cond(14) | op2(1)
	case ABNEG:
		return cond(6) | op2(1)
	case ABVC:
		return cond(15) | op2(1)
	case ABVS:
		return cond(7) | op2(1)

	// Branch on integer register with prediction (BPr).
	case ABRZ:
		return cond(1) | op2(3)
	case ABRLEZ:
		return cond(2) | op2(3)
	case ABRLZ:
		return cond(3) | op2(3)
	case ABRNZ:
		return cond(5) | op2(3)
	case ABRGZ:
		return cond(6) | op2(3)
	case ABRGEZ:
		return cond(7) | op2(3)

	// Call and link
	case obj.ACALL:
		return op(1)

	case ACASW:
		return op3(3, 0x3C) | 1<<13
	case ACASD:
		return op3(3, 0x3E) | 1<<13

	case AFABSS:
		return op3(2, 0x34) | opf(9)
	case AFABSD:
		return op3(2, 0x34) | opf(10)

	// Branch on floating-point condition codes (FBfcc).
	case AFBA:
		return cond(8) | op2(6)
	case AFBN:
		return cond(0) | op2(6)
	case AFBU:
		return cond(7) | op2(6)
	case AFBG:
		return cond(6) | op2(6)
	case AFBUG:
		return cond(5) | op2(6)
	case AFBL:
		return cond(4) | op2(6)
	case AFBUL:
		return cond(3) | op2(6)
	case AFBLG:
		return cond(2) | op2(6)
	case AFBNE:
		return cond(1) | op2(6)
	case AFBE:
		return cond(9) | op2(6)
	case AFBUE:
		return cond(10) | op2(6)
	case AFBGE:
		return cond(11) | op2(6)
	case AFBUGE:
		return cond(12) | op2(6)
	case AFBLE:
		return cond(13) | op2(6)
	case AFBULE:
		return cond(14) | op2(6)
	case AFBO:
		return cond(15) | op2(6)

	// Floating-point compare.
	case AFCMPS:
		return op3(2, 0x35) | opf(0x51)
	case AFCMPD:
		return op3(2, 0x35) | opf(0x52)

	// Convert 32-bit integer to floating point.
	case AFITOS:
		return op3(2, 0x34) | opf(0xC4)
	case AFITOD:
		return op3(2, 0x34) | opf(0xC8)

	case AFLUSH:
		return op3(2, 0x3B)

	// Floating-point move.
	case AFMOVS:
		return op3(2, 0x34) | opf(1)
	case AFMOVD:
		return op3(2, 0x34) | opf(2)

	// Floating-point negate.
	case AFNEGS:
		return op3(2, 0x34) | opf(5)
	case AFNEGD:
		return op3(2, 0x34) | opf(6)

	// Floating-point square root.
	case AFSQRTS:
		return op3(2, 0x34) | opf(0x29)
	case AFSQRTD:
		return op3(2, 0x34) | opf(0x2A)

	// Convert floating-point to integer.
	case AFSTOX:
		return op3(2, 0x34) | opf(0x81)
	case AFDTOX:
		return op3(2, 0x34) | opf(0x82)
	case AFSTOI:
		return op3(2, 0x34) | opf(0xD1)
	case AFDTOI:
		return op3(2, 0x34) | opf(0xD2)

	// Convert between floating-point formats.
	case AFSTOD:
		return op3(2, 0x34) | opf(0xC9)
	case AFDTOS:
		return op3(2, 0x34) | opf(0xC6)

	// Convert 64-bit integer to floating point.
	case AFXTOS:
		return op3(2, 0x34) | opf(0x84)
	case AFXTOD:
		return op3(2, 0x34) | opf(0x88)

	// Jump and link.
	case AJMPL:
		return op3(2, 0x38)

	// Memory Barrier.
	case AMEMBAR:
		return op3(2, 0x28) | 0xF<<14 | 1<<13

	case ASETHI:
		return op2(4)

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
}

func pcrelclass(offset int64) int8 {
	if -1<<18 <= offset && offset <= 1<<18-1 {
		return ClassPCREL19
	}
	return ClassUnknown
}

func oregclass(offset int64) int8 {
	if offset == 0 {
		return ClassIndir0
	}
	if -4096 <= offset && offset <= 4095 {
		return ClassIndir13
	}
	return ClassIndir
}

func addrclass(offset int64) int8 {
	if -4096 <= offset && offset <= 4095 {
		return ClassEffectiveAddr13
	}
	return ClassEffectiveAddr
}

func constclass(offset int64) int8 {
	if 0 <= offset && offset <= 31 {
		return ClassConst5
	}
	if 0 <= offset && offset <= 63 {
		return ClassConst6
	}
	if -4096 <= offset && offset <= 4095 {
		return ClassConst13
	}
	if -1<<31 <= offset && offset < 0 {
		return ClassConst31_
	}
	if 0 <= offset && offset <= 1<<31-1 {
		return ClassConst31
	}
	if 0 <= offset && offset <= 1<<32-1 {
		return ClassConst32
	}
	return ClassConst
}

func rclass(r int16) int8 {
	switch {
	case r == RegZero:
		return ClassZero
	case REG_R1 <= r && r <= REG_R31:
		return ClassReg
	case REG_F0 <= r && r <= REG_F31:
		return ClassFloatReg
	case REG_D0 <= r && r <= REG_D62:
		return ClassDoubleReg
	case r == REG_BSP || r == REG_BFP:
		return ClassBiased
	case r == REG_ICC || r == REG_XCC:
		return ClassCondReg
	case REG_FCC0 <= r && r <= REG_FCC3:
		return ClassFloatCondReg
	case r >= REG_SPECIAL:
		return ClassSpecialReg
	}
	return ClassUnknown
}

func aclass(a *obj.Addr) int8 {
	switch a.Type {
	case obj.TYPE_NONE:
		return ClassNone

	case obj.TYPE_REG:
		return rclass(a.Reg)

	case obj.TYPE_REGREG:
		return ClassPairComma

	case obj.TYPE_MEM:
		switch a.Name {
		case obj.NAME_EXTERN, obj.NAME_STATIC:
			if a.Sym == nil {
				return ClassUnknown
			}
			return ClassMem

		case obj.NAME_AUTO, obj.NAME_PARAM:
			panic("unimplemented")

		case obj.TYPE_NONE:
			if a.Scale == 1 {
				return ClassPairPlus
			}
			return oregclass(a.Offset)
		}

	case obj.TYPE_FCONST:
		return ClassFloatConst

	case obj.TYPE_TEXTSIZE:
		return ClassTextSize

	case obj.TYPE_CONST, obj.TYPE_ADDR:
		switch a.Name {
		case obj.TYPE_NONE:
			if a.Reg != 0 {
				if a.Reg == RegZero && a.Offset == 0 {
					return ClassZero
				}
				return addrclass(a.Offset)
			}
			return constclass(a.Offset)

		case obj.NAME_EXTERN, obj.NAME_STATIC:
			if a.Sym == nil {
				return ClassUnknown
			}
			return ClassAddr

		case obj.NAME_AUTO, obj.NAME_PARAM:
			panic("unimplemented")
		}
	case obj.TYPE_BRANCH:
		if a.Sym == nil && a.Val != nil {
			return pcrelclass(a.Offset)
		}
	}
	return ClassUnknown
}

func span(ctxt *obj.Link, cursym *obj.LSym) {
	if cursym.Text == nil || cursym.Text.Link == nil { // handle external functions and ELF section symbols
		return
	}

	var pc int64      // relative to entry point
	var text []uint32 // actual assembled bytes
	for p := cursym.Text.Link; p != nil; p = p.Link {
		o, err := oplook(p)
		if err != nil {
			ctxt.Diag(err.Error())
		}
		out, err := asmout(p, o)
		if err != nil {
			ctxt.Diag(err.Error())
		}
		p.Pc = pc
		pc += int64(len(out)) * 4
		text = append(text, out...)
	}
	pc += -pc & (16 - 1)
	cursym.Size = pc
	obj.Symgrow(ctxt, cursym, pc)
	bp := cursym.P
	for _, v := range text {
		ctxt.Arch.ByteOrder.PutUint32(bp, v)
		bp = bp[4:]
	}
}

func asmout(p *obj.Prog, o int) (out []uint32, err error) {
	out = make([]uint32, 2)
	o1 := &out[0]
	o2 := &out[1]
	size := 1
	switch o {
	default:
		return nil, fmt.Errorf("unknown asm %d", o)

	// op Rs,       Rd	-> Rd = Rs op Rd
	// op Rs1, Rs2, Rd	-> Rd = Rs1 op Rs2
	case 1:
		reg := p.To.Reg
		if p.From3 != nil {
			reg = p.From3.Reg
		}
		*o1 = opalu(p.As) | rrr(p.From.Reg, 0, reg, p.To.Reg)

	// MOVD Rs, Rd
	case 2:
		*o1 = opalu(p.As) | rrr(p.From.Reg, 0, RegZero, p.To.Reg)

	// op Rs, $imm13, Rd	-> Rd = Rs op $imm13
	case 3:
		*o1 = opalu(p.As) | rsr(p.From.Reg, p.From3.Offset, p.To.Reg)

	// MOVD $imm13, Rd
	case 4:
		*o1 = opalu(p.As) | rsr(RegZero, p.From.Offset, p.To.Reg)

	// LDD (R1+R2), R	-> R = *(R1+R2)
	case 5:
		*o1 = opload(p.As) | rrr(p.From.Reg, 0, p.From.Index, p.To.Reg)

	// STD R, (R1+R2)	-> *(R1+R2) = R
	case 6:
		*o1 = opstore(p.As) | rrr(p.To.Reg, 0, p.To.Index, p.From.Reg)

	// LDD $imm13(Rs), R	-> R = *(Rs+$imm13)
	case 7:
		*o1 = opload(p.As) | rsr(p.From.Reg, p.From.Offset, p.To.Reg)

	// STD Rs, $imm13(R)	-> *(R+$imm13) = Rs
	case 8:
		*o1 = opstore(p.As) | rsr(p.To.Reg, p.To.Offset, p.From.Reg)

	// RD Rspecial, R
	case 9:
		*o1 = oprd(p.As) | uint32(p.From.Reg&0x1f)<<14 | rd(p.To.Reg)

	// CASD/CASW
	case 10:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.From3.Reg, p.To.Reg)

	// fop Fs, Fd
	case 11:
		*o1 = opcode(p.As) | rrr(0, 0, p.From.Reg, p.To.Reg)

	// SETHI $const, R
	case 12:
		if p.From.Offset&0x3FF != 0 {
			return nil, errors.New("SETHI constant not mod 1024")
		}
		*o1 = opcode(p.As) | ir(uint32(p.From.Offset)>>10, p.To.Reg)

	// MEMBAR $mask
	case 13:
		if p.From.Offset > 127 {
			return nil, errors.New("MEMBAR mask out of range")
		}
		*o1 = opcode(p.As) | uint32(p.From.Offset)

	// FCMPD FCC, F, F
	case 14:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.From3.Reg, p.To.Reg&3)

	// MOVD $imm32, R ->
	// 	SETHI hi($imm32), R
	// 	OR R, lo($imm32), R
	case 15:
		*o1 = opcode(ASETHI) | ir(uint32(p.From.Offset)>>10, p.To.Reg)
		if p.From.Offset&0x3FF == 0 {
			break
		}
		size = 2
		*o2 = opalu(AOR) | rsr(p.To.Reg, int64(p.From.Offset&0x3FF), p.To.Reg)

	// MOVD -$imm31, R ->
	// 	SETHI hi(^$imm32), R
	// 	XOR R, lo($imm32)|0x1C00, R
	case 16:
		size = 2
		*o1 = opcode(ASETHI) | ir(^(uint32(p.From.Offset))>>10, p.To.Reg)
		if p.From.Offset&0x3FF == 0 {
			*o2 = opalu(ASRAD) | rrr(p.To.Reg, 0, RegZero, p.To.Reg)
			break
		}
		*o2 = opalu(AXOR) | rsr(p.To.Reg, int64(uint32(p.From.Offset)&0x3ff|0x1C00), p.To.Reg)
	}

	return out[:size], nil
}

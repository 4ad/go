// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
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

	Optab{AADD, ClassConst13, ClassNone, ClassReg}:  2,
	Optab{AAND, ClassConst13, ClassNone, ClassReg}:  2,
	Optab{AMULD, ClassConst13, ClassNone, ClassReg}: 2,
	Optab{AADD, ClassConst13, ClassReg, ClassReg}:   2,
	Optab{AAND, ClassConst13, ClassReg, ClassReg}:   2,
	Optab{AMULD, ClassConst13, ClassReg, ClassReg}:  2,

	Optab{ALDD, ClassPairPlus, ClassNone, ClassReg}: 3,
	Optab{ASTD, ClassReg, ClassNone, ClassPairPlus}: 4,

	Optab{ALDD, ClassIndir13, ClassNone, ClassReg}: 5,
	Optab{ASTD, ClassReg, ClassNone, ClassIndir13}: 6,

	Optab{ARDPC, ClassNone, ClassNone, ClassReg}: 7,

	Optab{ACASD, ClassReg, ClassReg, ClassIndir0}: 8,
}

// Compatible classes, if something accepts a $hugeconst, it
// can also accept $smallconst, $0 and ZR. Something that accepts a
// register, can also accept $0, etc.
var cc = map[int8][]int8{
	ClassReg:           {ClassZero},
	ClassConst13:       {ClassZero},
	ClassConst:         {ClassConst13, ClassZero},
	ClassEffectiveAddr: {ClassEffectiveAddr13},
	ClassIndir13:       {ClassIndir0},
	ClassIndir:         {ClassIndir13, ClassIndir0},
}

// Compatible instructions, if an asm* function accepts AADD,
// it accepts ASUBCCC too.
var ci = map[int16][]int16{
	AADD:     {AADDCC, AADDC, AADDCCC, ASUB, ASUBCC, ASUBC, ASUBCCC},
	AAND:     {AANDCC, AANDN, AANDNCC, AOR, AORCC, AORN, AORNCC, AXOR, AXORCC, AXNOR, AXNORCC},
	obj.AJMP: {ABN, ABNE, ABE, ABG, ABLE, ABGE, ABL, ABGU, ABLEU, ABCC, ABCS, ABPOS, ABNEG, ABVC, ABVS},
	ABRZ:     {ABRLEZ, ABRLZ, ABRNZ, ABRGZ, ABRGEZ},
	ACASD:    {ACASW},
	AFABSD:   {AFABSS},
	AFADDD:   {AFADDS, AFSUBS, AFSUBD},
	AFBA:     {AFBN, AFBU, AFBG, AFBUG, AFBL, AFBUL, AFBLG, AFBNE, AFBE, AFBUE, AFBGE, AFBUGE, AFBLE, AFBULE, AFBO},
	AFCMPD:   {AFCMPS},
	AFDIVD:   {AFDIVS},
	AFWTOD:   {AFWTOS},
	AFMOVD:   {AFMOVS},
	AFMULD:   {AFMULS, AFSMULD},
	AFNEGD:   {AFNEGS},
	AFSQRTD:  {AFSQRTS},
	AFSTOXD:  {AFDTOXD, AFSTOXW, AFDTOXW},
	AFSTOD:   {AFDTOS},
	AFXTOD:   {AFXTOS},
	ALDD:     {ALDSB, ALDSH, ALDSW, ALDUB, ALDUH, ALDUW},
	ALDDF:    {ALDSF},
	AMULD:    {ASDIVD, AUDIVD},
	ARDPC:    {ARDTICK, ARDCCR},
	ASLLD:    {ASLLW, ASRLW, ASRAW, ASRLD, ASRAD},
	ASTD:     {ASTB, ASTH, ASTW},
	ASTDF:    {ASTSF},
}

func init() {
	// For each line in optab, duplicate it so that we'll also
	// have a line that will accept compatible instructions, but
	// only if there isn't an already existent line with the same
	// key.
	for o, v := range optab {
		for _, c := range ci[o.as] {
			do := o
			do.as = c
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
	o := Optab{as: p.As, a1: p.From.Class, a2: rclass(p.Reg), a3: p.To.Class}
	if p.Reg == 0 {
		o.a2 = ClassNone
	}
	v, ok := optab[o]
	if !ok {
		return 0, fmt.Errorf("illegal combination %v %v %v %v, %d %d", p, DRconv(o.a1), DRconv(o.a2), DRconv(o.a3), p.From.Type, p.To.Type)
	}
	return v, nil
}

func ir(imm22, rd int) uint32 {
	return uint32(rd&31<<25 | imm22&(1<<23-1))
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

func rrr(rs2, imm_asi, rs1, rd int16) uint32 {
	return uint32(rd&31<<25 | rs1&31<<14 | imm_asi&255<<5 | rs2&31)
}

func srr(simm13 int64, rs1, rd int16) uint32 {
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

func opcode(a int16) uint32 {
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

	case AFADDS:
		return op3(2, 0x34) | opf(0x41)
	case AFADDD:
		return op3(2, 0x34) | opf(0x42)

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

	// Floating-point divide.
	case AFDIVS:
		return op3(2, 0x34) | opf(0x4D)
	case AFDIVD:
		return op3(2, 0x34) | opf(0x4E)

	// Convert 32-bit integer to floating point.
	case AFWTOS:
		return op3(2, 0x34) | opf(0xC4)
	case AFWTOD:
		return op3(2, 0x34) | opf(0xC8)

	case AFLUSH:
		return op3(2, 0x3B)

	// Floating-point move.
	case AFMOVS:
		return op3(2, 0x34) | opf(1)
	case AFMOVD:
		return op3(2, 0x34) | opf(2)

	// Floating-point multiply.
	case AFMULS:
		return op3(2, 0x34) | opf(0x49)
	case AFMULD:
		return op3(2, 0x34) | opf(0x4A)
	case AFSMULD:
		return op3(2, 0x34) | opf(0x69)

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
	case AFSTOXD:
		return op3(2, 0x34) | opf(0x81)
	case AFDTOXD:
		return op3(2, 0x34) | opf(0x82)
	case AFSTOXW:
		return op3(2, 0x34) | opf(0xD1)
	case AFDTOXW:
		return op3(2, 0x34) | opf(0xD2)

	// Convert between floating-point formats.
	case AFSTOD:
		return op3(2, 0x34) | opf(0xC9)
	case AFDTOS:
		return op3(2, 0x34) | opf(0xC6)

	// Floating-point subtract.
	case AFSUBS:
		return op3(2, 0x34) | opf(0x45)
	case AFSUBD:
		return op3(2, 0x34) | opf(0x46)

	// Convert 64-bit integer to floating point.
	case AFXTOS:
		return op3(2, 0x34) | opf(0x84)
	case AFXTOD:
		return op3(2, 0x34) | opf(0x88)

	// Jump and link.
	case AJMPL:
		return op3(2, 0x38)

	// Load integer.
	case ALDSB:
		return op3(3, 9)
	case ALDSH:
		return op3(3, 10)
	case ALDSW:
		return op3(3, 8)
	case ALDUB:
		return op3(3, 1)
	case ALDUH:
		return op3(3, 2)
	case ALDUW:
		return op3(3, 0)
	case ALDD:
		return op3(3, 11)

	// Load floating-point register.
	case ALDSF:
		return op3(3, 0x20)
	case ALDDF:
		return op3(3, 0x23)

	// Memory Barrier.
	case AMEMBAR:
		return op3(2, 0x28) | 0xF<<14

	// Multiply and divide.
	case AMULD:
		return op3(2, 9)
	case ASDIVD:
		return op3(2, 0x2D)
	case AUDIVD:
		return op3(2, 0xD)

	// OR logical operation.
	case AOR:
		return op3(2, 2)
	case AORCC:
		return op3(2, 18)
	case AORN:
		return op3(2, 6)
	case AORNCC:
		return op3(2, 22)

	// Read ancillary state register.
	case ARDCCR:
		return op3(2, 0x28) | 2<<14
	case ARDTICK:
		return op3(2, 0x28) | 4<<14
	case ARDPC:
		return op3(2, 0x28) | 5<<14

	case ASETHI:
		return op2(4)

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

	// Store Integer.
	case ASTB:
		return op3(3, 5)
	case ASTH:
		return op3(3, 6)
	case ASTW:
		return op3(3, 4)
	case ASTD:
		return op3(3, 14)

	// Store floating-point.
	case ASTSF:
		return op3(3, 0x24)
	case ASTDF:
		return op3(3, 0x27)

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

	default:
		panic("unknown instruction: " + obj.Aconv(int(a)))
	}
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
	if -4096 <= offset && offset <= 4095 {
		return ClassConst13
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
	case r == REG_BSP || r == REG_BFP:
		return ClassBiased
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
		pc += int64(len(out)) * 4
		p.Pc = pc
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
	size := 1
	switch o {
	default:
		return nil, fmt.Errorf("unknown asm %d", o)

	// op Rs,        Rd	-> Rd = Rd  op Rs
	// op Rs1, Rs2,  Rd	-> Rd = Rs2 op Rs1
	case 1:
		reg := p.Reg
		if reg == 0 {
			reg = p.To.Reg
		}
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, reg, p.To.Reg)

	// op $imm13,     Rd	-> Rd = Rd op $imm13
	// op $imm13, Rs, Rd	-> Rd = Rs op $imm13
	case 2:
		reg := p.Reg
		if reg == 0 {
			reg = p.To.Reg
		}
		*o1 = opcode(p.As) | srr(p.From.Offset, reg, p.To.Reg)

	// LDD (R1+R2), R	-> R = *(R1+R2)
	case 3:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.From.Index, p.To.Reg)

	// STD R, (R1+R2)	-> *(R1+R2) = R
	case 4:
		*o1 = opcode(p.As) | rrr(p.To.Reg, 0, p.To.Index, p.From.Reg)

	// LDD $imm13(Rs), R	-> R = *(Rs+$imm13)
	case 5:
		*o1 = opcode(p.As) | srr(p.From.Offset, p.From.Reg, p.To.Reg)

	// STD Rs, $imm13(R)	-> *(R+$imm13) = Rs
	case 6:
		*o1 = opcode(p.As) | srr(p.To.Offset, p.To.Reg, p.From.Reg)

	// RDPC R
	case 7:
		*o1 = opcode(p.As) | rd(p.To.Reg)

	// CASD/CASW
	case 8:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.To.Reg, p.Reg)
	}

	return out[:size], nil
}

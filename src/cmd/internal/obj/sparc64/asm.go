// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import "cmd/internal/obj"

func ri(rd, imm22 int) uint32 {
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

func rrr(rd, rs1, imm_asi, rs2 int) uint32 {
	return uint32(rd&31<<25 | rs1&31<<14 | imm_asi&255<<5 | rs2&31)
}

func rrs(rd, rs1, simm13 int) uint32 {
	return uint32(rd&31<<25 | rs1&31<<14 | 1<<13 | simm13&(1<<14-1))
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

func opcode(a int) uint32 {
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

	case ACASAW:
		return op3(3, 0x3C)
	case ACASA:
		return op3(3, 0x3E)

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
		return op3(2, 0x35) | opf(0x31)
	case AFCMPD:
		return op3(2, 0x35) | opf(0x32)

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
		panic("unknown instruction: " + obj.Aconv(a))
	}
}

func rclass(r int16) int {
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

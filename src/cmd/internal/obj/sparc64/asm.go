// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"errors"
	"fmt"
	"sort"
)

type Optab struct {
	as int16 // instruction
	a1 int8  // from
	a2 int8  // reg
	a3 int8  // from3
	a4 int8  // to
}

type OptabSlice []Optab

func (tab OptabSlice) Len() int { return len(tab) }

func (tab OptabSlice) Swap(i, j int) { tab[i], tab[j] = tab[j], tab[i] }

func (tab OptabSlice) Less(i, j int) bool {
	return ocmp(tab[i], tab[j])
}

func ocmp(o1, o2 Optab) bool {
	if o1.as != o2.as {
		return o1.as < o2.as
	}
	if o1.a1 != o2.a1 {
		return o1.a1 < o2.a1
	}
	if o1.a2 != o2.a2 {
		return o1.a2 < o2.a2
	}
	if o1.a3 != o2.a3 {
		return o1.a3 < o2.a3
	}
	return o1.a4 < o2.a4
}

type Opval struct {
	op   int8 // selects case in asmout switch
	size int8 // *not* including delay-slot
}

var optab = map[Optab]Opval{
	Optab{obj.ATEXT, ClassAddr, ClassNone, ClassNone, ClassTextSize}: {0, 0},
	Optab{obj.AFUNCDATA, ClassConst, ClassNone, ClassNone, ClassMem}: {0, 0},
	Optab{obj.APCDATA, ClassConst, ClassNone, ClassNone, ClassConst}: {0, 0},

	Optab{AADD, ClassReg, ClassNone, ClassNone, ClassReg}:  {1, 4},
	Optab{AAND, ClassReg, ClassNone, ClassNone, ClassReg}:  {1, 4},
	Optab{AMULD, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4},
	Optab{ASLLD, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4},
	Optab{ASLLW, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4},
	Optab{AADD, ClassReg, ClassReg, ClassNone, ClassReg}:   {1, 4},
	Optab{AAND, ClassReg, ClassReg, ClassNone, ClassReg}:   {1, 4},
	Optab{AMULD, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4},
	Optab{ASLLD, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4},
	Optab{ASLLW, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4},

	Optab{AFADDD, ClassDReg, ClassNone, ClassNone, ClassDReg}:  {1, 4},
	Optab{AFADDD, ClassDReg, ClassDReg, ClassNone, ClassDReg}:  {1, 4},
	Optab{AFSMULD, ClassFReg, ClassFReg, ClassNone, ClassDReg}: {1, 4},

	Optab{AMOVD, ClassReg, ClassNone, ClassNone, ClassReg}: {2, 4},

	Optab{AADD, ClassConst13, ClassNone, ClassNone, ClassReg}:  {3, 4},
	Optab{AAND, ClassConst13, ClassNone, ClassNone, ClassReg}:  {3, 4},
	Optab{AMULD, ClassConst13, ClassNone, ClassNone, ClassReg}: {3, 4},
	Optab{ASLLD, ClassConst6, ClassNone, ClassNone, ClassReg}:  {3, 4},
	Optab{ASLLW, ClassConst5, ClassNone, ClassNone, ClassReg}:  {3, 4},
	Optab{AADD, ClassConst13, ClassReg, ClassNone, ClassReg}:   {3, 4},
	Optab{AAND, ClassConst13, ClassReg, ClassNone, ClassReg}:   {3, 4},
	Optab{AMULD, ClassConst13, ClassReg, ClassNone, ClassReg}:  {3, 4},
	Optab{ASLLD, ClassConst6, ClassReg, ClassNone, ClassReg}:   {3, 4},
	Optab{ASLLW, ClassConst5, ClassReg, ClassNone, ClassReg}:   {3, 4},

	Optab{AMOVD, ClassConst13, ClassNone, ClassNone, ClassReg}: {4, 4},

	Optab{ALDD, ClassIndirRegReg, ClassNone, ClassNone, ClassReg}:   {5, 4},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndirRegReg}:   {6, 4},
	Optab{ALDDF, ClassIndirRegReg, ClassNone, ClassNone, ClassDReg}: {5, 4},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassIndirRegReg}: {6, 4},

	Optab{ALDD, ClassIndir13, ClassNone, ClassNone, ClassReg}:   {7, 4},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndir13}:   {8, 4},
	Optab{ALDDF, ClassIndir13, ClassNone, ClassNone, ClassDReg}: {7, 4},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassIndir13}: {8, 4},

	Optab{ARD, ClassSpcReg, ClassNone, ClassNone, ClassReg}: {9, 4},

	Optab{ACASD, ClassIndir0, ClassReg, ClassNone, ClassReg}: {10, 4},

	Optab{AFSTOD, ClassFReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFDTOS, ClassDReg, ClassNone, ClassNone, ClassFReg}: {11, 4},

	Optab{AFMOVD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},

	Optab{AFXTOD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFITOD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFXTOS, ClassDReg, ClassNone, ClassNone, ClassFReg}: {11, 4},
	Optab{AFITOS, ClassFReg, ClassNone, ClassNone, ClassFReg}: {11, 4},

	Optab{AFSTOX, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFDTOX, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFDTOI, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},
	Optab{AFSTOI, ClassFReg, ClassNone, ClassNone, ClassFReg}: {11, 4},

	Optab{AFABSD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4},

	Optab{ASETHI, ClassConst32, ClassNone, ClassNone, ClassReg}: {12, 4},
	Optab{ARNOP, ClassNone, ClassNone, ClassNone, ClassNone}:    {12, 4},

	Optab{AMEMBAR, ClassConst, ClassNone, ClassNone, ClassNone}: {13, 4},

	Optab{AFCMPD, ClassDReg, ClassDReg, ClassNone, ClassFCond}: {14, 4},
	Optab{AFCMPD, ClassDReg, ClassDReg, ClassNone, ClassNone}:  {14, 4},

	Optab{AMOVD, ClassConst32, ClassNone, ClassNone, ClassReg}:  {15, 8},
	Optab{AMOVD, ClassConst31_, ClassNone, ClassNone, ClassReg}: {16, 8},

	Optab{obj.AJMP, ClassNone, ClassNone, ClassNone, ClassBranch}: {17, 4},
	Optab{ABN, ClassCond, ClassNone, ClassNone, ClassBranch}:      {17, 4},
	Optab{ABNW, ClassNone, ClassNone, ClassNone, ClassBranch}:     {17, 4},
	Optab{ABRZ, ClassReg, ClassNone, ClassNone, ClassBranch}:      {18, 4},
	Optab{AFBA, ClassNone, ClassNone, ClassNone, ClassBranch}:     {19, 4},

	Optab{AJMPL, ClassReg, ClassNone, ClassNone, ClassReg}:        {20, 4},
	Optab{AJMPL, ClassRegConst13, ClassNone, ClassNone, ClassReg}: {20, 4},
	Optab{AJMPL, ClassRegReg, ClassNone, ClassNone, ClassReg}:     {21, 4},

	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassMem}:     {22, 4},
	Optab{obj.ADUFFZERO, ClassNone, ClassNone, ClassNone, ClassMem}: {22, 4},
	Optab{obj.ADUFFCOPY, ClassNone, ClassNone, ClassNone, ClassMem}: {22, 4},

	Optab{AMOVD, ClassAddr, ClassNone, ClassNone, ClassReg}: {23, 8},

	Optab{ALDD, ClassMem, ClassNone, ClassNone, ClassReg}:   {24, 12},
	Optab{ALDDF, ClassMem, ClassNone, ClassNone, ClassDReg}: {24, 12},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassMem}:   {25, 12},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassMem}: {25, 12},

	Optab{obj.ARET, ClassNone, ClassNone, ClassNone, ClassNone}: {26, 4},

	Optab{ATA, ClassConst13, ClassNone, ClassNone, ClassNone}: {27, 4},

	Optab{AMOVD, ClassRegConst13, ClassNone, ClassNone, ClassReg}: {28, 4},

	Optab{AMOVUB, ClassReg, ClassNone, ClassNone, ClassReg}: {29, 4},
	Optab{AMOVUH, ClassReg, ClassNone, ClassNone, ClassReg}: {30, 8},
	Optab{AMOVUW, ClassReg, ClassNone, ClassNone, ClassReg}: {31, 4},

	Optab{AMOVB, ClassReg, ClassNone, ClassNone, ClassReg}: {32, 8},
	Optab{AMOVH, ClassReg, ClassNone, ClassNone, ClassReg}: {33, 8},
	Optab{AMOVW, ClassReg, ClassNone, ClassNone, ClassReg}: {34, 4},

	Optab{ANEG, ClassReg, ClassNone, ClassNone, ClassReg}: {35, 4},

	Optab{ACMP, ClassReg, ClassReg, ClassNone, ClassNone}:     {36, 4},
	Optab{ACMP, ClassConst13, ClassReg, ClassNone, ClassNone}: {37, 4},

	Optab{ABND, ClassNone, ClassNone, ClassNone, ClassBranch}: {38, 4},

	Optab{obj.AUNDEF, ClassNone, ClassNone, ClassNone, ClassNone}: {39, 4},

	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassReg}:    {40, 4},
	Optab{obj.ACALL, ClassReg, ClassNone, ClassNone, ClassReg}:     {40, 4},
	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassIndir0}: {40, 4},
	Optab{obj.ACALL, ClassReg, ClassNone, ClassNone, ClassIndir0}:  {40, 4},

	Optab{AADD, ClassConst32, ClassNone, ClassNone, ClassReg}: {41, 12},
	Optab{AAND, ClassConst32, ClassNone, ClassNone, ClassReg}: {41, 12},
	Optab{AADD, ClassConst32, ClassReg, ClassNone, ClassReg}:  {41, 12},
	Optab{AAND, ClassConst32, ClassReg, ClassNone, ClassReg}:  {41, 12},

	Optab{AMOVD, ClassRegConst, ClassNone, ClassNone, ClassReg}: {42, 12},

	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndir}: {43, 12},
	Optab{ALDD, ClassIndir, ClassNone, ClassNone, ClassReg}: {44, 12},
}

// Compatible classes, if something accepts a $hugeconst, it
// can also accept $smallconst, $0 and ZR. Something that accepts a
// register, can also accept $0, etc.
var cc = map[int8][]int8{
	ClassReg:      {ClassZero},
	ClassConst6:   {ClassConst5, ClassZero},
	ClassConst13:  {ClassConst6, ClassConst5, ClassZero},
	ClassConst31:  {ClassConst6, ClassConst5, ClassZero},
	ClassConst32:  {ClassConst31_, ClassConst31, ClassConst13, ClassConst6, ClassConst5, ClassZero},
	ClassConst:    {ClassConst32, ClassConst31_, ClassConst31, ClassConst13, ClassConst6, ClassConst5, ClassZero},
	ClassRegConst: {ClassRegConst13},
	ClassIndir13:  {ClassIndir0},
	ClassIndir:    {ClassIndir13, ClassIndir0},
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
	AADD:   {AADDCC, AADDC, AADDCCC, ASUB, ASUBCC, ASUBC, ASUBCCC},
	AAND:   {AANDCC, AANDN, AANDNCC, AOR, AORCC, AORN, AORNCC, AXOR, AXORCC, AXNOR, AXNORCC},
	ABN:    {ABNE, ABE, ABG, ABLE, ABGE, ABL, ABGU, ABLEU, ABCC, ABCS, ABPOS, ABNEG, ABVC, ABVS},
	ABNW:   {ABNEW, ABEW, ABGW, ABLEW, ABGEW, ABLW, ABGUW, ABLEUW, ABCCW, ABCSW, ABPOSW, ABNEGW, ABVCW, ABVSW},
	ABND:   {ABNED, ABED, ABGD, ABLED, ABGED, ABLD, ABGUD, ABLEUD, ABCCD, ABCSD, ABPOSD, ABNEGD, ABVCD, ABVSD},
	ABRZ:   {ABRLEZ, ABRLZ, ABRNZ, ABRGZ, ABRGEZ},
	ACASD:  {ACASW},
	AFABSD: {AFABSS, AFNEGD, AFNEGS, AFSQRTD, AFNEGS},
	AFADDD: {AFADDS, AFSUBS, AFSUBD, AFMULD, AFMULS, AFSMULD, AFDIVD, AFDIVS},
	AFBA:   {AFBN, AFBU, AFBG, AFBUG, AFBL, AFBUL, AFBLG, AFBNE, AFBE, AFBUE, AFBGE, AFBUGE, AFBLE, AFBULE, AFBO},
	AFCMPD: {AFCMPS},
	AFITOD: {AFITOS},
	AFMOVD: {AFMOVS},
	AFSTOD: {AFDTOS},
	AFXTOD: {AFXTOS},
	ALDD:   {ALDSB, ALDSH, ALDSW, ALDUB, ALDUH, ALDUW, AMOVB, AMOVH, AMOVW, AMOVUB, AMOVUH, AMOVUW, AMOVD},
	ALDDF:  {ALDSF, AFMOVD, AFMOVS},
	AMULD:  {ASDIVD, AUDIVD},
	ARD:    {AMOVD},
	ASLLD:  {ASRLD, ASRAD},
	ASLLW:  {ASLLW, ASRLW, ASRAW},
	ASTD:   {ASTB, ASTH, ASTW, AMOVB, AMOVH, AMOVW, AMOVUB, AMOVUH, AMOVUW, AMOVD},
	ASTDF:  {ASTSF, AFMOVD, AFMOVS},
}

func opkeys() OptabSlice {
	keys := make(OptabSlice, 0, len(optab))
	// create sorted map index by keys
	for k := range optab {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func init() {
	// For each line in optab, duplicate it so that we'll also
	// have a line that will accept compatible instructions, but
	// only if there isn't an already existent line with the same
	// key. Also change operand type, if the instruction is a double.
	for _, o := range opkeys() {
		for _, c := range ci[o.as] {
			do := o
			do.as = c
			if isInstDouble[o.as] && isInstFloat[do.as] {
				if do.a1 == ClassDReg {
					do.a1 = ClassFReg
				}
				if do.a2 == ClassDReg {
					do.a2 = ClassFReg
				}
				if do.a3 == ClassDReg {
					do.a3 = ClassFReg
				}
				if do.a4 == ClassDReg {
					do.a4 = ClassFReg
				}
			}
			_, ok := optab[do]
			if !ok {
				optab[do] = optab[o]
			}
		}
	}
	// For each line in optab that accepts a large-class operand,
	// duplicate it so that we'll also have a line that accepts a
	// small-class operand, but do it only if there isn't an already
	// existent line with the same key.
	for _, o := range opkeys() {
		for _, c := range cc[o.a1] {
			do := o
			do.a1 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = optab[o]
			}
		}
	}
	for _, o := range opkeys() {
		for _, c := range cc[o.a2] {
			do := o
			do.a2 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = optab[o]
			}
		}
	}
	for _, o := range opkeys() {
		for _, c := range cc[o.a3] {
			do := o
			do.a3 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = optab[o]
			}
		}
	}
	for _, o := range opkeys() {
		for _, c := range cc[o.a4] {
			do := o
			do.a4 = c
			_, ok := optab[do]
			if !ok {
				optab[do] = optab[o]
			}
		}
	}
}

func oplook(p *obj.Prog) (Opval, error) {
	var a2, a3 int8 = ClassNone, ClassNone
	if p.Reg != 0 {
		a2 = rclass(p.Reg)
	}
	var type3 int16
	if p.From3 != nil {
		a3 = p.From3.Class
		type3 = p.From3.Type
	}
	o := Optab{as: p.As, a1: p.From.Class, a2: a2, a3: a3, a4: p.To.Class}
	v, ok := optab[o]
	if !ok {
		return Opval{}, fmt.Errorf("illegal combination %v %v %v %v %v, %d %d %d %d", p, DRconv(o.a1), DRconv(o.a2), DRconv(o.a3), DRconv(o.a4), p.From.Type, p.Reg, type3, p.To.Type)
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
	case ALDSB, AMOVB:
		return op3(3, 9)
	case ALDSH, AMOVH:
		return op3(3, 10)
	case ALDSW, AMOVW:
		return op3(3, 8)
	case ALDUB, AMOVUB:
		return op3(3, 1)
	case ALDUH, AMOVUH:
		return op3(3, 2)
	case ALDUW, AMOVUW:
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
	case ASTB, AMOVB, AMOVUB:
		return op3(3, 5)
	case ASTH, AMOVH, AMOVUH:
		return op3(3, 6)
	case ASTW, AMOVW, AMOVUW:
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
	case ABN, ABNW, ABND:
		return cond(0) | op2(1)
	case ABNE, ABNEW, ABNED:
		return cond(9) | op2(1)
	case ABE, ABEW, ABED:
		return cond(1) | op2(1)
	case ABG, ABGW, ABGD:
		return cond(10) | op2(1)
	case ABLE, ABLEW, ABLED:
		return cond(2) | op2(1)
	case ABGE, ABGEW, ABGED:
		return cond(11) | op2(1)
	case ABL, ABLW, ABLD:
		return cond(3) | op2(1)
	case ABGU, ABGUW, ABGUD:
		return cond(12) | op2(1)
	case ABLEU, ABLEUW, ABLEUD:
		return cond(4) | op2(1)
	case ABCC, ABCCW, ABCCD:
		return cond(13) | op2(1)
	case ABCS, ABCSW, ABCSD:
		return cond(5) | op2(1)
	case ABPOS, ABPOSW, ABPOSD:
		return cond(14) | op2(1)
	case ABNEG, ABNEGW, ABNEGD:
		return cond(6) | op2(1)
	case ABVC, ABVCW, ABVCD:
		return cond(15) | op2(1)
	case ABVS, ABVSW, ABVSD:
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
	case obj.ACALL, obj.ADUFFCOPY, obj.ADUFFZERO:
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

	case ASETHI, ARNOP:
		return op2(4)

	// Trap on Integer Condition Codes (Tcc).
	case ATA:
		return op3(2, 0x3A)

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
		return ClassRegConst13
	}
	return ClassRegConst
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
	case r == REG_ZR:
		return ClassZero
	case REG_R1 <= r && r <= REG_R31:
		return ClassReg
	case REG_F0 <= r && r <= REG_F31:
		return ClassFReg
	case REG_D0 <= r && r <= REG_D62:
		return ClassDReg
	case r == REG_ICC || r == REG_XCC:
		return ClassCond
	case REG_FCC0 <= r && r <= REG_FCC3:
		return ClassFCond
	case r >= REG_SPECIAL:
		return ClassSpcReg
	}
	return ClassUnknown
}

func aclass(a *obj.Addr) int8 {
	switch a.Type {
	case obj.TYPE_NONE:
		return ClassNone

	case obj.TYPE_REG:
		return rclass(a.Reg)

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
				return ClassIndirRegReg
			}
			return oregclass(a.Offset)
		}

	case obj.TYPE_FCONST:
		return ClassFConst

	case obj.TYPE_TEXTSIZE:
		return ClassTextSize

	case obj.TYPE_CONST, obj.TYPE_ADDR:
		switch a.Name {
		case obj.TYPE_NONE:
			if a.Reg != 0 {
				if a.Reg == REG_ZR && a.Offset == 0 {
					return ClassZero
				}
				if a.Scale == 1 {
					return ClassRegReg
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
		return ClassBranch
	}
	return ClassUnknown
}

func span(ctxt *obj.Link, cursym *obj.LSym) {
	if cursym.Text == nil || cursym.Text.Link == nil { // handle external functions and ELF section symbols
		return
	}

	var pc int64 // relative to entry point
	for p := cursym.Text.Link; p != nil; p = p.Link {
		o, err := oplook(p)
		if err != nil {
			ctxt.Diag(err.Error())
		}
		p.Pc = pc
		pc += int64(o.size)
	}
	pc += -pc & (16 - 1)
	cursym.Size = pc
	obj.Symgrow(ctxt, cursym, pc)

	var text []uint32 // actual assembled bytes
	for p := cursym.Text.Link; p != nil; p = p.Link {
		o, _ := oplook(p)
		out, _ := asmout(p, o, cursym)
		text = append(text, out...)
	}

	bp := cursym.P
	for _, v := range text {
		ctxt.Arch.ByteOrder.PutUint32(bp, v)
		bp = bp[4:]
	}
}

// bigmove assembles a move of the constant part of addr into reg.
func bigmove(addr *obj.Addr, reg int16) (out []uint32) {
	out = make([]uint32, 2)
	class := aclass(addr)
	switch class {
	case ClassRegConst, ClassIndir:
		class = constclass(addr.Offset)
	}
	switch class {
	// MOVD $imm32, R ->
	// 	SETHI hi($imm32), R
	// 	OR R, lo($imm32), R
	case ClassConst31, ClassConst32:
		out[0] = opcode(ASETHI) | ir(uint32(addr.Offset)>>10, reg)
		out[1] = opalu(AOR) | rsr(reg, int64(addr.Offset&0x3FF), reg)

	// MOVD -$imm31, R ->
	// 	SETHI hi(^$imm32), R
	// 	XOR R, lo($imm32)|0x1C00, R
	case ClassConst31_:
		out[0] = opcode(ASETHI) | ir(^(uint32(addr.Offset))>>10, reg)
		if addr.Offset&0x3FF == 0 {
			out[1] = opalu(ASRAD) | rrr(reg, 0, REG_ZR, reg)
			return out
		}
		out[1] = opalu(AXOR) | rsr(reg, int64(uint32(addr.Offset)&0x3ff|0x1C00), reg)
	default:
		panic("unexpected operand class: " + DRconv(class))
	}
	return out
}

func asmout(p *obj.Prog, o Opval, cursym *obj.LSym) (out []uint32, err error) {
	out = make([]uint32, 3)
	o1 := &out[0]
	o2 := &out[1]
	o3 := &out[2]
	switch o.op {
	default:
		return nil, fmt.Errorf("unknown asm %d", o)

	// op Rs,       Rd	-> Rd = Rs op Rd
	// op Rs1, Rs2, Rd	-> Rd = Rs2 op Rs1
	case 1:
		reg := p.To.Reg
		if p.Reg != 0 {
			reg = p.Reg
		}
		*o1 = opalu(p.As) | rrr(reg, 0, p.From.Reg, p.To.Reg)

	// MOVD Rs, Rd
	case 2:
		*o1 = opalu(p.As) | rrr(REG_ZR, 0, p.From.Reg, p.To.Reg)

	// op $imm13, Rs, Rd	-> Rd = Rs op $imm13
	case 3:
		reg := p.To.Reg
		if p.Reg != 0 {
			reg = p.Reg
		}
		*o1 = opalu(p.As) | rsr(reg, p.From.Offset, p.To.Reg)

	// MOVD $imm13, Rd
	case 4:
		*o1 = opalu(p.As) | rsr(REG_ZR, p.From.Offset, p.To.Reg)

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
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.Reg, p.To.Reg)

	// fop Fs, Fd
	case 11:
		*o1 = opcode(p.As) | rrr(0, 0, p.From.Reg, p.To.Reg)

	// SETHI $const, R
	// RNOP
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

	// FCMPD F, F, FCC
	case 14:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.Reg, p.To.Reg&3)

	// MOVD $imm32, R
	// MOVD -$imm31, R
	case 15, 16:
		out := bigmove(&p.From, p.To.Reg)
		return out, nil

	// BLE XCC, n(PC)
	// JMP n(PC)
	case 17:
		offset := p.Pcond.Pc - p.Pc
		if offset < -1<<22 || offset > 1<<22-1 {
			return nil, errors.New("branch target out of range")
		}
		if offset%4 != 0 {
			return nil, errors.New("branch target not mod 4")
		}
		*o1 = opcode(p.As) | uint32(p.From.Reg&3)<<20 | uint32(offset>>2)&(1<<19-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}

	// BRZ R, n(PC)
	case 18:
		offset := p.Pcond.Pc - p.Pc
		if offset < -1<<19 || offset > 1<<19-1 {
			return nil, errors.New("branch target out of range")
		}
		if offset%4 != 0 {
			return nil, errors.New("branch target not mod 4")
		}
		*o1 = opcode(p.As) | uint32((offset>>14)&3)<<20 | uint32(p.From.Reg&31)<<14 | uint32(offset>>2)&(1<<14-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}

	// FBA n(PC)
	case 19:
		offset := p.Pcond.Pc - p.Pc
		if offset < -1<<25 || offset > 1<<25-1 {
			return nil, errors.New("branch target out of range")
		}
		if offset%4 != 0 {
			return nil, errors.New("branch target not mod 4")
		}
		*o1 = opcode(p.As) | uint32(offset>>2)&(1<<22-1)

	// JMPL $imm13(Rs1), Rd
	case 20:
		*o1 = opcode(p.As) | rsr(p.From.Reg, p.From.Offset, p.To.Reg)

	// JMPL $(R1+R2), Rd
	case 21:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.From.Index, p.To.Reg)

	// CALL sym(SB)
	// DUFFCOPY, DUFFZERO
	case 22:
		*o1 = opcode(p.As)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 4
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_CALLSPARC64

	// MOVD $sym(SB), R ->
	// 	SETHI hi($sym), R
	// 	OR R, lo($sym), R
	case 23:
		*o1 = opcode(ASETHI) | ir(0, p.To.Reg)
		*o2 = opalu(AOR) | rsr(p.To.Reg, 0, p.To.Reg)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64

	// MOV sym(SB), R ->
	// 	SETHI hi($sym), R
	// 	OR R, lo($sym), R
	//	MOV (R), R
	case 24:
		*o1 = opcode(ASETHI) | ir(0, p.To.Reg)
		*o2 = opalu(AOR) | rsr(p.To.Reg, 0, p.To.Reg)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64
		*o3 = opload(p.As) | rsr(p.To.Reg, 0, p.To.Reg)

	// MOV R, sym(SB) ->
	// 	SETHI hi($sym), TMP
	// 	OR R, lo($sym), TMP
	//	MOV R, (TMP)
	case 25:
		*o1 = opcode(ASETHI) | ir(0, REG_TMP)
		*o2 = opalu(AOR) | rsr(REG_TMP, 0, REG_TMP)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_ADDRSPARC64
		*o3 = opstore(p.As) | rsr(REG_TMP, 0, p.From.Reg)

	// RET
	case 26:
		*o1 = opcode(AJMPL) | rsr(REG_LR, 8, REG_ZR)

	// TA $tn
	case 27:
		if p.From.Offset > 255 {
			return nil, errors.New("trap number too big")
		}
		*o1 = cond(8) | opcode(p.As) | 1<<13 | uint32(p.From.Offset&0xff)

	// MOVD	$imm13(R), Rd -> ADD R, $imm13, Rd
	case 28:
		*o1 = opalu(AADD) | rsr(p.From.Reg, p.From.Offset, p.To.Reg)

	// MOVUB Rs, Rd
	case 29:
		*o1 = opalu(AAND) | rsr(p.From.Reg, 0xff, p.To.Reg)

	// AMOVUH Rs, Rd
	case 30:
		*o1 = opalu(ASLLD) | rsr(p.From.Reg, 48, p.To.Reg)
		*o2 = opalu(ASRLD) | rsr(p.To.Reg, 48, p.To.Reg)

	// AMOVUW Rs, Rd
	case 31:
		*o1 = opalu(ASRLW) | rsr(p.From.Reg, 0, p.To.Reg)

	// AMOVB Rs, Rd
	case 32:
		*o1 = opalu(ASLLD) | rsr(p.From.Reg, 56, p.To.Reg)
		*o2 = opalu(ASRAD) | rsr(p.To.Reg, 56, p.To.Reg)

	// AMOVH Rs, Rd
	case 33:
		*o1 = opalu(ASLLD) | rsr(p.From.Reg, 48, p.To.Reg)
		*o2 = opalu(ASRAD) | rsr(p.To.Reg, 48, p.To.Reg)

	// AMOVW Rs, Rd
	case 34:
		*o1 = opalu(ASRAW) | rsr(p.From.Reg, 0, p.To.Reg)

	// ANEG Rs, Rd
	case 35:
		*o1 = opalu(ASUB) | rrr(REG_ZR, 0, p.From.Reg, p.To.Reg)

	// CMP R1, R2
	case 36:
		*o1 = opalu(ASUBCC) | rrr(p.Reg, 0, p.From.Reg, REG_ZR)

	// CMP $42, R2
	case 37:
		*o1 = opalu(ASUBCC) | rsr(p.Reg, p.From.Offset, REG_ZR)

	// BLED, n(PC)
	// JMP n(PC)
	case 38:
		offset := p.Pcond.Pc - p.Pc
		if offset < -1<<22 || offset > 1<<22-1 {
			return nil, errors.New("branch target out of range")
		}
		if offset%4 != 0 {
			return nil, errors.New("branch target not mod 4")
		}
		*o1 = opcode(p.As) | 2<<20 | uint32(offset>>2)&(1<<19-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}

	// UNDEF
	// This is supposed to be something that stops execution.
	// It's not supposed to be reached, ever, but if it is, we'd
	// like to be able to tell how we got there.  Assemble as
	// 0xdead0 which is guaranteed to raise undefined instruction
	// exception.
	case 39:
		*o1 = 0xdead0 // ILLTRAP

	// CALL R
	// CALL (R)
	// CALL R, R
	case 40:
		*o1 = opcode(AJMPL) | rsr(p.To.Reg, 0, REG_LR)

	// ADD $huge, Rd
	// AND $huge, Rs, Rd
	case 41:
		move := bigmove(&p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		reg := p.To.Reg
		if p.Reg != 0 {
			reg = p.Reg
		}
		*o3 = opalu(p.As) | rrr(reg, 0, REG_TMP, p.To.Reg)

	// AMOVD $huge(R), R
	case 42:
		move := bigmove(&p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opalu(AADD) | rrr(p.From.Reg, 0, REG_TMP, p.To.Reg)

	// AMOVD R, huge(R)
	case 43:
		move := bigmove(&p.To, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opstore(p.As) | rrr(p.To.Reg, 0, REG_TMP, p.From.Reg)

	// AMOVD huge(R), R
	case 44:
		move := bigmove(&p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opload(p.As) | rrr(p.From.Reg, 0, REG_TMP, p.To.Reg)
	}

	return out[:o.size/4], nil
}

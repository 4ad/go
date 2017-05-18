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
	as obj.As // instruction
	a1 int8   // from
	a2 int8   // reg
	a3 int8   // from3
	a4 int8   // to
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
	op     int8 // selects case in asmout switch
	size   int8 // *not* including delay-slot
	OpInfo      // information about the instruction
}

type OpInfo int8

const (
	ClobberTMP OpInfo = 1 << iota
)

var optab = map[Optab]Opval{
	Optab{obj.ATEXT, ClassAddr, ClassNone, ClassNone, ClassTextSize}: {0, 0, 0},
	Optab{obj.AFUNCDATA, ClassConst, ClassNone, ClassNone, ClassMem}: {0, 0, 0},
	Optab{obj.APCDATA, ClassConst, ClassNone, ClassNone, ClassConst}: {0, 0, 0},

	Optab{AADD, ClassReg, ClassNone, ClassNone, ClassReg}:  {1, 4, 0},
	Optab{AAND, ClassReg, ClassNone, ClassNone, ClassReg}:  {1, 4, 0},
	Optab{AMULD, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4, 0},
	Optab{ASLLD, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4, 0},
	Optab{ASLLW, ClassReg, ClassNone, ClassNone, ClassReg}: {1, 4, 0},
	Optab{AADD, ClassReg, ClassReg, ClassNone, ClassReg}:   {1, 4, 0},
	Optab{AAND, ClassReg, ClassReg, ClassNone, ClassReg}:   {1, 4, 0},
	Optab{AMULD, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4, 0},
	Optab{ASLLD, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4, 0},
	Optab{ASLLW, ClassReg, ClassReg, ClassNone, ClassReg}:  {1, 4, 0},

	Optab{ASAVE, ClassReg, ClassReg, ClassNone, ClassReg}: {1, 4, 0},
	Optab{ASAVE, ClassReg, ClassReg | ClassBias, ClassNone, ClassReg | ClassBias}: {1, 4, 0},

	Optab{AFADDD, ClassDReg, ClassNone, ClassNone, ClassDReg}:  {1, 4, 0},
	Optab{AFADDD, ClassDReg, ClassDReg, ClassNone, ClassDReg}:  {1, 4, 0},
	Optab{AFSMULD, ClassFReg, ClassFReg, ClassNone, ClassDReg}: {1, 4, 0},

	Optab{AMOVD, ClassReg, ClassNone, ClassNone, ClassReg}: {2, 4, 0},

	Optab{AADD, ClassConst13, ClassNone, ClassNone, ClassReg}:  {3, 4, 0},
	Optab{AAND, ClassConst13, ClassNone, ClassNone, ClassReg}:  {3, 4, 0},
	Optab{AMULD, ClassConst13, ClassNone, ClassNone, ClassReg}: {3, 4, 0},
	Optab{ASLLD, ClassConst6, ClassNone, ClassNone, ClassReg}:  {3, 4, 0},
	Optab{ASLLW, ClassConst5, ClassNone, ClassNone, ClassReg}:  {3, 4, 0},
	Optab{AADD, ClassConst13, ClassReg, ClassNone, ClassReg}:   {3, 4, 0},
	Optab{AAND, ClassConst13, ClassReg, ClassNone, ClassReg}:   {3, 4, 0},
	Optab{AMULD, ClassConst13, ClassReg, ClassNone, ClassReg}:  {3, 4, 0},
	Optab{ASLLD, ClassConst6, ClassReg, ClassNone, ClassReg}:   {3, 4, 0},
	Optab{ASLLW, ClassConst5, ClassReg, ClassNone, ClassReg}:   {3, 4, 0},

	Optab{ASAVE, ClassConst13, ClassReg, ClassNone, ClassReg}: {3, 4, 0},
	Optab{ASAVE, ClassConst13, ClassReg | ClassBias, ClassNone, ClassReg | ClassBias}: {3, 4, 0},

	Optab{AMOVD, ClassConst13, ClassNone, ClassNone, ClassReg}: {4, 4, 0},
	Optab{AMOVW, ClassConst13, ClassNone, ClassNone, ClassReg}: {4, 4, 0},

	Optab{ALDD, ClassIndirRegReg, ClassNone, ClassNone, ClassReg}:   {5, 4, 0},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndirRegReg}:   {6, 4, 0},
	Optab{ALDDF, ClassIndirRegReg, ClassNone, ClassNone, ClassDReg}: {5, 4, 0},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassIndirRegReg}: {6, 4, 0},

	Optab{ALDD, ClassIndir13, ClassNone, ClassNone, ClassReg}:   {7, 4, 0},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndir13}:   {8, 4, 0},
	Optab{ALDDF, ClassIndir13, ClassNone, ClassNone, ClassDReg}: {7, 4, 0},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassIndir13}: {8, 4, 0},

	Optab{ARD, ClassSpcReg, ClassNone, ClassNone, ClassReg}: {9, 4, 0},

	Optab{ACASD, ClassIndir0, ClassReg, ClassNone, ClassReg}: {10, 4, 0},

	Optab{AFSTOD, ClassFReg, ClassNone, ClassNone, ClassDReg}: {11, 4, 0},
	Optab{AFDTOS, ClassDReg, ClassNone, ClassNone, ClassFReg}: {11, 4, 0},

	Optab{AFMOVD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4, 0},

	Optab{AFXTOD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4, ClobberTMP},
	Optab{AFITOD, ClassFReg, ClassNone, ClassNone, ClassDReg}: {11, 4, 0},
	Optab{AFXTOS, ClassDReg, ClassNone, ClassNone, ClassFReg}: {11, 4, ClobberTMP},
	Optab{AFITOS, ClassFReg, ClassNone, ClassNone, ClassFReg}: {11, 4, 0},

	Optab{AFSTOX, ClassFReg, ClassNone, ClassNone, ClassDReg}: {11, 4, 0},
	Optab{AFDTOX, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4, ClobberTMP},
	Optab{AFDTOI, ClassDReg, ClassNone, ClassNone, ClassFReg}: {11, 4, 0},
	Optab{AFSTOI, ClassFReg, ClassNone, ClassNone, ClassFReg}: {11, 4, 0},

	Optab{AFABSD, ClassDReg, ClassNone, ClassNone, ClassDReg}: {11, 4, 0},

	Optab{ASETHI, ClassConst32, ClassNone, ClassNone, ClassReg}: {12, 4, 0},
	Optab{ARNOP, ClassNone, ClassNone, ClassNone, ClassNone}:    {12, 4, 0},
	Optab{AFLUSHW, ClassNone, ClassNone, ClassNone, ClassNone}:  {12, 4, 0},

	Optab{AMEMBAR, ClassConst, ClassNone, ClassNone, ClassNone}: {13, 4, 0},

	Optab{AFCMPD, ClassDReg, ClassDReg, ClassNone, ClassFCond}: {14, 4, 0},
	Optab{AFCMPD, ClassDReg, ClassDReg, ClassNone, ClassNone}:  {14, 4, 0},

	Optab{AMOVW, ClassConst32, ClassNone, ClassNone, ClassReg}:  {15, 8, 0},
	Optab{AMOVW, ClassConst31_, ClassNone, ClassNone, ClassReg}: {16, 8, 0},

	Optab{AMOVD, ClassConst32, ClassNone, ClassNone, ClassReg}:  {15, 8, 0},
	Optab{AMOVD, ClassConst31_, ClassNone, ClassNone, ClassReg}: {16, 8, 0},

	Optab{obj.AJMP, ClassNone, ClassNone, ClassNone, ClassBranch}: {17, 4, 0},
	Optab{ABN, ClassCond, ClassNone, ClassNone, ClassBranch}:      {17, 4, 0},
	Optab{ABNW, ClassNone, ClassNone, ClassNone, ClassBranch}:     {17, 4, 0},
	Optab{ABRZ, ClassReg, ClassNone, ClassNone, ClassBranch}:      {18, 4, 0},
	Optab{AFBA, ClassNone, ClassNone, ClassNone, ClassBranch}:     {19, 4, 0},

	Optab{AJMPL, ClassReg, ClassNone, ClassNone, ClassReg}:        {20, 4, 0},
	Optab{AJMPL, ClassRegConst13, ClassNone, ClassNone, ClassReg}: {20, 4, 0},
	Optab{AJMPL, ClassRegReg, ClassNone, ClassNone, ClassReg}:     {21, 4, 0},

	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassMem}:     {22, 4, 0},
	Optab{obj.ADUFFZERO, ClassNone, ClassNone, ClassNone, ClassMem}: {22, 4, 0},
	Optab{obj.ADUFFCOPY, ClassNone, ClassNone, ClassNone, ClassMem}: {22, 4, 0},

	Optab{AMOVD, ClassAddr, ClassNone, ClassNone, ClassReg}: {23, 24, ClobberTMP},

	Optab{ALDD, ClassMem, ClassNone, ClassNone, ClassReg}:   {24, 28, ClobberTMP},
	Optab{ALDDF, ClassMem, ClassNone, ClassNone, ClassDReg}: {24, 28, ClobberTMP},
	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassMem}:   {25, 28, ClobberTMP},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassMem}: {25, 28, ClobberTMP},

	Optab{obj.ARET, ClassNone, ClassNone, ClassNone, ClassNone}: {26, 4, 0},

	Optab{ATA, ClassConst13, ClassNone, ClassNone, ClassNone}: {27, 4, 0},

	Optab{AMOVD, ClassRegConst13, ClassNone, ClassNone, ClassReg}: {28, 4, 0},

	Optab{AMOVUB, ClassReg, ClassNone, ClassNone, ClassReg}: {29, 4, 0},
	Optab{AMOVUH, ClassReg, ClassNone, ClassNone, ClassReg}: {30, 8, 0},
	Optab{AMOVUW, ClassReg, ClassNone, ClassNone, ClassReg}: {31, 4, 0},

	Optab{AMOVB, ClassReg, ClassNone, ClassNone, ClassReg}: {32, 8, 0},
	Optab{AMOVH, ClassReg, ClassNone, ClassNone, ClassReg}: {33, 8, 0},
	Optab{AMOVW, ClassReg, ClassNone, ClassNone, ClassReg}: {34, 4, 0},

	Optab{ANEG, ClassReg, ClassNone, ClassNone, ClassReg}: {35, 4, 0},

	Optab{ACMP, ClassReg, ClassReg, ClassNone, ClassNone}:     {36, 4, 0},
	Optab{ACMP, ClassConst13, ClassReg, ClassNone, ClassNone}: {37, 4, 0},

	Optab{ABND, ClassNone, ClassNone, ClassNone, ClassBranch}: {38, 4, 0},

	Optab{obj.AUNDEF, ClassNone, ClassNone, ClassNone, ClassNone}: {39, 4, 0},

	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassReg}:    {40, 4, 0},
	Optab{obj.ACALL, ClassReg, ClassNone, ClassNone, ClassReg}:     {40, 4, 0},
	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassIndir0}: {40, 4, 0},
	Optab{obj.ACALL, ClassReg, ClassNone, ClassNone, ClassIndir0}:  {40, 4, 0},

	Optab{AADD, ClassConst32, ClassNone, ClassNone, ClassReg}: {41, 12, ClobberTMP},
	Optab{AAND, ClassConst32, ClassNone, ClassNone, ClassReg}: {41, 12, ClobberTMP},
	Optab{AADD, ClassConst32, ClassReg, ClassNone, ClassReg}:  {41, 12, ClobberTMP},
	Optab{AAND, ClassConst32, ClassReg, ClassNone, ClassReg}:  {41, 12, ClobberTMP},

	Optab{AMOVD, ClassRegConst, ClassNone, ClassNone, ClassReg}: {42, 12, ClobberTMP},

	Optab{ASTD, ClassReg, ClassNone, ClassNone, ClassIndir}:   {43, 12, ClobberTMP},
	Optab{ASTDF, ClassDReg, ClassNone, ClassNone, ClassIndir}: {43, 12, ClobberTMP},
	Optab{ALDD, ClassIndir, ClassNone, ClassNone, ClassReg}:   {44, 12, ClobberTMP},
	Optab{ALDDF, ClassIndir, ClassNone, ClassNone, ClassDReg}: {44, 12, ClobberTMP},

	Optab{obj.AJMP, ClassNone, ClassNone, ClassNone, ClassMem}: {45, 28, ClobberTMP},

	Optab{AMOVA, ClassCond, ClassNone, ClassConst11, ClassReg}: {46, 4, 0},
	Optab{AMOVA, ClassCond, ClassReg, ClassNone, ClassReg}:     {47, 4, 0},

	Optab{AMOVFA, ClassFCond, ClassNone, ClassConst11, ClassReg}: {46, 4, 0},
	Optab{AMOVFA, ClassFCond, ClassReg, ClassNone, ClassReg}:     {47, 4, 0},

	Optab{AMOVRZ, ClassReg, ClassNone, ClassConst10, ClassReg}: {48, 4, 0},
	Optab{AMOVRZ, ClassReg, ClassReg, ClassNone, ClassReg}:     {49, 4, 0},

	Optab{AMOVD, ClassTLSAddr, ClassNone, ClassNone, ClassReg}: {50, 12, 0},

	Optab{ARETRESTORE, ClassNone, ClassNone, ClassNone, ClassNone}: {51, 12, 0},

	Optab{obj.AJMP, ClassNone, ClassNone, ClassNone, ClassLargeBranch}: {52, 28, ClobberTMP},
	Optab{ABN, ClassCond, ClassNone, ClassNone, ClassLargeBranch}:  {53, 48, ClobberTMP},
	Optab{ABNW, ClassNone, ClassNone, ClassNone, ClassLargeBranch}: {53, 48, ClobberTMP},
	Optab{ABRZ, ClassReg, ClassNone, ClassNone, ClassLargeBranch}:  {54, 48, ClobberTMP},
	Optab{AFBA, ClassNone, ClassNone, ClassNone, ClassLargeBranch}: {55, 48, ClobberTMP},
	Optab{ABND, ClassNone, ClassNone, ClassNone, ClassLargeBranch}: {56, 48, ClobberTMP},

	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassBranch}:      {57, 4, 0},
	Optab{obj.ACALL, ClassNone, ClassNone, ClassNone, ClassLargeBranch}: {57, 4, 0},
}

// Compatible classes, if something accepts a $hugeconst, it
// can also accept $smallconst, $0 and ZR. Something that accepts a
// register, can also accept $0, etc.
var cc = map[int8][]int8{
	ClassReg:         {ClassZero},
	ClassConst6:      {ClassConst5, ClassZero},
	ClassConst10:     {ClassConst6, ClassConst5, ClassZero},
	ClassConst11:     {ClassConst10, ClassConst6, ClassConst5, ClassZero},
	ClassConst13:     {ClassConst11, ClassConst10, ClassConst6, ClassConst5, ClassZero},
	ClassConst31:     {ClassConst6, ClassConst5, ClassZero},
	ClassConst32:     {ClassConst31_, ClassConst31, ClassConst13, ClassConst11, ClassConst10, ClassConst6, ClassConst5, ClassZero},
	ClassConst:       {ClassConst32, ClassConst31_, ClassConst31, ClassConst13, ClassConst11, ClassConst10, ClassConst6, ClassConst5, ClassZero},
	ClassRegConst:    {ClassRegConst13},
	ClassIndir13:     {ClassIndir0},
	ClassIndir:       {ClassIndir13, ClassIndir0},
	ClassLargeBranch: {ClassBranch},
}

func isAddrCompatible(ctxt *obj.Link, a *obj.Addr, class int8) bool {
	cls := aclass(ctxt, a)
	cls &= ^ClassBias
	if cls == class {
		return true
	}
	for _, v := range cc[class] {
		if cls == v {
			return true
		}
	}
	return false
}

var isInstDouble = map[obj.As]bool{
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

var isInstFloat = map[obj.As]bool{
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

var isSrcDouble = map[obj.As]bool{
	AFXTOD: true,
	AFXTOS: true,
	AFDTOX: true,
	AFDTOI: true,
	AFDTOS: true,
}

var isSrcFloat = map[obj.As]bool{
	AFITOD: true,
	AFITOS: true,
	AFSTOX: true,
	AFSTOI: true,
	AFSTOD: true,
}

var isDstDouble = map[obj.As]bool{
	AFXTOD: true,
	AFITOD: true,
	AFSTOX: true,
	AFDTOX: true,
	AFSTOD: true,
}

var isDstFloat = map[obj.As]bool{
	AFXTOS: true,
	AFITOS: true,
	AFDTOI: true,
	AFSTOI: true,
	AFDTOS: true,
}

// Compatible instructions, if an asm* function accepts AADD,
// it accepts ASUBCCC too.
var ci = map[obj.As][]obj.As{
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
	AMOVA:  {AMOVN, AMOVNE, AMOVE, AMOVG, AMOVLE, AMOVGE, AMOVL, AMOVGU, AMOVLEU, AMOVCC, AMOVCS, AMOVPOS, AMOVNEG, AMOVVC, AMOVVS},
	AMOVFA: {AMOVFN, AMOVFU, AMOVFG, AMOVFUG, AMOVFL, AMOVFUL, AMOVFLG, AMOVFNE, AMOVFE, AMOVFUE, AMOVFGE, AMOVFUGE, AMOVFLE, AMOVFULE, AMOVFO},
	AMOVRZ: {AMOVRLEZ, AMOVRLZ, AMOVRNZ, AMOVRGZ, AMOVRGEZ},
	AMULD:  {ASDIVD, AUDIVD},
	ARD:    {AMOVD},
	ASLLD:  {ASRLD, ASRAD},
	ASLLW:  {ASLLW, ASRLW, ASRAW},
	ASTD:   {ASTB, ASTH, ASTW, AMOVB, AMOVH, AMOVW, AMOVUB, AMOVUH, AMOVUW, AMOVD},
	ASTDF:  {ASTSF, AFMOVD, AFMOVS},
	ASAVE:  {ARESTORE},
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
	var type3 obj.AddrType
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

func opload(a obj.As) uint32 {
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
		panic("unknown instruction: " + a.String())
	}
}

func opstore(a obj.As) uint32 {
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
		panic("unknown instruction: " + a.String())
	}
}

func oprd(a obj.As) uint32 {
	switch a {
	// Read ancillary state register.
	case ARD, AMOVD:
		return op3(2, 0x28)

	default:
		panic("unknown instruction: " + a.String())
	}
}

func opalu(a obj.As) uint32 {
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
	case AOR, AMOVD, AMOVW:
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

	case ASAVE:
		return op3(2, 0x3C)
	case ARESTORE:
		return op3(2, 0x3D)

	default:
		panic("unknown instruction: " + a.String())
	}
}

func opcode(a obj.As) uint32 {
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
		return op3(3, 0x3C)
	case ACASD:
		return op3(3, 0x3E)

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

	case AFLUSHW:
		return op3(2, 0x2B)

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

	// Move Integer Register on Condition (MOVcc).
	case AMOVA:
		return op3(2, 0x2C) | 8<<14 | 1<<18
	case AMOVN:
		return op3(2, 0x2C) | 0<<14 | 1<<18
	case AMOVNE:
		return op3(2, 0x2C) | 9<<14 | 1<<18
	case AMOVE:
		return op3(2, 0x2C) | 1<<14 | 1<<18
	case AMOVG:
		return op3(2, 0x2C) | 10<<14 | 1<<18
	case AMOVLE:
		return op3(2, 0x2C) | 2<<14 | 1<<18
	case AMOVGE:
		return op3(2, 0x2C) | 11<<14 | 1<<18
	case AMOVL:
		return op3(2, 0x2C) | 3<<14 | 1<<18
	case AMOVGU:
		return op3(2, 0x2C) | 12<<14 | 1<<18
	case AMOVLEU:
		return op3(2, 0x2C) | 4<<14 | 1<<18
	case AMOVCC:
		return op3(2, 0x2C) | 13<<14 | 1<<18
	case AMOVCS:
		return op3(2, 0x2C) | 5<<14 | 1<<18
	case AMOVPOS:
		return op3(2, 0x2C) | 14<<14 | 1<<18
	case AMOVNEG:
		return op3(2, 0x2C) | 6<<14 | 1<<18
	case AMOVVC:
		return op3(2, 0x2C) | 15<<14 | 1<<18
	case AMOVVS:
		return op3(2, 0x2C) | 7<<14 | 1<<18

	// Move Integer Register on Floating-Point Condition (MOVcc).
	case AMOVFA:
		return op3(2, 0x2C) | 8<<14 | 0<<18
	case AMOVFN:
		return op3(2, 0x2C) | 0<<14 | 0<<18
	case AMOVFU:
		return op3(2, 0x2C) | 7<<14 | 0<<18
	case AMOVFG:
		return op3(2, 0x2C) | 6<<14 | 0<<18
	case AMOVFUG:
		return op3(2, 0x2C) | 5<<14 | 0<<18
	case AMOVFL:
		return op3(2, 0x2C) | 4<<14 | 0<<18
	case AMOVFUL:
		return op3(2, 0x2C) | 3<<14 | 0<<18
	case AMOVFLG:
		return op3(2, 0x2C) | 2<<14 | 0<<18
	case AMOVFNE:
		return op3(2, 0x2C) | 1<<14 | 0<<18
	case AMOVFE:
		return op3(2, 0x2C) | 9<<14 | 0<<18
	case AMOVFUE:
		return op3(2, 0x2C) | 10<<14 | 0<<18
	case AMOVFGE:
		return op3(2, 0x2C) | 11<<14 | 0<<18
	case AMOVFUGE:
		return op3(2, 0x2C) | 12<<14 | 0<<18
	case AMOVFLE:
		return op3(2, 0x2C) | 13<<14 | 0<<18
	case AMOVFULE:
		return op3(2, 0x2C) | 14<<14 | 0<<18
	case AMOVFO:
		return op3(2, 0x2C) | 15<<14 | 0<<18

	// Move Integer Register on Register Condition (MOVr).
	case AMOVRZ:
		return op3(2, 0x2f) | 1<<10
	case AMOVRLEZ:
		return op3(2, 0x2f) | 2<<10
	case AMOVRLZ:
		return op3(2, 0x2f) | 3<<10
	case AMOVRNZ:
		return op3(2, 0x2f) | 5<<10
	case AMOVRGZ:
		return op3(2, 0x2f) | 6<<10
	case AMOVRGEZ:
		return op3(2, 0x2f) | 7<<10

	// Memory Barrier.
	case AMEMBAR:
		return op3(2, 0x28) | 0xF<<14 | 1<<13

	case ASETHI, ARNOP:
		return op2(4)

	// Trap on Integer Condition Codes (Tcc).
	case ATA:
		return op3(2, 0x3A)

	default:
		panic("unknown instruction: " + a.String())
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
	if -512 <= offset && offset <= 513 {
		return ClassConst10
	}
	if -1024 <= offset && offset <= 1023 {
		return ClassConst11
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
	case REG_G1 <= r && r <= REG_I7:
		return ClassReg
	case REG_F0 <= r && r <= REG_F31:
		return ClassFReg
	case REG_D0 <= r && r <= REG_D62:
		return ClassDReg
	case r == REG_ICC || r == REG_XCC:
		return ClassCond
	case REG_FCC0 <= r && r <= REG_FCC3:
		return ClassFCond
	case r == REG_BSP || r == REG_BFP:
		return ClassReg | ClassBias
	case r >= REG_SPECIAL:
		return ClassSpcReg
	}
	return ClassUnknown
}

func aclass(ctxt *obj.Link, a *obj.Addr) int8 {
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
			if a.Sym.Type == obj.STLSBSS {
				return ClassTLSMem
			}
			return ClassMem

		case obj.NAME_AUTO, obj.NAME_PARAM:
			return aclass(ctxt, autoeditaddr(ctxt, a))

		case obj.NAME_NONE:
			if a.Scale == 1 {
				return ClassIndirRegReg
			}
			return oregclass(a.Offset) | rclass(a.Reg)&ClassBias
		}

	case obj.TYPE_FCONST:
		return ClassFConst

	case obj.TYPE_TEXTSIZE:
		return ClassTextSize

	case obj.TYPE_CONST, obj.TYPE_ADDR:
		switch a.Name {
		case obj.NAME_NONE:
			if a.Reg != 0 {
				if a.Reg == REG_ZR && a.Offset == 0 {
					return ClassZero
				}
				if a.Scale == 1 {
					return ClassRegReg
				}
				return addrclass(a.Offset) | rclass(a.Reg)&ClassBias
			}
			return constclass(a.Offset)

		case obj.NAME_EXTERN, obj.NAME_STATIC:
			if a.Sym == nil {
				return ClassUnknown
			}
			if a.Sym.Type == obj.STLSBSS {
				return ClassTLSAddr
			}
			return ClassAddr

		case obj.NAME_AUTO, obj.NAME_PARAM:
			return aclass(ctxt, autoeditaddr(ctxt, a))
		}
	case obj.TYPE_BRANCH:
		if a.Class == ClassLargeBranch {
			// Set by span() after initial pcs have been calculated.
			return ClassLargeBranch
		}
		return ClassBranch
	}
	return ClassUnknown
}

// Assign Pcs and reclassify branches that exceed the standard 21-bit signed
// maximum offset and recalculate; multiple invocations may be required.
func assignPc(ctxt *obj.Link, cursym *obj.LSym, prevpc int64) (pc int64) {
	expandBranch := prevpc > 0
	for p := cursym.Text.Link; p != nil; p = p.Link {
		if expandBranch && p.To.Type == obj.TYPE_BRANCH && p.To.Class == ClassBranch {
			var offset int64
			if p.Pcond != nil {
				offset = p.Pcond.Pc - p.Pc
			} else {
				// obj.brloop will set p.Pcond to nil for jumps
				// to the same instruction.
				offset = p.To.Val.(*obj.Prog).Pc - p.Pc
			}
			if offset < -1<<20 || offset > 1<<20-1 {
				// Ideally, this would be done in aclass(), but
				// we don't have access to p there or the pc
				// (yet) in most cases. oplook will use this to
				// transform the branch appropriately so that
				// asmout will perform a "large" branch.
				p.To.Class = ClassLargeBranch
			}
		}

		o, err := oplook(autoeditprog(ctxt, p))
		if err != nil {
			ctxt.Diag(err.Error())
		}

		p.Pc = pc
		pc += int64(o.size)
	}

	if prevpc == 0 {
		// After initial Pc assignment, reassign until Pc no longer
		// increases.
		prevpc := pc
		for {
			pc = assignPc(ctxt, cursym, prevpc)
			if pc <= prevpc {
				break
			}
			prevpc = pc
		}
	}
	return
}

func span(ctxt *obj.Link, cursym *obj.LSym) {
	if cursym.Text == nil || cursym.Text.Link == nil { // handle external functions and ELF section symbols
		return
	}

	var pc = assignPc(ctxt, cursym, 0)
	cursym.Size = pc
	cursym.Grow(cursym.Size)

	var text []uint32 // actual assembled bytes
	for p := cursym.Text.Link; p != nil; p = p.Link {
		p1 := autoeditprog(ctxt, p)
		o, _ := oplook(p1)
		out, err := asmout(p1, o, cursym)
		if err != nil {
			ctxt.Diag("span: can't assemble: %v\n\t%v", err, p)
		}
		text = append(text, out...)
	}

	bp := cursym.P
	for _, v := range text {
		ctxt.Arch.ByteOrder.PutUint32(bp, v)
		bp = bp[4:]
	}
}

// bigmove assembles a move of the constant part of addr into reg.
func bigmove(ctxt *obj.Link, addr *obj.Addr, reg int16) (out []uint32) {
	out = make([]uint32, 2)
	class := aclass(ctxt, addr)
	switch class {
	case ClassRegConst, ClassIndir:
		class = constclass(addr.Offset)
	}
	switch class {
	// MOV[WD] $imm32, R ->
	// 	SETHI hi($imm32), R
	// 	OR R, lo($imm32), R
	case ClassConst31, ClassConst32:
		out[0] = opcode(ASETHI) | ir(uint32(addr.Offset)>>10, reg)
		out[1] = opalu(AOR) | rsr(reg, int64(addr.Offset&0x3FF), reg)

	// MOV[WD] -$imm31, R ->
	// 	SETHI hi(^$imm32), R
	// 	XOR R, lo($imm32)|0x1C00, R
	case ClassConst31_:
		out[0] = opcode(ASETHI) | ir(^(uint32(addr.Offset))>>10, reg)
		out[1] = opalu(AXOR) | rsr(reg, int64(uint32(addr.Offset)&0x3ff|0x1C00), reg)
	default:
		panic("unexpected operand class: " + DRconv(class))
	}
	return out
}

func usesRegs(a *obj.Addr) bool {
	if a == nil {
		return false
	}
	switch a.Class {
	case ClassReg, ClassFReg, ClassDReg, ClassCond, ClassFCond, ClassSpcReg, ClassZero, ClassRegReg, ClassRegConst13, ClassRegConst, ClassIndirRegReg, ClassIndir0, ClassIndir13, ClassIndir:
		return true
	}
	return false
}

func isTMP(r int16) bool {
	return r == REG_TMP || r == REG_TMP2
}

func usesTMP(a *obj.Addr) bool {
	return usesRegs(a) && (isTMP(a.Reg) || isTMP(a.Index))
}

func srcCount(p *obj.Prog) (c int) {
	if p.From.Type != obj.TYPE_NONE {
		c++
	}
	if p.Reg != obj.REG_NONE {
		c++
	}
	if p.From3Type() != obj.TYPE_NONE {
		c++
	}
	return c
}

// largebranch assembles a branch to a pc that exceeds a 21-bit signed displacement
func largebranch(offset int64) ([]uint32, error) {
	if offset%4 != 0 {
		return nil, errors.New("branch target not mod 4")
	}

	out := make([]uint32, 7)
	// We don't know where we are, and we don't want to emit a
	// reloc, so save %o7 since we may be in the function prologue,
	// then do a pc-relative call to determine current address,
	// then restore %o7 so that we can use the current address plus
	// the calculated offset to perform a "large" jump to the
	// desired location.
	out[0] = opalu(AMOVD) | rrr(REG_ZR, 0, REG_OLR, REG_TMP2)
	out[1] = opcode(obj.ACALL) | d30(1)
	out[2] = opalu(AMOVD) | rrr(REG_ZR, 0, REG_OLR, REG_TMP)
	out[3] = opalu(AMOVD) | rrr(REG_ZR, 0, REG_TMP2, REG_OLR)
	offset -= 4 // make branch relative to call
	class := constclass(offset)
	switch class {
	// 	SETHI hi($imm32), R
	// 	OR R, lo($imm32), R
	case ClassConst31, ClassConst32:
		out[4] = opcode(ASETHI) | ir(uint32(offset)>>10, REG_TMP2)
		out[5] = opalu(AOR) | rsr(REG_TMP2, int64(offset&0x3FF), REG_TMP2)

	// 	SETHI hi(^$imm32), R
	// 	XOR R, lo($imm32)|0x1C00, R
	case ClassConst31_:
		out[4] = opcode(ASETHI) | ir(^(uint32(offset))>>10, REG_TMP2)
		out[5] = opalu(AXOR) | rsr(REG_TMP2, int64(uint32(offset)&0x3ff|0x1C00), REG_TMP2)
	default:
		panic("unexpected operand class: " + DRconv(class))
	}
	out[6] = opcode(AJMPL) | rrr(REG_TMP, 0, REG_TMP2, REG_ZR)
	return out, nil
}

func asmout(p *obj.Prog, o Opval, cursym *obj.LSym) (out []uint32, err error) {
	out = make([]uint32, 12)
	o1 := &out[0]
	o2 := &out[1]
	o3 := &out[2]
	o4 := &out[3]
	o5 := &out[4]
	o6 := &out[5]
	o7 := &out[6]
	o8 := &out[7]
	o9 := &out[8]
	o10 := &out[9]
	o11 := &out[10]
	o12 := &out[11]
	if o.OpInfo == ClobberTMP {
		if usesTMP(&p.From) {
			return nil, fmt.Errorf("asmout: %q not allowed: synthetic instruction clobbers temporary registers", obj.Mconv(&p.From))
		}
		if isTMP(p.Reg) {
			return nil, fmt.Errorf("asmout: %q not allowed: synthetic instruction clobbers temporary registers", Rconv(int(p.Reg)))
		}
		if usesTMP(p.From3) {
			return nil, fmt.Errorf("asmout: %q not allowed: synthetic instruction clobbers temporary registers", obj.Mconv(p.From3))
		}
		if usesTMP(&p.To) {
			if p.From.Type == obj.TYPE_NONE || srcCount(p) < 2 {
				return nil, fmt.Errorf("asmout: illegal use of temporary register: synthetic instruction clobbers temporary registers")
			}
		}
	}
	switch o.op {
	default:
		return nil, fmt.Errorf("unknown asm %d in %v", o, p)

	case 0: /* pseudo ops */
		break

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
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0x80, p.Reg, p.To.Reg)

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
		*o1 = opcode(p.As) | rrr(p.Reg, 0, p.From.Reg, p.To.Reg&3)

	// MOVW $imm32, R
	// MOVW -$imm31, R
	// MOVD $imm32, R
	// MOVD -$imm31, R
	case 15, 16:
		out := bigmove(p.Ctxt, &p.From, p.To.Reg)
		return out, nil

	// BLE XCC, n(PC)
	// JMP n(PC)
	case 17:
		var offset int64
		if p.Pcond != nil {
			offset = p.Pcond.Pc - p.Pc
		} else {
			// obj.brloop will set p.Pcond to nil for jumps to the same instruction.
			offset = p.To.Val.(*obj.Prog).Pc - p.Pc
		}
		if offset < -1<<20 || offset > 1<<20-1 {
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
		if offset < -1<<24 || offset > 1<<24-1 {
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
	// 	SETHI hh($sym), TMP
	// 	OR TMP, hm($sym), TMP
	//	SLLD	$32, TMP, TMP
	// 	SETHI hi($sym), R
	// 	OR R, lo($sym), R
	// 	OR TMP, R, R
	case 23:
		*o1 = opcode(ASETHI) | ir(0, REG_TMP)
		*o2 = opalu(AOR) | rsr(REG_TMP, 0, REG_TMP)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64HI
		*o3 = opalu(ASLLD) | rsr(REG_TMP, 32, REG_TMP)
		*o4 = opcode(ASETHI) | ir(0, p.To.Reg)
		*o5 = opalu(AOR) | rsr(p.To.Reg, 0, p.To.Reg)
		rel = obj.Addrel(cursym)
		rel.Off = int32(p.Pc + 12)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64LO
		*o6 = opalu(AOR) | rrr(REG_TMP, 0, p.To.Reg, p.To.Reg)

	// MOV sym(SB), R ->
	// 	SETHI hh($sym), TMP
	// 	OR TMP, hm($sym), TMP
	//	SLLD	$32, TMP, TMP
	// 	SETHI hi($sym), TMP2
	// 	OR TMP2, lo($sym), TMP2
	// 	OR TMP, TMP2, TMP2
	//	MOV (TMP2), R
	case 24:
		*o1 = opcode(ASETHI) | ir(0, REG_TMP)
		*o2 = opalu(AOR) | rsr(REG_TMP, 0, REG_TMP)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64HI
		*o3 = opalu(ASLLD) | rsr(REG_TMP, 32, REG_TMP)
		*o4 = opcode(ASETHI) | ir(0, REG_TMP2)
		*o5 = opalu(AOR) | rsr(REG_TMP2, 0, REG_TMP2)
		rel = obj.Addrel(cursym)
		rel.Off = int32(p.Pc + 12)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_ADDRSPARC64LO
		*o6 = opalu(AOR) | rrr(REG_TMP, 0, REG_TMP2, REG_TMP2)
		*o7 = opload(p.As) | rsr(REG_TMP2, 0, p.To.Reg)

	// MOV R, sym(SB) ->
	// 	SETHI hh($sym), TMP
	// 	OR TMP, hm($sym), TMP
	//	SLLD	$32, TMP, TMP
	// 	SETHI hi($sym), TMP2
	// 	OR TMP2, lo($sym), TMP2
	// 	OR TMP, TMP2, TMP2
	//	MOV R, (TMP2)
	case 25:
		*o1 = opcode(ASETHI) | ir(0, REG_TMP)
		*o2 = opalu(AOR) | rsr(REG_TMP, 0, REG_TMP)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_ADDRSPARC64HI
		*o3 = opalu(ASLLD) | rsr(REG_TMP, 32, REG_TMP)
		*o4 = opcode(ASETHI) | ir(0, REG_TMP2)
		*o5 = opalu(AOR) | rsr(REG_TMP2, 0, REG_TMP2)
		rel = obj.Addrel(cursym)
		rel.Off = int32(p.Pc + 12)
		rel.Siz = 8
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_ADDRSPARC64LO
		*o6 = opalu(AOR) | rrr(REG_TMP, 0, REG_TMP2, REG_TMP2)
		*o7 = opstore(p.As) | rsr(REG_TMP2, 0, p.From.Reg)

	// RET
	case 26:
		*o1 = opcode(AJMPL) | rsr(REG_OLR, 8, REG_ZR)

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
		if offset < -1<<20 || offset > 1<<20-1 {
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
		*o1 = opcode(AJMPL) | rsr(p.To.Reg, 0, REG_OLR)

	// ADD $huge, Rd
	// AND $huge, Rs, Rd
	case 41:
		move := bigmove(p.Ctxt, &p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		reg := p.To.Reg
		if p.Reg != 0 {
			reg = p.Reg
		}
		*o3 = opalu(p.As) | rrr(reg, 0, REG_TMP, p.To.Reg)

	// AMOVD $huge(R), R
	case 42:
		move := bigmove(p.Ctxt, &p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opalu(AADD) | rrr(p.From.Reg, 0, REG_TMP, p.To.Reg)

	// AMOVD R, huge(R)
	case 43:
		move := bigmove(p.Ctxt, &p.To, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opstore(p.As) | rrr(p.To.Reg, 0, REG_TMP, p.From.Reg)

	// AMOVD huge(R), R
	case 44:
		move := bigmove(p.Ctxt, &p.From, REG_TMP)
		*o1, *o2 = move[0], move[1]
		*o3 = opload(p.As) | rrr(p.From.Reg, 0, REG_TMP, p.To.Reg)

	// JMP sym(SB) ->
	//	MOVD	$sym(SB), TMP2 ->
	// 		SETHI hh($sym), TMP
	// 		OR TMP, hm($sym), TMP
	//		SLLD	$32, TMP, TMP
	// 		SETHI hi($sym), TMP2
	// 		OR TMP2, lo($sym), TMP2
	// 		OR TMP, TMP2, TMP2
	//	JMPL	TMP2, ZR
	case 45:
		*o1 = opcode(ASETHI) | ir(0, REG_TMP)
		*o2 = opalu(AOR) | rsr(REG_TMP, 0, REG_TMP)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_ADDRSPARC64HI
		*o3 = opalu(ASLLD) | rsr(REG_TMP, 32, REG_TMP)
		*o4 = opcode(ASETHI) | ir(0, REG_TMP2)
		*o5 = opalu(AOR) | rsr(REG_TMP2, 0, REG_TMP2)
		rel = obj.Addrel(cursym)
		rel.Off = int32(p.Pc + 12)
		rel.Siz = 8
		rel.Sym = p.To.Sym
		rel.Add = p.To.Offset
		rel.Type = obj.R_ADDRSPARC64LO
		*o6 = opalu(AOR) | rrr(REG_TMP, 0, REG_TMP2, REG_TMP2)
		*o7 = opcode(AJMPL) | rsr(REG_TMP2, 0, REG_ZR)

	// MOV[F]A ICC/XCC/FCC, $simm11, R
	case 46:
		*o1 = opcode(p.As) | rsr(0, p.From3.Offset, p.To.Reg) | 1<<13 | uint32(p.From.Reg&3<<11)

	// MOV[F]A ICC/XCC/FCC, R, R
	case 47:
		*o1 = opcode(p.As) | rrr(0, 0, p.Reg, p.To.Reg) | uint32(p.From.Reg&3<<11)

	// MOVRZ	R, $simm10, Rd
	case 48:
		*o1 = opcode(p.As) | rsr(p.From.Reg, p.From3.Offset, p.To.Reg) | 1<<13

	// MOVRZ	R, Rs, Rd
	case 49:
		*o1 = opcode(p.As) | rrr(p.From.Reg, 0, p.Reg, p.To.Reg)

	// MOVD $tlssym, R
	case 50:
		*o1 = opcode(ASETHI) | ir(0, p.To.Reg)
		*o2 = opalu(AXOR) | rsr(p.To.Reg, 0, p.To.Reg)
		rel := obj.Addrel(cursym)
		rel.Off = int32(p.Pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Add = p.From.Offset
		rel.Type = obj.R_SPARC64_TLS_LE
		*o3 = opalu(AADD) | rrr(REG_TLS, 0, p.To.Reg, p.To.Reg)

	// RETRESTORE
	case 51:
		*o1 = opload(AMOVD) | rsr(REG_RSP, StackBias+120, REG_ILR)
		*o2 = opcode(AJMPL) | rsr(REG_ILR, 8, REG_ZR)
		*o3 = opalu(ARESTORE) | rsr(REG_ZR, 0, REG_ZR)

	// JMP $huge(n(PC)) ->
	//	MOVD	OLR, TMP2
	//	CALL	+0x4
	//	MOVD	OLR, TMP
	//	MOVD	TMP2, OLR
	//	MOVD	$huge(n(PC)), TMP2
	//	...
	//	JMPL	TMP + TMP2
	case 52:
		var offset int64
		if p.Pcond != nil {
			offset = p.Pcond.Pc - p.Pc
		} else {
			// obj.brloop will set p.Pcond to nil for jumps to the same instruction.
			offset = p.To.Val.(*obj.Prog).Pc - p.Pc
		}

		branch, err := largebranch(offset)
		if err != nil {
			return nil, err
		}
		*o1, *o2, *o3, *o4, *o5, *o6, *o7 =
			branch[0], branch[1], branch[2],
			branch[3], branch[4], branch[5],
			branch[6]

	// BLE XCC, $huge(n(PC)) ->
	//	BLE	XCC, 4(PC)
	//	NOP
	//	BA	10(PC)
	//	NOP
	//	MOVD	OLR, TMP2
	//	CALL	+0x4
	//	MOVD	OLR, TMP
	//	MOVD	TMP2, OLR
	//	MOVD	$huge(n(PC)), TMP2
	//	...
	//	JMP	TMP + TMP2
	//	NOP
	case 53:
		offset := int64(16)
		*o1 = opcode(p.As) | uint32(p.From.Reg&3)<<20 | uint32(offset>>2)&(1<<19-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}

		if p.Pcond != nil {
			offset = p.Pcond.Pc - p.Pc
		} else {
			// obj.brloop will set p.Pcond to nil for jumps to the same instruction.
			offset = p.To.Val.(*obj.Prog).Pc - p.Pc
		}
		*o2 = opcode(ARNOP)
		*o3 = opcode(obj.AJMP) | uint32(10)&(1<<22-1)
		*o4 = opcode(ARNOP)

		offset = p.Pcond.Pc - p.Pc
		offset -= 16 // make branch relative to first instruction
		branch, err := largebranch(offset)
		if err != nil {
			return nil, err
		}
		*o5, *o6, *o7, *o8, *o9, *o10, *o11 =
			branch[0], branch[1], branch[2],
			branch[3], branch[4], branch[5],
			branch[6]
		*o12 = opcode(ARNOP)

	// BRZ R, $huge(n(PC)) ->
	//	BRZ	R, 4(PC)
	//	NOP
	//	BA	10(PC)
	//	NOP
	//	MOVD	OLR, TMP2
	//	CALL	+0x4
	//	MOVD	OLR, TMP
	//	MOVD	TMP2, OLR
	//	MOVD	$huge(n(PC)), TMP2
	//	...
	//	JMP	TMP + TMP2
	//	NOP
	case 54:
		offset := int64(16)
		*o1 = opcode(p.As) | uint32((offset>>14)&3)<<20 | uint32(p.From.Reg&31)<<14 | uint32(offset>>2)&(1<<14-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}
		*o2 = opcode(ARNOP)
		*o3 = opcode(obj.AJMP) | uint32(10)&(1<<22-1)
		*o4 = opcode(ARNOP)

		offset = p.Pcond.Pc - p.Pc
		offset -= 16 // make branch relative to first instruction
		branch, err := largebranch(offset)
		if err != nil {
			return nil, err
		}
		*o5, *o6, *o7, *o8, *o9, *o10, *o11 =
			branch[0], branch[1], branch[2],
			branch[3], branch[4], branch[5],
			branch[6]
		*o12 = opcode(ARNOP)

	// FBA $huge(n(PC)) ->
	//	FBA	4(PC)
	//	NOP
	//	BA	10(PC)
	//	NOP
	//	MOVD	OLR, TMP2
	//	CALL	+0x4
	//	MOVD	OLR, TMP
	//	MOVD	TMP2, OLR
	//	MOVD	$huge(n(PC)), TMP2
	//	...
	//	JMP	TMP + TMP2
	//	NOP
	case 55:
		offset := int64(16)
		*o1 = opcode(p.As) | uint32(offset>>2)&(1<<22-1)
		*o2 = opcode(ARNOP)
		*o3 = opcode(obj.AJMP) | uint32(10)&(1<<22-1)
		*o4 = opcode(ARNOP)

		offset = p.Pcond.Pc - p.Pc
		offset -= 16 // make branch relative to first instruction
		branch, err := largebranch(offset)
		if err != nil {
			return nil, err
		}
		*o5, *o6, *o7, *o8, *o9, *o10, *o11 =
			branch[0], branch[1], branch[2],
			branch[3], branch[4], branch[5],
			branch[6]
		*o12 = opcode(ARNOP)

	// BLED, $huge(n(PC)) ->
	//	BLED, 4(PC)
	//	NOP
	//	BA	10(PC)
	//	NOP
	//	MOVD	OLR, TMP2
	//	CALL	+0x4
	//	MOVD	OLR, TMP
	//	MOVD	TMP2, OLR
	//	MOVD	$huge(n(PC)), TMP2
	//	...
	//	JMP	TMP + TMP2
	//	NOP
	case 56:
		offset := int64(16)
		*o1 = opcode(p.As) | 2<<20 | uint32(offset>>2)&(1<<19-1)
		// default is to predict branch taken
		if p.Scond == 0 {
			*o1 |= 1 << 19
		}
		*o2 = opcode(ARNOP)
		*o3 = opcode(obj.AJMP) | uint32(10)&(1<<22-1)
		*o4 = opcode(ARNOP)

		offset = p.Pcond.Pc - p.Pc
		offset -= 16 // make branch relative to first instruction
		branch, err := largebranch(offset)
		if err != nil {
			return nil, err
		}
		*o5, *o6, *o7, *o8, *o9, *o10, *o11 =
			branch[0], branch[1], branch[2],
			branch[3], branch[4], branch[5],
			branch[6]
		*o12 = opcode(ARNOP)

	// CALL n(PC)
	// CALL $huge(n(PC))
	case 57:
		var offset int64
		if p.Pcond != nil {
			offset = p.Pcond.Pc - p.Pc
		} else {
			// obj.brloop will set p.Pcond to nil for jumps to the same instruction.
			offset = p.To.Val.(*obj.Prog).Pc - p.Pc
		}
		if offset < -1<<31 || offset > 1<<31-4 {
			return nil, errors.New("branch target out of range")
		}
		if offset%4 != 0 {
			return nil, errors.New("branch target not mod 4")
		}
		*o1 = opcode(obj.ACALL) | d30(int(offset>>2))
	}

	return out[:o.size/4], nil
}

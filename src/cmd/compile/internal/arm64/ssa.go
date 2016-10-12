// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arm64

import (
	"math"

	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/arm64"
)

var ssaRegToReg = []int16{
	arm64.REG_R0,
	arm64.REG_R1,
	arm64.REG_R2,
	arm64.REG_R3,
	arm64.REG_R4,
	arm64.REG_R5,
	arm64.REG_R6,
	arm64.REG_R7,
	arm64.REG_R8,
	arm64.REG_R9,
	arm64.REG_R10,
	arm64.REG_R11,
	arm64.REG_R12,
	arm64.REG_R13,
	arm64.REG_R14,
	arm64.REG_R15,
	arm64.REG_R16,
	arm64.REG_R17,
	arm64.REG_R18, // platform register, not used
	arm64.REG_R19,
	arm64.REG_R20,
	arm64.REG_R21,
	arm64.REG_R22,
	arm64.REG_R23,
	arm64.REG_R24,
	arm64.REG_R25,
	arm64.REG_R26,
	// R27 = REGTMP not used in regalloc
	arm64.REGG,    // R28
	arm64.REG_R29, // frame pointer, not used
	// R30 = REGLINK not used in regalloc
	arm64.REGSP, // R31

	arm64.REG_F0,
	arm64.REG_F1,
	arm64.REG_F2,
	arm64.REG_F3,
	arm64.REG_F4,
	arm64.REG_F5,
	arm64.REG_F6,
	arm64.REG_F7,
	arm64.REG_F8,
	arm64.REG_F9,
	arm64.REG_F10,
	arm64.REG_F11,
	arm64.REG_F12,
	arm64.REG_F13,
	arm64.REG_F14,
	arm64.REG_F15,
	arm64.REG_F16,
	arm64.REG_F17,
	arm64.REG_F18,
	arm64.REG_F19,
	arm64.REG_F20,
	arm64.REG_F21,
	arm64.REG_F22,
	arm64.REG_F23,
	arm64.REG_F24,
	arm64.REG_F25,
	arm64.REG_F26,
	arm64.REG_F27,
	arm64.REG_F28,
	arm64.REG_F29,
	arm64.REG_F30,
	arm64.REG_F31,

	0, // SB isn't a real register.  We fill an Addr.Reg field with 0 in this case.
}

// Smallest possible faulting page at address zero,
// see ../../../../runtime/mheap.go:/minPhysPageSize
const minZeroPage = 4096

// loadByType returns the load instruction of the given type.
func loadByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return arm64.AFMOVS
		case 8:
			return arm64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return arm64.AMOVB
			} else {
				return arm64.AMOVBU
			}
		case 2:
			if t.IsSigned() {
				return arm64.AMOVH
			} else {
				return arm64.AMOVHU
			}
		case 4:
			if t.IsSigned() {
				return arm64.AMOVW
			} else {
				return arm64.AMOVWU
			}
		case 8:
			return arm64.AMOVD
		}
	}
	panic("bad load type")
}

// storeByType returns the store instruction of the given type.
func storeByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return arm64.AFMOVS
		case 8:
			return arm64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			return arm64.AMOVB
		case 2:
			return arm64.AMOVH
		case 4:
			return arm64.AMOVW
		case 8:
			return arm64.AMOVD
		}
	}
	panic("bad store type")
}

// makeshift encodes a register shifted by a constant, used as an Offset in Prog
func makeshift(reg int16, typ int64, s int64) int64 {
	return int64(reg&31)<<16 | typ | (s&63)<<10
}

// genshift generates a Prog for r = r0 op (r1 shifted by s)
func genshift(as obj.As, r0, r1, r int16, typ int64, s int64) *obj.Prog {
	p := gc.Prog(as)
	p.From.Type = obj.TYPE_SHIFT
	p.From.Offset = makeshift(r1, typ, s)
	p.Reg = r0
	if r != 0 {
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	}
	return p
}

func ssaGenValue(s *gc.SSAGenState, v *ssa.Value) {
	s.SetLineno(v.Line)
	switch v.Op {
	case ssa.OpInitMem:
		// memory arg needs no code
	case ssa.OpArg:
		// input args need no code
	case ssa.OpSP, ssa.OpSB, ssa.OpGetG:
		// nothing to do
	case ssa.OpCopy, ssa.OpARM64MOVDconvert, ssa.OpARM64MOVDreg:
		if v.Type.IsMemory() {
			return
		}
		x := gc.SSARegNum(v.Args[0])
		y := gc.SSARegNum(v)
		if x == y {
			return
		}
		as := arm64.AMOVD
		if v.Type.IsFloat() {
			switch v.Type.Size() {
			case 4:
				as = arm64.AFMOVS
			case 8:
				as = arm64.AFMOVD
			default:
				panic("bad float size")
			}
		}
		p := gc.Prog(as)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x
		p.To.Type = obj.TYPE_REG
		p.To.Reg = y
	case ssa.OpARM64MOVDnop:
		if gc.SSARegNum(v) != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}
		// nothing to do
	case ssa.OpLoadReg:
		if v.Type.IsFlags() {
			v.Unimplementedf("load flags not implemented: %v", v.LongString())
			return
		}
		p := gc.Prog(loadByType(v.Type))
		n, off := gc.AutoVar(v.Args[0])
		p.From.Type = obj.TYPE_MEM
		p.From.Node = n
		p.From.Sym = gc.Linksym(n.Sym)
		p.From.Offset = off
		if n.Class == gc.PPARAM || n.Class == gc.PPARAMOUT {
			p.From.Name = obj.NAME_PARAM
			p.From.Offset += n.Xoffset
		} else {
			p.From.Name = obj.NAME_AUTO
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpPhi:
		gc.CheckLoweredPhi(v)
	case ssa.OpStoreReg:
		if v.Type.IsFlags() {
			v.Unimplementedf("store flags not implemented: %v", v.LongString())
			return
		}
		p := gc.Prog(storeByType(v.Type))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		n, off := gc.AutoVar(v)
		p.To.Type = obj.TYPE_MEM
		p.To.Node = n
		p.To.Sym = gc.Linksym(n.Sym)
		p.To.Offset = off
		if n.Class == gc.PPARAM || n.Class == gc.PPARAMOUT {
			p.To.Name = obj.NAME_PARAM
			p.To.Offset += n.Xoffset
		} else {
			p.To.Name = obj.NAME_AUTO
		}
	case ssa.OpARM64ADD,
		ssa.OpARM64SUB,
		ssa.OpARM64AND,
		ssa.OpARM64OR,
		ssa.OpARM64XOR,
		ssa.OpARM64BIC,
		ssa.OpARM64MUL,
		ssa.OpARM64MULW,
		ssa.OpARM64MULH,
		ssa.OpARM64UMULH,
		ssa.OpARM64MULL,
		ssa.OpARM64UMULL,
		ssa.OpARM64DIV,
		ssa.OpARM64UDIV,
		ssa.OpARM64DIVW,
		ssa.OpARM64UDIVW,
		ssa.OpARM64MOD,
		ssa.OpARM64UMOD,
		ssa.OpARM64MODW,
		ssa.OpARM64UMODW,
		ssa.OpARM64SLL,
		ssa.OpARM64SRL,
		ssa.OpARM64SRA,
		ssa.OpARM64FADDS,
		ssa.OpARM64FADDD,
		ssa.OpARM64FSUBS,
		ssa.OpARM64FSUBD,
		ssa.OpARM64FMULS,
		ssa.OpARM64FMULD,
		ssa.OpARM64FDIVS,
		ssa.OpARM64FDIVD:
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r2
		p.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpARM64ADDconst,
		ssa.OpARM64SUBconst,
		ssa.OpARM64ANDconst,
		ssa.OpARM64ORconst,
		ssa.OpARM64XORconst,
		ssa.OpARM64BICconst,
		ssa.OpARM64SLLconst,
		ssa.OpARM64SRLconst,
		ssa.OpARM64SRAconst,
		ssa.OpARM64RORconst,
		ssa.OpARM64RORWconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64ADDshiftLL,
		ssa.OpARM64SUBshiftLL,
		ssa.OpARM64ANDshiftLL,
		ssa.OpARM64ORshiftLL,
		ssa.OpARM64XORshiftLL,
		ssa.OpARM64BICshiftLL:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), gc.SSARegNum(v), arm64.SHIFT_LL, v.AuxInt)
	case ssa.OpARM64ADDshiftRL,
		ssa.OpARM64SUBshiftRL,
		ssa.OpARM64ANDshiftRL,
		ssa.OpARM64ORshiftRL,
		ssa.OpARM64XORshiftRL,
		ssa.OpARM64BICshiftRL:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), gc.SSARegNum(v), arm64.SHIFT_LR, v.AuxInt)
	case ssa.OpARM64ADDshiftRA,
		ssa.OpARM64SUBshiftRA,
		ssa.OpARM64ANDshiftRA,
		ssa.OpARM64ORshiftRA,
		ssa.OpARM64XORshiftRA,
		ssa.OpARM64BICshiftRA:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), gc.SSARegNum(v), arm64.SHIFT_AR, v.AuxInt)
	case ssa.OpARM64MOVDconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64FMOVSconst,
		ssa.OpARM64FMOVDconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_FCONST
		p.From.Val = math.Float64frombits(uint64(v.AuxInt))
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64CMP,
		ssa.OpARM64CMPW,
		ssa.OpARM64CMN,
		ssa.OpARM64CMNW,
		ssa.OpARM64FCMPS,
		ssa.OpARM64FCMPD:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.Reg = gc.SSARegNum(v.Args[0])
	case ssa.OpARM64CMPconst,
		ssa.OpARM64CMPWconst,
		ssa.OpARM64CMNconst,
		ssa.OpARM64CMNWconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.Reg = gc.SSARegNum(v.Args[0])
	case ssa.OpARM64CMPshiftLL:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), 0, arm64.SHIFT_LL, v.AuxInt)
	case ssa.OpARM64CMPshiftRL:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), 0, arm64.SHIFT_LR, v.AuxInt)
	case ssa.OpARM64CMPshiftRA:
		genshift(v.Op.Asm(), gc.SSARegNum(v.Args[0]), gc.SSARegNum(v.Args[1]), 0, arm64.SHIFT_AR, v.AuxInt)
	case ssa.OpARM64MOVDaddr:
		p := gc.Prog(arm64.AMOVD)
		p.From.Type = obj.TYPE_ADDR
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

		var wantreg string
		// MOVD $sym+off(base), R
		// the assembler expands it as the following:
		// - base is SP: add constant offset to SP (R13)
		//               when constant is large, tmp register (R11) may be used
		// - base is SB: load external address from constant pool (use relocation)
		switch v.Aux.(type) {
		default:
			v.Fatalf("aux is of unknown type %T", v.Aux)
		case *ssa.ExternSymbol:
			wantreg = "SB"
			gc.AddAux(&p.From, v)
		case *ssa.ArgSymbol, *ssa.AutoSymbol:
			wantreg = "SP"
			gc.AddAux(&p.From, v)
		case nil:
			// No sym, just MOVD $off(SP), R
			wantreg = "SP"
			p.From.Reg = arm64.REGSP
			p.From.Offset = v.AuxInt
		}
		if reg := gc.SSAReg(v.Args[0]); reg.Name() != wantreg {
			v.Fatalf("bad reg %s for symbol type %T, want %s", reg.Name(), v.Aux, wantreg)
		}
	case ssa.OpARM64MOVBload,
		ssa.OpARM64MOVBUload,
		ssa.OpARM64MOVHload,
		ssa.OpARM64MOVHUload,
		ssa.OpARM64MOVWload,
		ssa.OpARM64MOVWUload,
		ssa.OpARM64MOVDload,
		ssa.OpARM64FMOVSload,
		ssa.OpARM64FMOVDload:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64LDAR,
		ssa.OpARM64LDARW:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum0(v)
	case ssa.OpARM64MOVBstore,
		ssa.OpARM64MOVHstore,
		ssa.OpARM64MOVWstore,
		ssa.OpARM64MOVDstore,
		ssa.OpARM64FMOVSstore,
		ssa.OpARM64FMOVDstore,
		ssa.OpARM64STLR,
		ssa.OpARM64STLRW:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpARM64MOVBstorezero,
		ssa.OpARM64MOVHstorezero,
		ssa.OpARM64MOVWstorezero,
		ssa.OpARM64MOVDstorezero:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = arm64.REGZERO
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpARM64LoweredAtomicExchange64,
		ssa.OpARM64LoweredAtomicExchange32:
		// LDAXR	(Rarg0), Rout
		// STLXR	Rarg1, (Rarg0), Rtmp
		// CBNZ		Rtmp, -2(PC)
		ld := arm64.ALDAXR
		st := arm64.ASTLXR
		if v.Op == ssa.OpARM64LoweredAtomicExchange32 {
			ld = arm64.ALDAXRW
			st = arm64.ASTLXRW
		}
		r0 := gc.SSARegNum(v.Args[0])
		r1 := gc.SSARegNum(v.Args[1])
		out := gc.SSARegNum0(v)
		p := gc.Prog(ld)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = r0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = out
		p1 := gc.Prog(st)
		p1.From.Type = obj.TYPE_REG
		p1.From.Reg = r1
		p1.To.Type = obj.TYPE_MEM
		p1.To.Reg = r0
		p1.RegTo2 = arm64.REGTMP
		p2 := gc.Prog(arm64.ACBNZ)
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = arm64.REGTMP
		p2.To.Type = obj.TYPE_BRANCH
		gc.Patch(p2, p)
	case ssa.OpARM64LoweredAtomicAdd64,
		ssa.OpARM64LoweredAtomicAdd32:
		// LDAXR	(Rarg0), Rout
		// ADD		Rarg1, Rout
		// STLXR	Rout, (Rarg0), Rtmp
		// CBNZ		Rtmp, -3(PC)
		ld := arm64.ALDAXR
		st := arm64.ASTLXR
		if v.Op == ssa.OpARM64LoweredAtomicAdd32 {
			ld = arm64.ALDAXRW
			st = arm64.ASTLXRW
		}
		r0 := gc.SSARegNum(v.Args[0])
		r1 := gc.SSARegNum(v.Args[1])
		out := gc.SSARegNum0(v)
		p := gc.Prog(ld)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = r0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = out
		p1 := gc.Prog(arm64.AADD)
		p1.From.Type = obj.TYPE_REG
		p1.From.Reg = r1
		p1.To.Type = obj.TYPE_REG
		p1.To.Reg = out
		p2 := gc.Prog(st)
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = out
		p2.To.Type = obj.TYPE_MEM
		p2.To.Reg = r0
		p2.RegTo2 = arm64.REGTMP
		p3 := gc.Prog(arm64.ACBNZ)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = arm64.REGTMP
		p3.To.Type = obj.TYPE_BRANCH
		gc.Patch(p3, p)
	case ssa.OpARM64LoweredAtomicCas64,
		ssa.OpARM64LoweredAtomicCas32:
		// LDAXR	(Rarg0), Rtmp
		// CMP		Rarg1, Rtmp
		// BNE		3(PC)
		// STLXR	Rarg2, (Rarg0), Rtmp
		// CBNZ		Rtmp, -4(PC)
		// CSET		EQ, Rout
		ld := arm64.ALDAXR
		st := arm64.ASTLXR
		cmp := arm64.ACMP
		if v.Op == ssa.OpARM64LoweredAtomicCas32 {
			ld = arm64.ALDAXRW
			st = arm64.ASTLXRW
			cmp = arm64.ACMPW
		}
		r0 := gc.SSARegNum(v.Args[0])
		r1 := gc.SSARegNum(v.Args[1])
		r2 := gc.SSARegNum(v.Args[2])
		out := gc.SSARegNum0(v)
		p := gc.Prog(ld)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = r0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = arm64.REGTMP
		p1 := gc.Prog(cmp)
		p1.From.Type = obj.TYPE_REG
		p1.From.Reg = r1
		p1.Reg = arm64.REGTMP
		p2 := gc.Prog(arm64.ABNE)
		p2.To.Type = obj.TYPE_BRANCH
		p3 := gc.Prog(st)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = r2
		p3.To.Type = obj.TYPE_MEM
		p3.To.Reg = r0
		p3.RegTo2 = arm64.REGTMP
		p4 := gc.Prog(arm64.ACBNZ)
		p4.From.Type = obj.TYPE_REG
		p4.From.Reg = arm64.REGTMP
		p4.To.Type = obj.TYPE_BRANCH
		gc.Patch(p4, p)
		p5 := gc.Prog(arm64.ACSET)
		p5.From.Type = obj.TYPE_REG // assembler encodes conditional bits in Reg
		p5.From.Reg = arm64.COND_EQ
		p5.To.Type = obj.TYPE_REG
		p5.To.Reg = out
		gc.Patch(p2, p5)
	case ssa.OpARM64LoweredAtomicAnd8,
		ssa.OpARM64LoweredAtomicOr8:
		// LDAXRB	(Rarg0), Rtmp
		// AND/OR	Rarg1, Rtmp
		// STLXRB	Rtmp, (Rarg0), Rtmp
		// CBNZ		Rtmp, -3(PC)
		r0 := gc.SSARegNum(v.Args[0])
		r1 := gc.SSARegNum(v.Args[1])
		p := gc.Prog(arm64.ALDAXRB)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = r0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = arm64.REGTMP
		p1 := gc.Prog(v.Op.Asm())
		p1.From.Type = obj.TYPE_REG
		p1.From.Reg = r1
		p1.To.Type = obj.TYPE_REG
		p1.To.Reg = arm64.REGTMP
		p2 := gc.Prog(arm64.ASTLXRB)
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = arm64.REGTMP
		p2.To.Type = obj.TYPE_MEM
		p2.To.Reg = r0
		p2.RegTo2 = arm64.REGTMP
		p3 := gc.Prog(arm64.ACBNZ)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = arm64.REGTMP
		p3.To.Type = obj.TYPE_BRANCH
		gc.Patch(p3, p)
	case ssa.OpARM64MOVBreg,
		ssa.OpARM64MOVBUreg,
		ssa.OpARM64MOVHreg,
		ssa.OpARM64MOVHUreg,
		ssa.OpARM64MOVWreg,
		ssa.OpARM64MOVWUreg:
		a := v.Args[0]
		for a.Op == ssa.OpCopy || a.Op == ssa.OpARM64MOVDreg {
			a = a.Args[0]
		}
		if a.Op == ssa.OpLoadReg {
			t := a.Type
			switch {
			case v.Op == ssa.OpARM64MOVBreg && t.Size() == 1 && t.IsSigned(),
				v.Op == ssa.OpARM64MOVBUreg && t.Size() == 1 && !t.IsSigned(),
				v.Op == ssa.OpARM64MOVHreg && t.Size() == 2 && t.IsSigned(),
				v.Op == ssa.OpARM64MOVHUreg && t.Size() == 2 && !t.IsSigned(),
				v.Op == ssa.OpARM64MOVWreg && t.Size() == 4 && t.IsSigned(),
				v.Op == ssa.OpARM64MOVWUreg && t.Size() == 4 && !t.IsSigned():
				// arg is a proper-typed load, already zero/sign-extended, don't extend again
				if gc.SSARegNum(v) == gc.SSARegNum(v.Args[0]) {
					return
				}
				p := gc.Prog(arm64.AMOVD)
				p.From.Type = obj.TYPE_REG
				p.From.Reg = gc.SSARegNum(v.Args[0])
				p.To.Type = obj.TYPE_REG
				p.To.Reg = gc.SSARegNum(v)
				return
			default:
			}
		}
		fallthrough
	case ssa.OpARM64MVN,
		ssa.OpARM64NEG,
		ssa.OpARM64FNEGS,
		ssa.OpARM64FNEGD,
		ssa.OpARM64FSQRTD,
		ssa.OpARM64FCVTZSSW,
		ssa.OpARM64FCVTZSDW,
		ssa.OpARM64FCVTZUSW,
		ssa.OpARM64FCVTZUDW,
		ssa.OpARM64FCVTZSS,
		ssa.OpARM64FCVTZSD,
		ssa.OpARM64FCVTZUS,
		ssa.OpARM64FCVTZUD,
		ssa.OpARM64SCVTFWS,
		ssa.OpARM64SCVTFWD,
		ssa.OpARM64SCVTFS,
		ssa.OpARM64SCVTFD,
		ssa.OpARM64UCVTFWS,
		ssa.OpARM64UCVTFWD,
		ssa.OpARM64UCVTFS,
		ssa.OpARM64UCVTFD,
		ssa.OpARM64FCVTSD,
		ssa.OpARM64FCVTDS,
		ssa.OpARM64REV,
		ssa.OpARM64REVW,
		ssa.OpARM64REV16W,
		ssa.OpARM64RBIT,
		ssa.OpARM64RBITW,
		ssa.OpARM64CLZ,
		ssa.OpARM64CLZW:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64CSELULT,
		ssa.OpARM64CSELULT0:
		r1 := int16(arm64.REGZERO)
		if v.Op == ssa.OpARM64CSELULT {
			r1 = gc.SSARegNum(v.Args[1])
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG // assembler encodes conditional bits in Reg
		p.From.Reg = arm64.COND_LO
		p.Reg = gc.SSARegNum(v.Args[0])
		p.From3 = &obj.Addr{Type: obj.TYPE_REG, Reg: r1}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpARM64DUFFZERO:
		// runtime.duffzero expects start address - 8 in R16
		p := gc.Prog(arm64.ASUB)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 8
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = arm64.REG_R16
		p = gc.Prog(obj.ADUFFZERO)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Pkglookup("duffzero", gc.Runtimepkg))
		p.To.Offset = v.AuxInt
	case ssa.OpARM64LoweredZero:
		// MOVD.P	ZR, 8(R16)
		// CMP	Rarg1, R16
		// BLE	-2(PC)
		// arg1 is the address of the last element to zero
		p := gc.Prog(arm64.AMOVD)
		p.Scond = arm64.C_XPOST
		p.From.Type = obj.TYPE_REG
		p.From.Reg = arm64.REGZERO
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = arm64.REG_R16
		p.To.Offset = 8
		p2 := gc.Prog(arm64.ACMP)
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = gc.SSARegNum(v.Args[1])
		p2.Reg = arm64.REG_R16
		p3 := gc.Prog(arm64.ABLE)
		p3.To.Type = obj.TYPE_BRANCH
		gc.Patch(p3, p)
	case ssa.OpARM64LoweredMove:
		// MOVD.P	8(R16), Rtmp
		// MOVD.P	Rtmp, 8(R17)
		// CMP	Rarg2, R16
		// BLE	-3(PC)
		// arg2 is the address of the last element of src
		p := gc.Prog(arm64.AMOVD)
		p.Scond = arm64.C_XPOST
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = arm64.REG_R16
		p.From.Offset = 8
		p.To.Type = obj.TYPE_REG
		p.To.Reg = arm64.REGTMP
		p2 := gc.Prog(arm64.AMOVD)
		p2.Scond = arm64.C_XPOST
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = arm64.REGTMP
		p2.To.Type = obj.TYPE_MEM
		p2.To.Reg = arm64.REG_R17
		p2.To.Offset = 8
		p3 := gc.Prog(arm64.ACMP)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = gc.SSARegNum(v.Args[2])
		p3.Reg = arm64.REG_R16
		p4 := gc.Prog(arm64.ABLE)
		p4.To.Type = obj.TYPE_BRANCH
		gc.Patch(p4, p)
	case ssa.OpARM64CALLstatic:
		if v.Aux.(*gc.Sym) == gc.Deferreturn.Sym {
			// Deferred calls will appear to be returning to
			// the CALL deferreturn(SB) that we are about to emit.
			// However, the stack trace code will show the line
			// of the instruction byte before the return PC.
			// To avoid that being an unrelated instruction,
			// insert an actual hardware NOP that will have the right line number.
			// This is different from obj.ANOP, which is a virtual no-op
			// that doesn't make it into the instruction stream.
			ginsnop()
		}
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(v.Aux.(*gc.Sym))
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpARM64CALLclosure:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpARM64CALLdefer:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Deferproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpARM64CALLgo:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Newproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpARM64CALLinter:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpARM64LoweredNilCheck:
		// Optimization - if the subsequent block has a load or store
		// at the same address, we don't need to issue this instruction.
		mem := v.Args[1]
		for _, w := range v.Block.Succs[0].Block().Values {
			if w.Op == ssa.OpPhi {
				if w.Type.IsMemory() {
					mem = w
				}
				continue
			}
			if len(w.Args) == 0 || !w.Args[len(w.Args)-1].Type.IsMemory() {
				// w doesn't use a store - can't be a memory op.
				continue
			}
			if w.Args[len(w.Args)-1] != mem {
				v.Fatalf("wrong store after nilcheck v=%s w=%s", v, w)
			}
			switch w.Op {
			case ssa.OpARM64MOVBload, ssa.OpARM64MOVBUload, ssa.OpARM64MOVHload, ssa.OpARM64MOVHUload,
				ssa.OpARM64MOVWload, ssa.OpARM64MOVWUload, ssa.OpARM64MOVDload,
				ssa.OpARM64FMOVSload, ssa.OpARM64FMOVDload,
				ssa.OpARM64LDAR, ssa.OpARM64LDARW,
				ssa.OpARM64MOVBstore, ssa.OpARM64MOVHstore, ssa.OpARM64MOVWstore, ssa.OpARM64MOVDstore,
				ssa.OpARM64FMOVSstore, ssa.OpARM64FMOVDstore,
				ssa.OpARM64MOVBstorezero, ssa.OpARM64MOVHstorezero, ssa.OpARM64MOVWstorezero, ssa.OpARM64MOVDstorezero,
				ssa.OpARM64STLR, ssa.OpARM64STLRW,
				ssa.OpARM64LoweredAtomicExchange64, ssa.OpARM64LoweredAtomicExchange32,
				ssa.OpARM64LoweredAtomicAdd64, ssa.OpARM64LoweredAtomicAdd32,
				ssa.OpARM64LoweredAtomicCas64, ssa.OpARM64LoweredAtomicCas32,
				ssa.OpARM64LoweredAtomicAnd8, ssa.OpARM64LoweredAtomicOr8:
				// arg0 is ptr, auxint is offset (atomic ops have auxint 0)
				if w.Args[0] == v.Args[0] && w.Aux == nil && w.AuxInt >= 0 && w.AuxInt < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpARM64DUFFZERO, ssa.OpARM64LoweredZero:
				// arg0 is ptr
				if w.Args[0] == v.Args[0] {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpARM64LoweredMove:
				// arg0 is dst ptr, arg1 is src ptr
				if w.Args[0] == v.Args[0] || w.Args[1] == v.Args[0] {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			default:
			}
			if w.Type.IsMemory() || w.Type.IsTuple() && w.Type.FieldType(1).IsMemory() {
				if w.Op == ssa.OpVarDef || w.Op == ssa.OpVarKill || w.Op == ssa.OpVarLive {
					// these ops are OK
					mem = w
					continue
				}
				// We can't delay the nil check past the next store.
				break
			}
		}
		// Issue a load which will fault if arg is nil.
		p := gc.Prog(arm64.AMOVB)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = arm64.REGTMP
		if gc.Debug_checknil != 0 && v.Line > 1 { // v.Line==1 in generated wrappers
			gc.Warnl(v.Line, "generated nil check")
		}
	case ssa.OpVarDef:
		gc.Gvardef(v.Aux.(*gc.Node))
	case ssa.OpVarKill:
		gc.Gvarkill(v.Aux.(*gc.Node))
	case ssa.OpVarLive:
		gc.Gvarlive(v.Aux.(*gc.Node))
	case ssa.OpKeepAlive:
		gc.KeepAlive(v)
	case ssa.OpARM64Equal,
		ssa.OpARM64NotEqual,
		ssa.OpARM64LessThan,
		ssa.OpARM64LessEqual,
		ssa.OpARM64GreaterThan,
		ssa.OpARM64GreaterEqual,
		ssa.OpARM64LessThanU,
		ssa.OpARM64LessEqualU,
		ssa.OpARM64GreaterThanU,
		ssa.OpARM64GreaterEqualU:
		// generate boolean values using CSET
		p := gc.Prog(arm64.ACSET)
		p.From.Type = obj.TYPE_REG // assembler encodes conditional bits in Reg
		p.From.Reg = condBits[v.Op]
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpSelect0, ssa.OpSelect1:
		// nothing to do
	case ssa.OpARM64LoweredGetClosurePtr:
		// Closure pointer is R26 (arm64.REGCTXT).
		gc.CheckLoweredGetClosurePtr(v)
	case ssa.OpARM64FlagEQ,
		ssa.OpARM64FlagLT_ULT,
		ssa.OpARM64FlagLT_UGT,
		ssa.OpARM64FlagGT_ULT,
		ssa.OpARM64FlagGT_UGT:
		v.Fatalf("Flag* ops should never make it to codegen %v", v.LongString())
	case ssa.OpARM64InvertFlags:
		v.Fatalf("InvertFlags should never make it to codegen %v", v.LongString())
	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
}

var condBits = map[ssa.Op]int16{
	ssa.OpARM64Equal:         arm64.COND_EQ,
	ssa.OpARM64NotEqual:      arm64.COND_NE,
	ssa.OpARM64LessThan:      arm64.COND_LT,
	ssa.OpARM64LessThanU:     arm64.COND_LO,
	ssa.OpARM64LessEqual:     arm64.COND_LE,
	ssa.OpARM64LessEqualU:    arm64.COND_LS,
	ssa.OpARM64GreaterThan:   arm64.COND_GT,
	ssa.OpARM64GreaterThanU:  arm64.COND_HI,
	ssa.OpARM64GreaterEqual:  arm64.COND_GE,
	ssa.OpARM64GreaterEqualU: arm64.COND_HS,
}

var blockJump = map[ssa.BlockKind]struct {
	asm, invasm obj.As
}{
	ssa.BlockARM64EQ:  {arm64.ABEQ, arm64.ABNE},
	ssa.BlockARM64NE:  {arm64.ABNE, arm64.ABEQ},
	ssa.BlockARM64LT:  {arm64.ABLT, arm64.ABGE},
	ssa.BlockARM64GE:  {arm64.ABGE, arm64.ABLT},
	ssa.BlockARM64LE:  {arm64.ABLE, arm64.ABGT},
	ssa.BlockARM64GT:  {arm64.ABGT, arm64.ABLE},
	ssa.BlockARM64ULT: {arm64.ABLO, arm64.ABHS},
	ssa.BlockARM64UGE: {arm64.ABHS, arm64.ABLO},
	ssa.BlockARM64UGT: {arm64.ABHI, arm64.ABLS},
	ssa.BlockARM64ULE: {arm64.ABLS, arm64.ABHI},
}

func ssaGenBlock(s *gc.SSAGenState, b, next *ssa.Block) {
	s.SetLineno(b.Line)

	switch b.Kind {
	case ssa.BlockPlain, ssa.BlockCheck:
		if b.Succs[0].Block() != next {
			p := gc.Prog(obj.AJMP)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		}

	case ssa.BlockDefer:
		// defer returns in R0:
		// 0 if we should continue executing
		// 1 if we should jump to deferreturn call
		p := gc.Prog(arm64.ACMP)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 0
		p.Reg = arm64.REG_R0
		p = gc.Prog(arm64.ABNE)
		p.To.Type = obj.TYPE_BRANCH
		s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[1].Block()})
		if b.Succs[0].Block() != next {
			p := gc.Prog(obj.AJMP)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		}

	case ssa.BlockExit:
		gc.Prog(obj.AUNDEF) // tell plive.go that we never reach here

	case ssa.BlockRet:
		gc.Prog(obj.ARET)

	case ssa.BlockRetJmp:
		p := gc.Prog(obj.ARET)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(b.Aux.(*gc.Sym))

	case ssa.BlockARM64EQ, ssa.BlockARM64NE,
		ssa.BlockARM64LT, ssa.BlockARM64GE,
		ssa.BlockARM64LE, ssa.BlockARM64GT,
		ssa.BlockARM64ULT, ssa.BlockARM64UGT,
		ssa.BlockARM64ULE, ssa.BlockARM64UGE:
		jmp := blockJump[b.Kind]
		var p *obj.Prog
		switch next {
		case b.Succs[0].Block():
			p = gc.Prog(jmp.invasm)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[1].Block()})
		case b.Succs[1].Block():
			p = gc.Prog(jmp.asm)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		default:
			p = gc.Prog(jmp.asm)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
			q := gc.Prog(obj.AJMP)
			q.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: q, B: b.Succs[1].Block()})
		}

	default:
		b.Unimplementedf("branch not implemented: %s. Control: %s", b.LongString(), b.Control.LongString())
	}
}

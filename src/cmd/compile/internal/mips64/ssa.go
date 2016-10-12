// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mips64

import (
	"math"

	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/mips"
)

var ssaRegToReg = []int16{
	mips.REG_R0, // constant 0
	mips.REG_R1,
	mips.REG_R2,
	mips.REG_R3,
	mips.REG_R4,
	mips.REG_R5,
	mips.REG_R6,
	mips.REG_R7,
	mips.REG_R8,
	mips.REG_R9,
	mips.REG_R10,
	mips.REG_R11,
	mips.REG_R12,
	mips.REG_R13,
	mips.REG_R14,
	mips.REG_R15,
	mips.REG_R16,
	mips.REG_R17,
	mips.REG_R18,
	mips.REG_R19,
	mips.REG_R20,
	mips.REG_R21,
	mips.REG_R22,
	// R23 = REGTMP not used in regalloc
	mips.REG_R24,
	mips.REG_R25,
	// R26 reserved by kernel
	// R27 reserved by kernel
	// R28 = REGSB not used in regalloc
	mips.REGSP, // R29
	mips.REGG,  // R30
	// R31 = REGLINK not used in regalloc

	mips.REG_F0,
	mips.REG_F1,
	mips.REG_F2,
	mips.REG_F3,
	mips.REG_F4,
	mips.REG_F5,
	mips.REG_F6,
	mips.REG_F7,
	mips.REG_F8,
	mips.REG_F9,
	mips.REG_F10,
	mips.REG_F11,
	mips.REG_F12,
	mips.REG_F13,
	mips.REG_F14,
	mips.REG_F15,
	mips.REG_F16,
	mips.REG_F17,
	mips.REG_F18,
	mips.REG_F19,
	mips.REG_F20,
	mips.REG_F21,
	mips.REG_F22,
	mips.REG_F23,
	mips.REG_F24,
	mips.REG_F25,
	mips.REG_F26,
	mips.REG_F27,
	mips.REG_F28,
	mips.REG_F29,
	mips.REG_F30,
	mips.REG_F31,

	mips.REG_HI, // high bits of multiplication
	mips.REG_LO, // low bits of multiplication

	0, // SB isn't a real register.  We fill an Addr.Reg field with 0 in this case.
}

// Smallest possible faulting page at address zero,
// see ../../../../runtime/mheap.go:/minPhysPageSize
const minZeroPage = 4096

// isFPreg returns whether r is an FP register
func isFPreg(r int16) bool {
	return mips.REG_F0 <= r && r <= mips.REG_F31
}

// isHILO returns whether r is HI or LO register
func isHILO(r int16) bool {
	return r == mips.REG_HI || r == mips.REG_LO
}

// loadByType returns the load instruction of the given type.
func loadByType(t ssa.Type, r int16) obj.As {
	if isFPreg(r) {
		if t.Size() == 4 { // float32 or int32
			return mips.AMOVF
		} else { // float64 or int64
			return mips.AMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return mips.AMOVB
			} else {
				return mips.AMOVBU
			}
		case 2:
			if t.IsSigned() {
				return mips.AMOVH
			} else {
				return mips.AMOVHU
			}
		case 4:
			if t.IsSigned() {
				return mips.AMOVW
			} else {
				return mips.AMOVWU
			}
		case 8:
			return mips.AMOVV
		}
	}
	panic("bad load type")
}

// storeByType returns the store instruction of the given type.
func storeByType(t ssa.Type, r int16) obj.As {
	if isFPreg(r) {
		if t.Size() == 4 { // float32 or int32
			return mips.AMOVF
		} else { // float64 or int64
			return mips.AMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			return mips.AMOVB
		case 2:
			return mips.AMOVH
		case 4:
			return mips.AMOVW
		case 8:
			return mips.AMOVV
		}
	}
	panic("bad store type")
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
	case ssa.OpCopy, ssa.OpMIPS64MOVVconvert, ssa.OpMIPS64MOVVreg:
		if v.Type.IsMemory() {
			return
		}
		x := gc.SSARegNum(v.Args[0])
		y := gc.SSARegNum(v)
		if x == y {
			return
		}
		as := mips.AMOVV
		if isFPreg(x) && isFPreg(y) {
			as = mips.AMOVD
		}
		p := gc.Prog(as)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x
		p.To.Type = obj.TYPE_REG
		p.To.Reg = y
		if isHILO(x) && isHILO(y) || isHILO(x) && isFPreg(y) || isFPreg(x) && isHILO(y) {
			// cannot move between special registers, use TMP as intermediate
			p.To.Reg = mips.REGTMP
			p = gc.Prog(mips.AMOVV)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = mips.REGTMP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = y
		}
	case ssa.OpMIPS64MOVVnop:
		if gc.SSARegNum(v) != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}
		// nothing to do
	case ssa.OpLoadReg:
		if v.Type.IsFlags() {
			v.Unimplementedf("load flags not implemented: %v", v.LongString())
			return
		}
		r := gc.SSARegNum(v)
		p := gc.Prog(loadByType(v.Type, r))
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
		p.To.Reg = r
		if isHILO(r) {
			// cannot directly load, load to TMP and move
			p.To.Reg = mips.REGTMP
			p = gc.Prog(mips.AMOVV)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = mips.REGTMP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = r
		}
	case ssa.OpPhi:
		gc.CheckLoweredPhi(v)
	case ssa.OpStoreReg:
		if v.Type.IsFlags() {
			v.Unimplementedf("store flags not implemented: %v", v.LongString())
			return
		}
		r := gc.SSARegNum(v.Args[0])
		if isHILO(r) {
			// cannot directly store, move to TMP and store
			p := gc.Prog(mips.AMOVV)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r
			p.To.Type = obj.TYPE_REG
			p.To.Reg = mips.REGTMP
			r = mips.REGTMP
		}
		p := gc.Prog(storeByType(v.Type, r))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r
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
	case ssa.OpMIPS64ADDV,
		ssa.OpMIPS64SUBV,
		ssa.OpMIPS64AND,
		ssa.OpMIPS64OR,
		ssa.OpMIPS64XOR,
		ssa.OpMIPS64NOR,
		ssa.OpMIPS64SLLV,
		ssa.OpMIPS64SRLV,
		ssa.OpMIPS64SRAV,
		ssa.OpMIPS64ADDF,
		ssa.OpMIPS64ADDD,
		ssa.OpMIPS64SUBF,
		ssa.OpMIPS64SUBD,
		ssa.OpMIPS64MULF,
		ssa.OpMIPS64MULD,
		ssa.OpMIPS64DIVF,
		ssa.OpMIPS64DIVD:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64SGT,
		ssa.OpMIPS64SGTU:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64ADDVconst,
		ssa.OpMIPS64SUBVconst,
		ssa.OpMIPS64ANDconst,
		ssa.OpMIPS64ORconst,
		ssa.OpMIPS64XORconst,
		ssa.OpMIPS64NORconst,
		ssa.OpMIPS64SLLVconst,
		ssa.OpMIPS64SRLVconst,
		ssa.OpMIPS64SRAVconst,
		ssa.OpMIPS64SGTconst,
		ssa.OpMIPS64SGTUconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64MULV,
		ssa.OpMIPS64MULVU,
		ssa.OpMIPS64DIVV,
		ssa.OpMIPS64DIVVU:
		// result in hi,lo
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.Reg = gc.SSARegNum(v.Args[0])
	case ssa.OpMIPS64MOVVconst:
		r := gc.SSARegNum(v)
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
		if isFPreg(r) || isHILO(r) {
			// cannot move into FP or special registers, use TMP as intermediate
			p.To.Reg = mips.REGTMP
			p = gc.Prog(mips.AMOVV)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = mips.REGTMP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = r
		}
	case ssa.OpMIPS64MOVFconst,
		ssa.OpMIPS64MOVDconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_FCONST
		p.From.Val = math.Float64frombits(uint64(v.AuxInt))
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64CMPEQF,
		ssa.OpMIPS64CMPEQD,
		ssa.OpMIPS64CMPGEF,
		ssa.OpMIPS64CMPGED,
		ssa.OpMIPS64CMPGTF,
		ssa.OpMIPS64CMPGTD:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.Reg = gc.SSARegNum(v.Args[1])
	case ssa.OpMIPS64MOVVaddr:
		p := gc.Prog(mips.AMOVV)
		p.From.Type = obj.TYPE_ADDR
		var wantreg string
		// MOVV $sym+off(base), R
		// the assembler expands it as the following:
		// - base is SP: add constant offset to SP (R29)
		//               when constant is large, tmp register (R23) may be used
		// - base is SB: load external address with relocation
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
			// No sym, just MOVV $off(SP), R
			wantreg = "SP"
			p.From.Reg = mips.REGSP
			p.From.Offset = v.AuxInt
		}
		if reg := gc.SSAReg(v.Args[0]); reg.Name() != wantreg {
			v.Fatalf("bad reg %s for symbol type %T, want %s", reg.Name(), v.Aux, wantreg)
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64MOVBload,
		ssa.OpMIPS64MOVBUload,
		ssa.OpMIPS64MOVHload,
		ssa.OpMIPS64MOVHUload,
		ssa.OpMIPS64MOVWload,
		ssa.OpMIPS64MOVWUload,
		ssa.OpMIPS64MOVVload,
		ssa.OpMIPS64MOVFload,
		ssa.OpMIPS64MOVDload:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64MOVBstore,
		ssa.OpMIPS64MOVHstore,
		ssa.OpMIPS64MOVWstore,
		ssa.OpMIPS64MOVVstore,
		ssa.OpMIPS64MOVFstore,
		ssa.OpMIPS64MOVDstore:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpMIPS64MOVBstorezero,
		ssa.OpMIPS64MOVHstorezero,
		ssa.OpMIPS64MOVWstorezero,
		ssa.OpMIPS64MOVVstorezero:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = mips.REGZERO
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpMIPS64MOVBreg,
		ssa.OpMIPS64MOVBUreg,
		ssa.OpMIPS64MOVHreg,
		ssa.OpMIPS64MOVHUreg,
		ssa.OpMIPS64MOVWreg,
		ssa.OpMIPS64MOVWUreg:
		a := v.Args[0]
		for a.Op == ssa.OpCopy || a.Op == ssa.OpMIPS64MOVVreg {
			a = a.Args[0]
		}
		if a.Op == ssa.OpLoadReg {
			t := a.Type
			switch {
			case v.Op == ssa.OpMIPS64MOVBreg && t.Size() == 1 && t.IsSigned(),
				v.Op == ssa.OpMIPS64MOVBUreg && t.Size() == 1 && !t.IsSigned(),
				v.Op == ssa.OpMIPS64MOVHreg && t.Size() == 2 && t.IsSigned(),
				v.Op == ssa.OpMIPS64MOVHUreg && t.Size() == 2 && !t.IsSigned(),
				v.Op == ssa.OpMIPS64MOVWreg && t.Size() == 4 && t.IsSigned(),
				v.Op == ssa.OpMIPS64MOVWUreg && t.Size() == 4 && !t.IsSigned():
				// arg is a proper-typed load, already zero/sign-extended, don't extend again
				if gc.SSARegNum(v) == gc.SSARegNum(v.Args[0]) {
					return
				}
				p := gc.Prog(mips.AMOVV)
				p.From.Type = obj.TYPE_REG
				p.From.Reg = gc.SSARegNum(v.Args[0])
				p.To.Type = obj.TYPE_REG
				p.To.Reg = gc.SSARegNum(v)
				return
			default:
			}
		}
		fallthrough
	case ssa.OpMIPS64MOVWF,
		ssa.OpMIPS64MOVWD,
		ssa.OpMIPS64TRUNCFW,
		ssa.OpMIPS64TRUNCDW,
		ssa.OpMIPS64MOVVF,
		ssa.OpMIPS64MOVVD,
		ssa.OpMIPS64TRUNCFV,
		ssa.OpMIPS64TRUNCDV,
		ssa.OpMIPS64MOVFD,
		ssa.OpMIPS64MOVDF,
		ssa.OpMIPS64NEGF,
		ssa.OpMIPS64NEGD:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64NEGV:
		// SUB from REGZERO
		p := gc.Prog(mips.ASUBVU)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.Reg = mips.REGZERO
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpMIPS64DUFFZERO:
		// runtime.duffzero expects start address - 8 in R1
		p := gc.Prog(mips.ASUBVU)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 8
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = mips.REG_R1
		p = gc.Prog(obj.ADUFFZERO)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Pkglookup("duffzero", gc.Runtimepkg))
		p.To.Offset = v.AuxInt
	case ssa.OpMIPS64LoweredZero:
		// SUBV	$8, R1
		// MOVV	R0, 8(R1)
		// ADDV	$8, R1
		// BNE	Rarg1, R1, -2(PC)
		// arg1 is the address of the last element to zero
		var sz int64
		var mov obj.As
		switch {
		case v.AuxInt%8 == 0:
			sz = 8
			mov = mips.AMOVV
		case v.AuxInt%4 == 0:
			sz = 4
			mov = mips.AMOVW
		case v.AuxInt%2 == 0:
			sz = 2
			mov = mips.AMOVH
		default:
			sz = 1
			mov = mips.AMOVB
		}
		p := gc.Prog(mips.ASUBVU)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = sz
		p.To.Type = obj.TYPE_REG
		p.To.Reg = mips.REG_R1
		p2 := gc.Prog(mov)
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = mips.REGZERO
		p2.To.Type = obj.TYPE_MEM
		p2.To.Reg = mips.REG_R1
		p2.To.Offset = sz
		p3 := gc.Prog(mips.AADDVU)
		p3.From.Type = obj.TYPE_CONST
		p3.From.Offset = sz
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = mips.REG_R1
		p4 := gc.Prog(mips.ABNE)
		p4.From.Type = obj.TYPE_REG
		p4.From.Reg = gc.SSARegNum(v.Args[1])
		p4.Reg = mips.REG_R1
		p4.To.Type = obj.TYPE_BRANCH
		gc.Patch(p4, p2)
	case ssa.OpMIPS64LoweredMove:
		// SUBV	$8, R1
		// MOVV	8(R1), Rtmp
		// MOVV	Rtmp, (R2)
		// ADDV	$8, R1
		// ADDV	$8, R2
		// BNE	Rarg2, R1, -4(PC)
		// arg2 is the address of the last element of src
		var sz int64
		var mov obj.As
		switch {
		case v.AuxInt%8 == 0:
			sz = 8
			mov = mips.AMOVV
		case v.AuxInt%4 == 0:
			sz = 4
			mov = mips.AMOVW
		case v.AuxInt%2 == 0:
			sz = 2
			mov = mips.AMOVH
		default:
			sz = 1
			mov = mips.AMOVB
		}
		p := gc.Prog(mips.ASUBVU)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = sz
		p.To.Type = obj.TYPE_REG
		p.To.Reg = mips.REG_R1
		p2 := gc.Prog(mov)
		p2.From.Type = obj.TYPE_MEM
		p2.From.Reg = mips.REG_R1
		p2.From.Offset = sz
		p2.To.Type = obj.TYPE_REG
		p2.To.Reg = mips.REGTMP
		p3 := gc.Prog(mov)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = mips.REGTMP
		p3.To.Type = obj.TYPE_MEM
		p3.To.Reg = mips.REG_R2
		p4 := gc.Prog(mips.AADDVU)
		p4.From.Type = obj.TYPE_CONST
		p4.From.Offset = sz
		p4.To.Type = obj.TYPE_REG
		p4.To.Reg = mips.REG_R1
		p5 := gc.Prog(mips.AADDVU)
		p5.From.Type = obj.TYPE_CONST
		p5.From.Offset = sz
		p5.To.Type = obj.TYPE_REG
		p5.To.Reg = mips.REG_R2
		p6 := gc.Prog(mips.ABNE)
		p6.From.Type = obj.TYPE_REG
		p6.From.Reg = gc.SSARegNum(v.Args[2])
		p6.Reg = mips.REG_R1
		p6.To.Type = obj.TYPE_BRANCH
		gc.Patch(p6, p2)
	case ssa.OpMIPS64CALLstatic:
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
	case ssa.OpMIPS64CALLclosure:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpMIPS64CALLdefer:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Deferproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpMIPS64CALLgo:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Newproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpMIPS64CALLinter:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpMIPS64LoweredNilCheck:
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
			case ssa.OpMIPS64MOVBload, ssa.OpMIPS64MOVBUload, ssa.OpMIPS64MOVHload, ssa.OpMIPS64MOVHUload,
				ssa.OpMIPS64MOVWload, ssa.OpMIPS64MOVWUload, ssa.OpMIPS64MOVVload,
				ssa.OpMIPS64MOVFload, ssa.OpMIPS64MOVDload,
				ssa.OpMIPS64MOVBstore, ssa.OpMIPS64MOVHstore, ssa.OpMIPS64MOVWstore, ssa.OpMIPS64MOVVstore,
				ssa.OpMIPS64MOVFstore, ssa.OpMIPS64MOVDstore,
				ssa.OpMIPS64MOVBstorezero, ssa.OpMIPS64MOVHstorezero, ssa.OpMIPS64MOVWstorezero, ssa.OpMIPS64MOVVstorezero:
				// arg0 is ptr, auxint is offset
				if w.Args[0] == v.Args[0] && w.Aux == nil && w.AuxInt >= 0 && w.AuxInt < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpMIPS64DUFFZERO, ssa.OpMIPS64LoweredZero:
				// arg0 is ptr
				if w.Args[0] == v.Args[0] {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpMIPS64LoweredMove:
				// arg0 is dst ptr, arg1 is src ptr
				if w.Args[0] == v.Args[0] || w.Args[1] == v.Args[0] {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			default:
			}
			if w.Type.IsMemory() {
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
		p := gc.Prog(mips.AMOVB)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = mips.REGTMP
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
	case ssa.OpMIPS64FPFlagTrue,
		ssa.OpMIPS64FPFlagFalse:
		// MOVV	$0, r
		// BFPF	2(PC)
		// MOVV	$1, r
		branch := mips.ABFPF
		if v.Op == ssa.OpMIPS64FPFlagFalse {
			branch = mips.ABFPT
		}
		p := gc.Prog(mips.AMOVV)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = mips.REGZERO
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
		p2 := gc.Prog(branch)
		p2.To.Type = obj.TYPE_BRANCH
		p3 := gc.Prog(mips.AMOVV)
		p3.From.Type = obj.TYPE_CONST
		p3.From.Offset = 1
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = gc.SSARegNum(v)
		p4 := gc.Prog(obj.ANOP) // not a machine instruction, for branch to land
		gc.Patch(p2, p4)
	case ssa.OpSelect0, ssa.OpSelect1:
		// nothing to do
	case ssa.OpMIPS64LoweredGetClosurePtr:
		// Closure pointer is R22 (mips.REGCTXT).
		gc.CheckLoweredGetClosurePtr(v)
	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
}

var blockJump = map[ssa.BlockKind]struct {
	asm, invasm obj.As
}{
	ssa.BlockMIPS64EQ:  {mips.ABEQ, mips.ABNE},
	ssa.BlockMIPS64NE:  {mips.ABNE, mips.ABEQ},
	ssa.BlockMIPS64LTZ: {mips.ABLTZ, mips.ABGEZ},
	ssa.BlockMIPS64GEZ: {mips.ABGEZ, mips.ABLTZ},
	ssa.BlockMIPS64LEZ: {mips.ABLEZ, mips.ABGTZ},
	ssa.BlockMIPS64GTZ: {mips.ABGTZ, mips.ABLEZ},
	ssa.BlockMIPS64FPT: {mips.ABFPT, mips.ABFPF},
	ssa.BlockMIPS64FPF: {mips.ABFPF, mips.ABFPT},
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
		// defer returns in R1:
		// 0 if we should continue executing
		// 1 if we should jump to deferreturn call
		p := gc.Prog(mips.ABNE)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = mips.REGZERO
		p.Reg = mips.REG_R1
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
	case ssa.BlockMIPS64EQ, ssa.BlockMIPS64NE,
		ssa.BlockMIPS64LTZ, ssa.BlockMIPS64GEZ,
		ssa.BlockMIPS64LEZ, ssa.BlockMIPS64GTZ,
		ssa.BlockMIPS64FPT, ssa.BlockMIPS64FPF:
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
		if !b.Control.Type.IsFlags() {
			p.From.Type = obj.TYPE_REG
			p.From.Reg = gc.SSARegNum(b.Control)
		}
	default:
		b.Unimplementedf("branch not implemented: %s. Control: %s", b.LongString(), b.Control.LongString())
	}
}

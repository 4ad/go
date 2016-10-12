// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s390x

import (
	"math"

	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/s390x"
)

// Smallest possible faulting page at address zero.
const minZeroPage = 4096

// ssaRegToReg maps ssa register numbers to obj register numbers.
var ssaRegToReg = []int16{
	s390x.REG_R0,
	s390x.REG_R1,
	s390x.REG_R2,
	s390x.REG_R3,
	s390x.REG_R4,
	s390x.REG_R5,
	s390x.REG_R6,
	s390x.REG_R7,
	s390x.REG_R8,
	s390x.REG_R9,
	s390x.REG_R10,
	s390x.REG_R11,
	s390x.REG_R12,
	s390x.REG_R13,
	s390x.REG_R14,
	s390x.REG_R15,
	s390x.REG_F0,
	s390x.REG_F1,
	s390x.REG_F2,
	s390x.REG_F3,
	s390x.REG_F4,
	s390x.REG_F5,
	s390x.REG_F6,
	s390x.REG_F7,
	s390x.REG_F8,
	s390x.REG_F9,
	s390x.REG_F10,
	s390x.REG_F11,
	s390x.REG_F12,
	s390x.REG_F13,
	s390x.REG_F14,
	s390x.REG_F15,
	0, // SB isn't a real register.  We fill an Addr.Reg field with 0 in this case.
}

// markMoves marks any MOVXconst ops that need to avoid clobbering flags.
func ssaMarkMoves(s *gc.SSAGenState, b *ssa.Block) {
	flive := b.FlagsLiveAtEnd
	if b.Control != nil && b.Control.Type.IsFlags() {
		flive = true
	}
	for i := len(b.Values) - 1; i >= 0; i-- {
		v := b.Values[i]
		if flive && v.Op == ssa.OpS390XMOVDconst {
			// The "mark" is any non-nil Aux value.
			v.Aux = v
		}
		if v.Type.IsFlags() {
			flive = false
		}
		for _, a := range v.Args {
			if a.Type.IsFlags() {
				flive = true
			}
		}
	}
}

// loadByType returns the load instruction of the given type.
func loadByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return s390x.AFMOVS
		case 8:
			return s390x.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return s390x.AMOVB
			} else {
				return s390x.AMOVBZ
			}
		case 2:
			if t.IsSigned() {
				return s390x.AMOVH
			} else {
				return s390x.AMOVHZ
			}
		case 4:
			if t.IsSigned() {
				return s390x.AMOVW
			} else {
				return s390x.AMOVWZ
			}
		case 8:
			return s390x.AMOVD
		}
	}
	panic("bad load type")
}

// storeByType returns the store instruction of the given type.
func storeByType(t ssa.Type) obj.As {
	width := t.Size()
	if t.IsFloat() {
		switch width {
		case 4:
			return s390x.AFMOVS
		case 8:
			return s390x.AFMOVD
		}
	} else {
		switch width {
		case 1:
			return s390x.AMOVB
		case 2:
			return s390x.AMOVH
		case 4:
			return s390x.AMOVW
		case 8:
			return s390x.AMOVD
		}
	}
	panic("bad store type")
}

// moveByType returns the reg->reg move instruction of the given type.
func moveByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		return s390x.AFMOVD
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return s390x.AMOVB
			} else {
				return s390x.AMOVBZ
			}
		case 2:
			if t.IsSigned() {
				return s390x.AMOVH
			} else {
				return s390x.AMOVHZ
			}
		case 4:
			if t.IsSigned() {
				return s390x.AMOVW
			} else {
				return s390x.AMOVWZ
			}
		case 8:
			return s390x.AMOVD
		}
	}
	panic("bad load type")
}

// opregreg emits instructions for
//     dest := dest(To) op src(From)
// and also returns the created obj.Prog so it
// may be further adjusted (offset, scale, etc).
func opregreg(op obj.As, dest, src int16) *obj.Prog {
	p := gc.Prog(op)
	p.From.Type = obj.TYPE_REG
	p.To.Type = obj.TYPE_REG
	p.To.Reg = dest
	p.From.Reg = src
	return p
}

// opregregimm emits instructions for
//	dest := src(From) op off
// and also returns the created obj.Prog so it
// may be further adjusted (offset, scale, etc).
func opregregimm(op obj.As, dest, src int16, off int64) *obj.Prog {
	p := gc.Prog(op)
	p.From.Type = obj.TYPE_CONST
	p.From.Offset = off
	p.Reg = src
	p.To.Reg = dest
	p.To.Type = obj.TYPE_REG
	return p
}

func ssaGenValue(s *gc.SSAGenState, v *ssa.Value) {
	s.SetLineno(v.Line)
	switch v.Op {
	case ssa.OpS390XSLD, ssa.OpS390XSLW,
		ssa.OpS390XSRD, ssa.OpS390XSRW,
		ssa.OpS390XSRAD, ssa.OpS390XSRAW:
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])
		if r2 == s390x.REG_R0 {
			v.Fatalf("cannot use R0 as shift value %s", v.LongString())
		}
		p := opregreg(v.Op.Asm(), r, r2)
		if r != r1 {
			p.Reg = r1
		}
	case ssa.OpS390XADD, ssa.OpS390XADDW,
		ssa.OpS390XSUB, ssa.OpS390XSUBW,
		ssa.OpS390XAND, ssa.OpS390XANDW,
		ssa.OpS390XOR, ssa.OpS390XORW,
		ssa.OpS390XXOR, ssa.OpS390XXORW:
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])
		p := opregreg(v.Op.Asm(), r, r2)
		if r != r1 {
			p.Reg = r1
		}
	// 2-address opcode arithmetic
	case ssa.OpS390XMULLD, ssa.OpS390XMULLW,
		ssa.OpS390XMULHD, ssa.OpS390XMULHDU,
		ssa.OpS390XFADDS, ssa.OpS390XFADD, ssa.OpS390XFSUBS, ssa.OpS390XFSUB,
		ssa.OpS390XFMULS, ssa.OpS390XFMUL, ssa.OpS390XFDIVS, ssa.OpS390XFDIV:
		r := gc.SSARegNum(v)
		if r != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}
		opregreg(v.Op.Asm(), r, gc.SSARegNum(v.Args[1]))
	case ssa.OpS390XDIVD, ssa.OpS390XDIVW,
		ssa.OpS390XDIVDU, ssa.OpS390XDIVWU,
		ssa.OpS390XMODD, ssa.OpS390XMODW,
		ssa.OpS390XMODDU, ssa.OpS390XMODWU:

		// TODO(mundaym): use the temp registers every time like x86 does with AX?
		dividend := gc.SSARegNum(v.Args[0])
		divisor := gc.SSARegNum(v.Args[1])

		// CPU faults upon signed overflow, which occurs when most
		// negative int is divided by -1.
		var j *obj.Prog
		if v.Op == ssa.OpS390XDIVD || v.Op == ssa.OpS390XDIVW ||
			v.Op == ssa.OpS390XMODD || v.Op == ssa.OpS390XMODW {

			var c *obj.Prog
			c = gc.Prog(s390x.ACMP)
			j = gc.Prog(s390x.ABEQ)

			c.From.Type = obj.TYPE_REG
			c.From.Reg = divisor
			c.To.Type = obj.TYPE_CONST
			c.To.Offset = -1

			j.To.Type = obj.TYPE_BRANCH

		}

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = divisor
		p.Reg = 0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = dividend

		// signed division, rest of the check for -1 case
		if j != nil {
			j2 := gc.Prog(s390x.ABR)
			j2.To.Type = obj.TYPE_BRANCH

			var n *obj.Prog
			if v.Op == ssa.OpS390XDIVD || v.Op == ssa.OpS390XDIVW {
				// n * -1 = -n
				n = gc.Prog(s390x.ANEG)
				n.To.Type = obj.TYPE_REG
				n.To.Reg = dividend
			} else {
				// n % -1 == 0
				n = gc.Prog(s390x.AXOR)
				n.From.Type = obj.TYPE_REG
				n.From.Reg = dividend
				n.To.Type = obj.TYPE_REG
				n.To.Reg = dividend
			}

			j.To.Val = n
			j2.To.Val = s.Pc()
		}
	case ssa.OpS390XADDconst, ssa.OpS390XADDWconst:
		opregregimm(v.Op.Asm(), gc.SSARegNum(v), gc.SSARegNum(v.Args[0]), v.AuxInt)
	case ssa.OpS390XMULLDconst, ssa.OpS390XMULLWconst,
		ssa.OpS390XSUBconst, ssa.OpS390XSUBWconst,
		ssa.OpS390XANDconst, ssa.OpS390XANDWconst,
		ssa.OpS390XORconst, ssa.OpS390XORWconst,
		ssa.OpS390XXORconst, ssa.OpS390XXORWconst:
		r := gc.SSARegNum(v)
		if r != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XSLDconst, ssa.OpS390XSLWconst,
		ssa.OpS390XSRDconst, ssa.OpS390XSRWconst,
		ssa.OpS390XSRADconst, ssa.OpS390XSRAWconst,
		ssa.OpS390XRLLGconst, ssa.OpS390XRLLconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		if r != r1 {
			p.Reg = r1
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XSUBEcarrymask, ssa.OpS390XSUBEWcarrymask:
		r := gc.SSARegNum(v)
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XMOVDaddridx:
		r := gc.SSARegNum(v.Args[0])
		i := gc.SSARegNum(v.Args[1])
		p := gc.Prog(s390x.AMOVD)
		p.From.Scale = 1
		if i == s390x.REGSP {
			r, i = i, r
		}
		p.From.Type = obj.TYPE_ADDR
		p.From.Reg = r
		p.From.Index = i
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpS390XMOVDaddr:
		p := gc.Prog(s390x.AMOVD)
		p.From.Type = obj.TYPE_ADDR
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpS390XCMP, ssa.OpS390XCMPW, ssa.OpS390XCMPU, ssa.OpS390XCMPWU:
		opregreg(v.Op.Asm(), gc.SSARegNum(v.Args[1]), gc.SSARegNum(v.Args[0]))
	case ssa.OpS390XTESTB:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 0xFF
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s390x.REGTMP
	case ssa.OpS390XFCMPS, ssa.OpS390XFCMP:
		opregreg(v.Op.Asm(), gc.SSARegNum(v.Args[1]), gc.SSARegNum(v.Args[0]))
	case ssa.OpS390XCMPconst, ssa.OpS390XCMPWconst, ssa.OpS390XCMPUconst, ssa.OpS390XCMPWUconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_CONST
		p.To.Offset = v.AuxInt
	case ssa.OpS390XMOVDconst:
		x := gc.SSARegNum(v)
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x
	case ssa.OpS390XFMOVSconst, ssa.OpS390XFMOVDconst:
		x := gc.SSARegNum(v)
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_FCONST
		p.From.Val = math.Float64frombits(uint64(v.AuxInt))
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x
	case ssa.OpS390XMOVDload,
		ssa.OpS390XMOVWZload, ssa.OpS390XMOVHZload, ssa.OpS390XMOVBZload,
		ssa.OpS390XMOVDBRload, ssa.OpS390XMOVWBRload, ssa.OpS390XMOVHBRload,
		ssa.OpS390XMOVBload, ssa.OpS390XMOVHload, ssa.OpS390XMOVWload,
		ssa.OpS390XFMOVSload, ssa.OpS390XFMOVDload:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpS390XMOVBZloadidx, ssa.OpS390XMOVHZloadidx, ssa.OpS390XMOVWZloadidx, ssa.OpS390XMOVDloadidx,
		ssa.OpS390XMOVHBRloadidx, ssa.OpS390XMOVWBRloadidx, ssa.OpS390XMOVDBRloadidx,
		ssa.OpS390XFMOVSloadidx, ssa.OpS390XFMOVDloadidx:
		r := gc.SSARegNum(v.Args[0])
		i := gc.SSARegNum(v.Args[1])
		if i == s390x.REGSP {
			r, i = i, r
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = r
		p.From.Scale = 1
		p.From.Index = i
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpS390XMOVBstore, ssa.OpS390XMOVHstore, ssa.OpS390XMOVWstore, ssa.OpS390XMOVDstore,
		ssa.OpS390XFMOVSstore, ssa.OpS390XFMOVDstore:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpS390XMOVBstoreidx, ssa.OpS390XMOVHstoreidx, ssa.OpS390XMOVWstoreidx, ssa.OpS390XMOVDstoreidx,
		ssa.OpS390XFMOVSstoreidx, ssa.OpS390XFMOVDstoreidx:
		r := gc.SSARegNum(v.Args[0])
		i := gc.SSARegNum(v.Args[1])
		if i == s390x.REGSP {
			r, i = i, r
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[2])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = r
		p.To.Scale = 1
		p.To.Index = i
		gc.AddAux(&p.To, v)
	case ssa.OpS390XMOVDstoreconst, ssa.OpS390XMOVWstoreconst, ssa.OpS390XMOVHstoreconst, ssa.OpS390XMOVBstoreconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		sc := v.AuxValAndOff()
		p.From.Offset = sc.Val()
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux2(&p.To, v, sc.Off())
	case ssa.OpS390XMOVBreg, ssa.OpS390XMOVHreg, ssa.OpS390XMOVWreg,
		ssa.OpS390XMOVBZreg, ssa.OpS390XMOVHZreg, ssa.OpS390XMOVWZreg,
		ssa.OpS390XCEFBRA, ssa.OpS390XCDFBRA, ssa.OpS390XCEGBRA, ssa.OpS390XCDGBRA,
		ssa.OpS390XCFEBRA, ssa.OpS390XCFDBRA, ssa.OpS390XCGEBRA, ssa.OpS390XCGDBRA,
		ssa.OpS390XLDEBR, ssa.OpS390XLEDBR,
		ssa.OpS390XFNEG, ssa.OpS390XFNEGS:
		opregreg(v.Op.Asm(), gc.SSARegNum(v), gc.SSARegNum(v.Args[0]))
	case ssa.OpS390XCLEAR:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		sc := v.AuxValAndOff()
		p.From.Offset = sc.Val()
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux2(&p.To, v, sc.Off())
	case ssa.OpCopy, ssa.OpS390XMOVDconvert:
		if v.Type.IsMemory() {
			return
		}
		x := gc.SSARegNum(v.Args[0])
		y := gc.SSARegNum(v)
		if x != y {
			opregreg(moveByType(v.Type), y, x)
		}
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
	case ssa.OpPhi:
		gc.CheckLoweredPhi(v)
	case ssa.OpInitMem:
		// memory arg needs no code
	case ssa.OpArg:
		// input args need no code
	case ssa.OpS390XLoweredGetClosurePtr:
		// Closure pointer is R12 (already)
		gc.CheckLoweredGetClosurePtr(v)
	case ssa.OpS390XLoweredGetG:
		r := gc.SSARegNum(v)
		p := gc.Prog(s390x.AMOVD)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = s390x.REGG
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XCALLstatic:
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
	case ssa.OpS390XCALLclosure:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpS390XCALLdefer:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Deferproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpS390XCALLgo:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Newproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpS390XCALLinter:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpS390XNEG, ssa.OpS390XNEGW:
		r := gc.SSARegNum(v)
		p := gc.Prog(v.Op.Asm())
		r1 := gc.SSARegNum(v.Args[0])
		if r != r1 {
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r1
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XNOT, ssa.OpS390XNOTW:
		v.Fatalf("NOT/NOTW generated %s", v.LongString())
	case ssa.OpS390XMOVDEQ, ssa.OpS390XMOVDNE,
		ssa.OpS390XMOVDLT, ssa.OpS390XMOVDLE,
		ssa.OpS390XMOVDGT, ssa.OpS390XMOVDGE,
		ssa.OpS390XMOVDGTnoinv, ssa.OpS390XMOVDGEnoinv:
		r := gc.SSARegNum(v)
		if r != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpS390XFSQRT:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)
	case ssa.OpSP, ssa.OpSB:
		// nothing to do
	case ssa.OpVarDef:
		gc.Gvardef(v.Aux.(*gc.Node))
	case ssa.OpVarKill:
		gc.Gvarkill(v.Aux.(*gc.Node))
	case ssa.OpVarLive:
		gc.Gvarlive(v.Aux.(*gc.Node))
	case ssa.OpKeepAlive:
		gc.KeepAlive(v)
	case ssa.OpS390XInvertFlags:
		v.Fatalf("InvertFlags should never make it to codegen %v", v.LongString())
	case ssa.OpS390XFlagEQ, ssa.OpS390XFlagLT, ssa.OpS390XFlagGT:
		v.Fatalf("Flag* ops should never make it to codegen %v", v.LongString())
	case ssa.OpS390XLoweredNilCheck:
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
			case ssa.OpS390XMOVDload,
				ssa.OpS390XMOVBload, ssa.OpS390XMOVHload, ssa.OpS390XMOVWload,
				ssa.OpS390XMOVBZload, ssa.OpS390XMOVHZload, ssa.OpS390XMOVWZload,
				ssa.OpS390XMOVHBRload, ssa.OpS390XMOVWBRload, ssa.OpS390XMOVDBRload,
				ssa.OpS390XMOVBstore, ssa.OpS390XMOVHstore, ssa.OpS390XMOVWstore, ssa.OpS390XMOVDstore,
				ssa.OpS390XFMOVSload, ssa.OpS390XFMOVDload,
				ssa.OpS390XFMOVSstore, ssa.OpS390XFMOVDstore,
				ssa.OpS390XSTMG2, ssa.OpS390XSTMG3, ssa.OpS390XSTMG4,
				ssa.OpS390XSTM2, ssa.OpS390XSTM3, ssa.OpS390XSTM4:
				if w.Args[0] == v.Args[0] && w.Aux == nil && w.AuxInt >= 0 && w.AuxInt < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpS390XMOVDstoreconst, ssa.OpS390XMOVWstoreconst, ssa.OpS390XMOVHstoreconst, ssa.OpS390XMOVBstoreconst,
				ssa.OpS390XCLEAR:
				off := ssa.ValAndOff(v.AuxInt).Off()
				if w.Args[0] == v.Args[0] && w.Aux == nil && off >= 0 && off < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpS390XMVC:
				off := ssa.ValAndOff(v.AuxInt).Off()
				if (w.Args[0] == v.Args[0] || w.Args[1] == v.Args[0]) && w.Aux == nil && off >= 0 && off < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
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
		// Issue a load which will fault if the input is nil.
		p := gc.Prog(s390x.AMOVBZ)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s390x.REGTMP
		if gc.Debug_checknil != 0 && v.Line > 1 { // v.Line==1 in generated wrappers
			gc.Warnl(v.Line, "generated nil check")
		}
	case ssa.OpS390XMVC:
		vo := v.AuxValAndOff()
		p := gc.Prog(s390x.AMVC)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.From.Offset = vo.Off()
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		p.To.Offset = vo.Off()
		p.From3 = new(obj.Addr)
		p.From3.Type = obj.TYPE_CONST
		p.From3.Offset = vo.Val()
	case ssa.OpS390XSTMG2, ssa.OpS390XSTMG3, ssa.OpS390XSTMG4,
		ssa.OpS390XSTM2, ssa.OpS390XSTM3, ssa.OpS390XSTM4:
		for i := 2; i < len(v.Args)-1; i++ {
			if gc.SSARegNum(v.Args[i]) != gc.SSARegNum(v.Args[i-1])+1 {
				v.Fatalf("invalid store multiple %s", v.LongString())
			}
		}
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.Reg = gc.SSARegNum(v.Args[len(v.Args)-2])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
	case ssa.OpS390XLoweredMove:
		// Inputs must be valid pointers to memory,
		// so adjust arg0 and arg1 as part of the expansion.
		// arg2 should be src+size,
		//
		// mvc: MVC  $256, 0(R2), 0(R1)
		//      MOVD $256(R1), R1
		//      MOVD $256(R2), R2
		//      CMP  R2, Rarg2
		//      BNE  mvc
		//      MVC  $rem, 0(R2), 0(R1) // if rem > 0
		// arg2 is the last address to move in the loop + 256
		mvc := gc.Prog(s390x.AMVC)
		mvc.From.Type = obj.TYPE_MEM
		mvc.From.Reg = gc.SSARegNum(v.Args[1])
		mvc.To.Type = obj.TYPE_MEM
		mvc.To.Reg = gc.SSARegNum(v.Args[0])
		mvc.From3 = new(obj.Addr)
		mvc.From3.Type = obj.TYPE_CONST
		mvc.From3.Offset = 256

		for i := 0; i < 2; i++ {
			movd := gc.Prog(s390x.AMOVD)
			movd.From.Type = obj.TYPE_ADDR
			movd.From.Reg = gc.SSARegNum(v.Args[i])
			movd.From.Offset = 256
			movd.To.Type = obj.TYPE_REG
			movd.To.Reg = gc.SSARegNum(v.Args[i])
		}

		cmpu := gc.Prog(s390x.ACMPU)
		cmpu.From.Reg = gc.SSARegNum(v.Args[1])
		cmpu.From.Type = obj.TYPE_REG
		cmpu.To.Reg = gc.SSARegNum(v.Args[2])
		cmpu.To.Type = obj.TYPE_REG

		bne := gc.Prog(s390x.ABLT)
		bne.To.Type = obj.TYPE_BRANCH
		gc.Patch(bne, mvc)

		if v.AuxInt > 0 {
			mvc := gc.Prog(s390x.AMVC)
			mvc.From.Type = obj.TYPE_MEM
			mvc.From.Reg = gc.SSARegNum(v.Args[1])
			mvc.To.Type = obj.TYPE_MEM
			mvc.To.Reg = gc.SSARegNum(v.Args[0])
			mvc.From3 = new(obj.Addr)
			mvc.From3.Type = obj.TYPE_CONST
			mvc.From3.Offset = v.AuxInt
		}
	case ssa.OpS390XLoweredZero:
		// Input must be valid pointers to memory,
		// so adjust arg0 as part of the expansion.
		// arg1 should be src+size,
		//
		// clear: CLEAR $256, 0(R1)
		//        MOVD  $256(R1), R1
		//        CMP   R1, Rarg1
		//        BNE   clear
		//        CLEAR $rem, 0(R1) // if rem > 0
		// arg1 is the last address to zero in the loop + 256
		clear := gc.Prog(s390x.ACLEAR)
		clear.From.Type = obj.TYPE_CONST
		clear.From.Offset = 256
		clear.To.Type = obj.TYPE_MEM
		clear.To.Reg = gc.SSARegNum(v.Args[0])

		movd := gc.Prog(s390x.AMOVD)
		movd.From.Type = obj.TYPE_ADDR
		movd.From.Reg = gc.SSARegNum(v.Args[0])
		movd.From.Offset = 256
		movd.To.Type = obj.TYPE_REG
		movd.To.Reg = gc.SSARegNum(v.Args[0])

		cmpu := gc.Prog(s390x.ACMPU)
		cmpu.From.Reg = gc.SSARegNum(v.Args[0])
		cmpu.From.Type = obj.TYPE_REG
		cmpu.To.Reg = gc.SSARegNum(v.Args[1])
		cmpu.To.Type = obj.TYPE_REG

		bne := gc.Prog(s390x.ABLT)
		bne.To.Type = obj.TYPE_BRANCH
		gc.Patch(bne, clear)

		if v.AuxInt > 0 {
			clear := gc.Prog(s390x.ACLEAR)
			clear.From.Type = obj.TYPE_CONST
			clear.From.Offset = v.AuxInt
			clear.To.Type = obj.TYPE_MEM
			clear.To.Reg = gc.SSARegNum(v.Args[0])
		}
	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
}

var blockJump = [...]struct {
	asm, invasm obj.As
}{
	ssa.BlockS390XEQ:  {s390x.ABEQ, s390x.ABNE},
	ssa.BlockS390XNE:  {s390x.ABNE, s390x.ABEQ},
	ssa.BlockS390XLT:  {s390x.ABLT, s390x.ABGE},
	ssa.BlockS390XGE:  {s390x.ABGE, s390x.ABLT},
	ssa.BlockS390XLE:  {s390x.ABLE, s390x.ABGT},
	ssa.BlockS390XGT:  {s390x.ABGT, s390x.ABLE},
	ssa.BlockS390XGTF: {s390x.ABGT, s390x.ABLEU},
	ssa.BlockS390XGEF: {s390x.ABGE, s390x.ABLTU},
}

func ssaGenBlock(s *gc.SSAGenState, b, next *ssa.Block) {
	s.SetLineno(b.Line)

	switch b.Kind {
	case ssa.BlockPlain, ssa.BlockCheck:
		if b.Succs[0].Block() != next {
			p := gc.Prog(s390x.ABR)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		}
	case ssa.BlockDefer:
		// defer returns in R3:
		// 0 if we should continue executing
		// 1 if we should jump to deferreturn call
		p := gc.Prog(s390x.AAND)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 0xFFFFFFFF
		p.Reg = s390x.REG_R3
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s390x.REG_R3
		p = gc.Prog(s390x.ABNE)
		p.To.Type = obj.TYPE_BRANCH
		s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[1].Block()})
		if b.Succs[0].Block() != next {
			p := gc.Prog(s390x.ABR)
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		}
	case ssa.BlockExit:
		gc.Prog(obj.AUNDEF) // tell plive.go that we never reach here
	case ssa.BlockRet:
		gc.Prog(obj.ARET)
	case ssa.BlockRetJmp:
		p := gc.Prog(s390x.ABR)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(b.Aux.(*gc.Sym))
	case ssa.BlockS390XEQ, ssa.BlockS390XNE,
		ssa.BlockS390XLT, ssa.BlockS390XGE,
		ssa.BlockS390XLE, ssa.BlockS390XGT,
		ssa.BlockS390XGEF, ssa.BlockS390XGTF:
		jmp := blockJump[b.Kind]
		likely := b.Likely
		var p *obj.Prog
		switch next {
		case b.Succs[0].Block():
			p = gc.Prog(jmp.invasm)
			likely *= -1
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
			q := gc.Prog(s390x.ABR)
			q.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: q, B: b.Succs[1].Block()})
		}
	default:
		b.Unimplementedf("branch not implemented: %s. Control: %s", b.LongString(), b.Control.LongString())
	}
}

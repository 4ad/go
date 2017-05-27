// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
	"math"
)

var ssaRegToReg = []int16{
	// sparc64.REG_ZR,   // zero register, not used by the compiler
	sparc64.REG_RT1,  // for runtime, liblink and duff device
	sparc64.REG_CTXT, // environment for closures
	sparc64.REG_G,    // g register
	sparc64.REG_RT2,  // for runtime, liblink and duff device
	// sparc64.REG_TMP,  // reserved for runtime and liblink
	// sparc64.REG_G6,   // reserved for the operating system
	// sparc64.REG_TLS,  // reserved for the operating system
	sparc64.REG_O0,
	sparc64.REG_O1,
	sparc64.REG_O2,
	sparc64.REG_O3,
	sparc64.REG_O4,
	sparc64.REG_O5,
	sparc64.REG_RSP,  // machine stack pointer
	// sparc64.REG_OLR,  // the output link register
	// sparc64.REG_TMP2, // reserved for runtime and liblink
	sparc64.REG_L1,
	sparc64.REG_L2,
	sparc64.REG_L3,
	sparc64.REG_L4,
	sparc64.REG_L5,
	sparc64.REG_L6,
	sparc64.REG_L7,
	sparc64.REG_I0,
	sparc64.REG_I1,
	sparc64.REG_I2,
	sparc64.REG_I3,
	sparc64.REG_I4,
	sparc64.REG_I5,
	sparc64.REG_RFP, // frame pointer
	// sparc64.REG_ILR, // the input link register

	sparc64.REG_Y0,
	sparc64.REG_Y1,
	sparc64.REG_Y2,
	sparc64.REG_Y3,
	sparc64.REG_Y4,
	sparc64.REG_Y5,
	sparc64.REG_Y6,
	sparc64.REG_Y7,
	sparc64.REG_Y8,
	sparc64.REG_Y9,
	sparc64.REG_Y10,
	sparc64.REG_Y11,
	sparc64.REG_Y12,
	sparc64.REG_Y13,
	// sparc64.REG_YTWO, // uncertain if used
	// sparc64.REG_YTMP, // uncertain if used

	0, // SB, pseudo symbol static base
	1, // SP, pseudo stack pointer
	2, // FP, pseudo frame pointer
}

// Smallest possible faulting page at address zero,
// see ../../../../runtime/mheap.go:/minPhysPageSize
const minZeroPage = 4096

// loadByType returns the load instruction of the given type.
func loadByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return sparc64.AFMOVS
		case 8:
			return sparc64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return sparc64.AMOVB
			} else {
				return sparc64.AMOVUB
			}
		case 2:
			if t.IsSigned() {
				return sparc64.AMOVH
			} else {
				return sparc64.AMOVUH
			}
		case 4:
			if t.IsSigned() {
				return sparc64.AMOVW
			} else {
				return sparc64.AMOVUW
			}
		case 8:
			return sparc64.AMOVD
		}
	}
	panic("bad load type")
}

// storeByType returns the store instruction of the given type.
func storeByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return sparc64.AFMOVS
		case 8:
			return sparc64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			return sparc64.AMOVB
		case 2:
			return sparc64.AMOVH
		case 4:
			return sparc64.AMOVW
		case 8:
			return sparc64.AMOVD
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

	case ssa.OpCopy, ssa.OpSPARC64MOVDconvert:
		if v.Type.IsMemory() {
			return
		}
		x := gc.SSARegNum(v.Args[0])
		y := gc.SSARegNum(v)
		if x == y {
			return
		}
		as := sparc64.AMOVD
		if v.Type.IsFloat() {
			switch v.Type.Size() {
			case 4:
				as = sparc64.AFMOVS
			case 8:
				as = sparc64.AFMOVD
			default:
				panic("bad float size")
			}
		}
		p := gc.Prog(as)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x
		p.To.Type = obj.TYPE_REG
		p.To.Reg = y

	case ssa.OpLoadReg:
		loadOp := loadByType(v.Type)
		n, off := gc.AutoVar(v.Args[0])
		p := gc.Prog(loadOp)
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
		storeOp := storeByType(v.Type)
		n, off := gc.AutoVar(v)
		p := gc.Prog(storeOp)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
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

	case ssa.OpVarDef:
		gc.Gvardef(v.Aux.(*gc.Node))
	case ssa.OpVarKill:
		gc.Gvarkill(v.Aux.(*gc.Node))
	case ssa.OpVarLive:
		gc.Gvarlive(v.Aux.(*gc.Node))
	case ssa.OpKeepAlive:
		gc.KeepAlive(v)
	case ssa.OpPhi:
		gc.CheckLoweredPhi(v)

	case ssa.OpSPARC64SLLmax,
		ssa.OpSPARC64SRLmax:

		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r2
		p.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
		p2 := gc.Prog(sparc64.ACMP)
		p2.From.Type = obj.TYPE_CONST
		p2.From.Offset = v.AuxInt
		p2.Reg = r2
		p3 := gc.Prog(sparc64.AMOVGU)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = sparc64.REG_XCC
		p3.From3 = &obj.Addr{}
		p3.From3.Type = obj.TYPE_CONST
		p3.From3.Offset = 0
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = r

	case ssa.OpSPARC64SRAmax:

		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])

		p := gc.Prog(sparc64.AMOVD)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = sparc64.REG_TMP
		p2 := gc.Prog(sparc64.ACMP)
		p2.From.Type = obj.TYPE_CONST
		p2.From.Offset = v.AuxInt
		p2.Reg = sparc64.REG_TMP
		p3 := gc.Prog(sparc64.AMOVGU)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = sparc64.REG_XCC
		p3.From3 = &obj.Addr{}
		p3.From3.Type = obj.TYPE_CONST
		p3.From3.Offset = 63
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = sparc64.REG_TMP
		p4 := gc.Prog(v.Op.Asm())
		p4.From.Type = obj.TYPE_REG
		p4.From.Reg = r2
		p4.Reg = sparc64.REG_TMP
		p4.To.Type = obj.TYPE_REG
		p4.To.Reg = r

	case ssa.OpSPARC64ADD,
		ssa.OpSPARC64SUB,
		ssa.OpSPARC64AND,
		ssa.OpSPARC64OR,
		ssa.OpSPARC64XOR,
		ssa.OpSPARC64MULD,
		ssa.OpSPARC64SLL,
		ssa.OpSPARC64SRL,
		ssa.OpSPARC64SRA,
		ssa.OpSPARC64SDIVD,
		ssa.OpSPARC64UDIVD,
		ssa.OpSPARC64FADDS,
		ssa.OpSPARC64FADDD,
		ssa.OpSPARC64FSUBS,
		ssa.OpSPARC64FSUBD,
		ssa.OpSPARC64FMULS,
		ssa.OpSPARC64FMULD,
		ssa.OpSPARC64FDIVS,
		ssa.OpSPARC64FDIVD:

		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		r2 := gc.SSARegNum(v.Args[1])
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r2
		p.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r

	case ssa.OpSPARC64ADDconst,
		ssa.OpSPARC64SUBconst,
		ssa.OpSPARC64ANDconst,
		ssa.OpSPARC64ORconst,
		ssa.OpSPARC64XORconst,
		ssa.OpSPARC64SLLconst,
		ssa.OpSPARC64SRLconst,
		ssa.OpSPARC64SRAconst:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64MOVBreg,
		ssa.OpSPARC64MOVUBreg,
		ssa.OpSPARC64MOVHreg,
		ssa.OpSPARC64MOVUHreg,
		ssa.OpSPARC64MOVWreg,
		ssa.OpSPARC64MOVUWreg,
		ssa.OpSPARC64MOVDreg,
		ssa.OpSPARC64NEG,
		ssa.OpSPARC64FNEGS,
		ssa.OpSPARC64FNEGD,
		ssa.OpSPARC64FSQRTD,
		ssa.OpSPARC64FSTOD,
		ssa.OpSPARC64FDTOS:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64FSTOI,
		ssa.OpSPARC64FSTOX,
		ssa.OpSPARC64FDTOI:

		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		ft := v.Args[0].Type
		tt := v.Type

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = sparc64.REG_YTMP
		p = gc.Prog(storeByType(ft))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = sparc64.REG_YTMP
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = sparc64.REG_RSP
		p.To.Offset = -8 + sparc64.StackBias
		p = gc.Prog(loadByType(tt))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = sparc64.REG_RSP
		p.From.Offset = -8 + sparc64.StackBias
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r

	case ssa.OpSPARC64FDTOX:
	// algorithm is:
	//	if small enough, use native float64 -> int64 conversion.
	//	otherwise, subtract 2^63, convert, and add it back.
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		ft := v.Args[0].Type
		tt := v.Type

		if !tt.IsSigned() {
			p := gc.Prog(storeByType(ft))
			p.From.Type = obj.TYPE_FCONST
			p.From.Val = float64(uint64(1<<63))
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_YTMP
			p = gc.Prog(sparc64.AFCMPD)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r1
			p.Reg = sparc64.REG_YTMP
			q := gc.Prog(sparc64.AFBG)
			q.To.Type = obj.TYPE_BRANCH
			q.To.Val = nil
			p = gc.Prog(sparc64.AFSUBD)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r1
			p.Reg = sparc64.REG_YTMP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_YTMP
			q2 := gc.Prog(obj.AJMP)
			q2.To.Type = obj.TYPE_BRANCH
			q2.To.Val = nil
			gc.Patch(q, gc.Pc)
			p = gc.Prog(storeByType(ft))
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r1
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_YTMP
			gc.Patch(q2, gc.Pc)
			r1 = sparc64.REG_YTMP
		}

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = sparc64.REG_YTMP
		p = gc.Prog(storeByType(ft))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = sparc64.REG_YTMP
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = sparc64.REG_RSP
		p.To.Offset = -8 + sparc64.StackBias
		p = gc.Prog(loadByType(tt))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = sparc64.REG_RSP
		p.From.Offset = -8 + sparc64.StackBias
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r

		if !tt.IsSigned() {
			q := gc.Prog(sparc64.AFBG)
			q.To.Type = obj.TYPE_BRANCH
			q.To.Val = nil
			p := gc.Prog(sparc64.AMOVD)
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 1
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			p = gc.Prog(sparc64.ASLLD)
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 63
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			p = gc.Prog(sparc64.AADD)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = sparc64.REG_TMP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = r
			gc.Patch(q, gc.Pc)
		}

	case ssa.OpSPARC64FITOS,
		ssa.OpSPARC64FITOD:

		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		ft := v.Args[0].Type
		tt := v.Type

		p := gc.Prog(storeByType(ft))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r1
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = sparc64.REG_RSP
		p.To.Offset = -8 + sparc64.StackBias
		p = gc.Prog(loadByType(tt))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = sparc64.REG_RSP
		p.From.Offset = -8 + sparc64.StackBias
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
		p = gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r

	case ssa.OpSPARC64FXTOS,
		ssa.OpSPARC64FXTOD:
	// algorithm is:
	//	if small enough, use native int64 -> float64 conversion,
	//	otherwise halve (x -> (x>>1)|(x&1)), convert, and double.
		r := gc.SSARegNum(v)
		r1 := gc.SSARegNum(v.Args[0])
		ft := v.Args[0].Type
		tt := v.Type

		if !ft.IsSigned() {
			p := gc.Prog(sparc64.AMOVD)
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 1
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			p = gc.Prog(sparc64.ASLLD)
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 63
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			p = gc.Prog(sparc64.ACMP)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = r1
			p.Reg = sparc64.REG_TMP
			q := gc.Prog(sparc64.ABLEUD)
			q.To.Type = obj.TYPE_BRANCH
			q.To.Val = nil
			p = gc.Prog(sparc64.AAND)
			p.Reg = r1
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 1
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP2
			p = gc.Prog(sparc64.ASRLD)
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = 1
			p.Reg = r1
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			p = gc.Prog(sparc64.AOR)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = sparc64.REG_TMP
			p.Reg = sparc64.REG_TMP2
			p.To.Type = obj.TYPE_REG
			p.To.Reg = sparc64.REG_TMP
			gc.Patch(q, gc.Pc)
			r1 = sparc64.REG_TMP
		}

		p := gc.Prog(storeByType(ft))
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r1
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = sparc64.REG_RSP
		p.To.Offset = -8 + sparc64.StackBias
		p = gc.Prog(loadByType(tt))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = sparc64.REG_RSP
		p.From.Offset = -8 + sparc64.StackBias
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
		p = gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = r
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r

		if !ft.IsSigned() {
			q := gc.Prog(sparc64.ABLEUD)
			q.To.Type = obj.TYPE_BRANCH
			q.To.Val = nil
			p := gc.Prog(sparc64.AFMULD)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = sparc64.REG_YTWO
			p.To.Type = obj.TYPE_REG
			p.To.Reg = r
			gc.Patch(q, gc.Pc)
		}

	case ssa.OpSPARC64MOVDaddr:
		p := gc.Prog(sparc64.AMOVD)
		p.From.Type = obj.TYPE_ADDR
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

		var wantreg string
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
			p.From.Reg = sparc64.REG_RSP
			p.From.Offset = v.AuxInt + sparc64.StackBias
		}
		if reg := gc.SSAReg(v.Args[0]); reg.Name() != wantreg {
			v.Fatalf("bad reg %s for symbol type %T, want %s", reg.Name(), v.Aux, wantreg)
		}

	case ssa.OpSPARC64MOVDconst,
		ssa.OpSPARC64MOVWconst:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64FMOVDconst,
		ssa.OpSPARC64FMOVSconst:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_FCONST
		p.From.Val = math.Float64frombits(uint64(v.AuxInt))
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64MOVBload,
		ssa.OpSPARC64MOVUBload,
		ssa.OpSPARC64MOVHload,
		ssa.OpSPARC64MOVUHload,
		ssa.OpSPARC64MOVWload,
		ssa.OpSPARC64MOVUWload,
		ssa.OpSPARC64MOVDload,
		ssa.OpSPARC64FMOVSload,
		ssa.OpSPARC64FMOVDload:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64MOVDstore,
		ssa.OpSPARC64MOVWstore,
		ssa.OpSPARC64MOVHstore,
		ssa.OpSPARC64MOVBstore,
		ssa.OpSPARC64FMOVSstore,
		ssa.OpSPARC64FMOVDstore:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)

	case ssa.OpSPARC64LoweredZero:
		// loop:
		// 	MOVD	ZR, (RT1)
		// 	ADD	$8, RT1
		// 	CMP	Rarg1, RT1
		// 	BLED	loop
		// arg0 is address of dst memory
		// arg1 is the address of the last element to zero
		p := gc.Prog(sparc64.AMOVD)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = sparc64.REG_ZR
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		p2 := gc.Prog(sparc64.AADD)
		p2.From.Type = obj.TYPE_CONST
		p2.From.Offset = 8
		p2.To.Type = obj.TYPE_REG
		p2.To.Reg = gc.SSARegNum(v.Args[0])
		p3 := gc.Prog(sparc64.ACMP)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = gc.SSARegNum(v.Args[1])
		p3.Reg = gc.SSARegNum(v.Args[0])
		p4 := gc.Prog(sparc64.ABLED)
		p4.To.Type = obj.TYPE_BRANCH
		gc.Patch(p4, p)

	case ssa.OpSPARC64LoweredMove:
		// loop:
		// 	MOVD	(RT1), TMP
		// 	ADD	$8, RT1
		//	MOVD	TMP, (RT2)
		//	ADD	$8, RT2
		// 	CMP	Rarg2, RT1
		// 	BLED	loop
		// arg0 is address of dst memory
		// arg1 is address of src memory
		// arg2 is the address of the last element of src
		p := gc.Prog(sparc64.AMOVD)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = sparc64.REG_TMP
		p2 := gc.Prog(sparc64.AADD)
		p2.From.Type = obj.TYPE_CONST
		p2.From.Offset = 8
		p2.To.Type = obj.TYPE_REG
		p2.To.Reg = gc.SSARegNum(v.Args[1])
		p3 := gc.Prog(sparc64.AMOVD)
		p3.From.Type = obj.TYPE_REG
		p3.From.Reg = sparc64.REG_TMP
		p3.To.Type = obj.TYPE_MEM
		p3.To.Reg = gc.SSARegNum(v.Args[0])
		p4 := gc.Prog(sparc64.AADD)
		p4.From.Type = obj.TYPE_CONST
		p4.From.Offset = 8
		p4.To.Type = obj.TYPE_REG
		p4.To.Reg = gc.SSARegNum(v.Args[0])
		p5 := gc.Prog(sparc64.ACMP)
		p5.From.Type = obj.TYPE_REG
		p5.From.Reg = gc.SSARegNum(v.Args[2])
		p5.Reg = gc.SSARegNum(v.Args[1])
		p6 := gc.Prog(sparc64.ABLED)
		p6.To.Type = obj.TYPE_BRANCH
		gc.Patch(p6, p)

	case ssa.OpSPARC64CALLstatic:
		if v.Aux.(*gc.Sym) == gc.Deferreturn.Sym {
			// TODO(shawn): is this true on sparc due to pc/npc difference?
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
	case ssa.OpSPARC64CALLclosure:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpSPARC64CALLdefer:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Deferproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpSPARC64CALLgo:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Newproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}
	case ssa.OpSPARC64CALLinter:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 0
		p.To.Reg = gc.SSARegNum(v.Args[0])
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}

	case ssa.OpSPARC64LoweredNilCheck:
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
			case ssa.OpSPARC64MOVBload, ssa.OpSPARC64MOVUBload, ssa.OpSPARC64MOVHload, ssa.OpSPARC64MOVUHload,
				ssa.OpSPARC64FMOVSload, ssa.OpSPARC64FMOVDload,
				ssa.OpSPARC64MOVWload, ssa.OpSPARC64MOVUWload, ssa.OpSPARC64MOVDload,
				ssa.OpSPARC64FMOVSstore, ssa.OpSPARC64FMOVDstore,
				ssa.OpSPARC64MOVBstore, ssa.OpSPARC64MOVHstore, ssa.OpSPARC64MOVWstore, ssa.OpSPARC64MOVDstore:
				// arg0 is ptr, auxint is offset (atomic ops have auxint 0)
				if w.Args[0] == v.Args[0] && w.Aux == nil && w.AuxInt >= 0 && w.AuxInt < minZeroPage {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpSPARC64LoweredZero:
				// arg0 is ptr
				if w.Args[0] == v.Args[0] {
					if gc.Debug_checknil != 0 && int(v.Line) > 1 {
						gc.Warnl(v.Line, "removed nil check")
					}
					return
				}
			case ssa.OpSPARC64LoweredMove:
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
		p := gc.Prog(sparc64.AMOVB)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = sparc64.REG_TMP
		if gc.Debug_checknil != 0 && v.Line > 1 { // v.Line==1 in generated wrappers
			gc.Warnl(v.Line, "generated nil check")
		}
	case ssa.OpSPARC64Equal32,
		ssa.OpSPARC64Equal64,
		ssa.OpSPARC64EqualF,
		ssa.OpSPARC64NotEqual32,
		ssa.OpSPARC64NotEqual64,
		ssa.OpSPARC64NotEqualF,
		ssa.OpSPARC64LessThan32,
		ssa.OpSPARC64LessThan32U,
		ssa.OpSPARC64LessThan64,
		ssa.OpSPARC64LessThan64U,
		ssa.OpSPARC64LessThanF,
		ssa.OpSPARC64LessEqual32,
		ssa.OpSPARC64LessEqual32U,
		ssa.OpSPARC64LessEqual64,
		ssa.OpSPARC64LessEqual64U,
		ssa.OpSPARC64LessEqualF,
		ssa.OpSPARC64GreaterThan32,
		ssa.OpSPARC64GreaterThan32U,
		ssa.OpSPARC64GreaterThan64,
		ssa.OpSPARC64GreaterThan64U,
		ssa.OpSPARC64GreaterThanF,
		ssa.OpSPARC64GreaterEqual32,
		ssa.OpSPARC64GreaterEqual32U,
		ssa.OpSPARC64GreaterEqual64,
		ssa.OpSPARC64GreaterEqual64U,
		ssa.OpSPARC64GreaterEqualF:

		p := gc.Prog(sparc64.AMOVD)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = sparc64.REG_ZR
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

		p = gc.Prog(condOps[v.Op])
		p.From.Type = obj.TYPE_REG
		p.From.Reg = condBits[v.Op]
		p.From3 = &obj.Addr{}
		p.From3.Type = obj.TYPE_CONST
		p.From3.Offset = 1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64CMP,
		ssa.OpSPARC64FCMPS,
		ssa.OpSPARC64FCMPD:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.Reg = gc.SSARegNum(v.Args[0])

	case ssa.OpSPARC64CMPconst:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.AuxInt
		p.Reg = gc.SSARegNum(v.Args[0])

	case ssa.OpSPARC64LoweredGetClosurePtr:
		// Closure pointer is sparc64.REG_CTXT.
		gc.CheckLoweredGetClosurePtr(v)

	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
}

var condBits = map[ssa.Op]int16{
	ssa.OpSPARC64Equal32: sparc64.REG_ICC,
	ssa.OpSPARC64NotEqual32: sparc64.REG_ICC,
	ssa.OpSPARC64LessThan32: sparc64.REG_ICC,
	ssa.OpSPARC64LessThan32U: sparc64.REG_ICC,
	ssa.OpSPARC64LessEqual32: sparc64.REG_ICC,
	ssa.OpSPARC64LessEqual32U: sparc64.REG_ICC,
	ssa.OpSPARC64GreaterThan32: sparc64.REG_ICC,
	ssa.OpSPARC64GreaterThan32U: sparc64.REG_ICC,
	ssa.OpSPARC64GreaterEqual32: sparc64.REG_ICC,
	ssa.OpSPARC64GreaterEqual32U: sparc64.REG_ICC,

	ssa.OpSPARC64Equal64: sparc64.REG_XCC,
	ssa.OpSPARC64NotEqual64: sparc64.REG_XCC,
	ssa.OpSPARC64LessThan64: sparc64.REG_XCC,
	ssa.OpSPARC64LessThan64U: sparc64.REG_XCC,
	ssa.OpSPARC64LessEqual64: sparc64.REG_XCC,
	ssa.OpSPARC64LessEqual64U: sparc64.REG_XCC,
	ssa.OpSPARC64GreaterThan64: sparc64.REG_XCC,
	ssa.OpSPARC64GreaterThan64U: sparc64.REG_XCC,
	ssa.OpSPARC64GreaterEqual64: sparc64.REG_XCC,
	ssa.OpSPARC64GreaterEqual64U: sparc64.REG_XCC,

	ssa.OpSPARC64EqualF: sparc64.REG_FCC0,
	ssa.OpSPARC64NotEqualF: sparc64.REG_FCC0,
	ssa.OpSPARC64LessThanF: sparc64.REG_FCC0,
	ssa.OpSPARC64LessEqualF: sparc64.REG_FCC0,
	ssa.OpSPARC64GreaterThanF: sparc64.REG_FCC0,
	ssa.OpSPARC64GreaterEqualF: sparc64.REG_FCC0,
}

var condOps = map[ssa.Op]obj.As{
	ssa.OpSPARC64Equal32: sparc64.AMOVE,
	ssa.OpSPARC64Equal64: sparc64.AMOVE,
	ssa.OpSPARC64NotEqual32: sparc64.AMOVNE,
	ssa.OpSPARC64NotEqual64: sparc64.AMOVNE,
	ssa.OpSPARC64LessThan32: sparc64.AMOVL,
	ssa.OpSPARC64LessThan64: sparc64.AMOVL,
	ssa.OpSPARC64LessThan32U: sparc64.AMOVCS,
	ssa.OpSPARC64LessThan64U: sparc64.AMOVCS,
	ssa.OpSPARC64LessEqual32: sparc64.AMOVLE,
	ssa.OpSPARC64LessEqual64: sparc64.AMOVLE,
	ssa.OpSPARC64LessEqual32U: sparc64.AMOVLEU,
	ssa.OpSPARC64LessEqual64U: sparc64.AMOVLEU,
	ssa.OpSPARC64GreaterThan32: sparc64.AMOVG,
	ssa.OpSPARC64GreaterThan64: sparc64.AMOVG,
	ssa.OpSPARC64GreaterThan32U: sparc64.AMOVGU,
	ssa.OpSPARC64GreaterThan64U: sparc64.AMOVGU,
	ssa.OpSPARC64GreaterEqual32: sparc64.AMOVGE,
	ssa.OpSPARC64GreaterEqual64: sparc64.AMOVGE,
	ssa.OpSPARC64GreaterEqual32U: sparc64.AMOVCC,
	ssa.OpSPARC64GreaterEqual64U: sparc64.AMOVCC,
	ssa.OpSPARC64EqualF: sparc64.AMOVFE,
	ssa.OpSPARC64NotEqualF: sparc64.AMOVFNE,
	ssa.OpSPARC64LessThanF: sparc64.AMOVFL,
	ssa.OpSPARC64LessEqualF: sparc64.AMOVFLE,
	ssa.OpSPARC64GreaterThanF: sparc64.AMOVFG,
	ssa.OpSPARC64GreaterEqualF: sparc64.AMOVFGE,
}

var blockJump = map[ssa.BlockKind]struct {
	asm, invasm obj.As
}{
	ssa.BlockSPARC64ND:  {sparc64.ABND, obj.AJMP},
	ssa.BlockSPARC64NED:  {sparc64.ABNED, sparc64.ABED},
	ssa.BlockSPARC64ED:  {sparc64.ABED, sparc64.ABNED},
	ssa.BlockSPARC64GD:  {sparc64.ABGD, sparc64.ABLED},
	ssa.BlockSPARC64LED:  {sparc64.ABLED, sparc64.ABGD},
	ssa.BlockSPARC64GED:  {sparc64.ABGED, sparc64.ABLD},
	ssa.BlockSPARC64LD:  {sparc64.ABLD, sparc64.ABGED},
	ssa.BlockSPARC64GUD:  {sparc64.ABGUD, sparc64.ABLEUD},
	ssa.BlockSPARC64LEUD:  {sparc64.ABLEUD, sparc64.ABGUD},
	ssa.BlockSPARC64CCD:  {sparc64.ABCCD, sparc64.ABCSD},
	ssa.BlockSPARC64CSD:  {sparc64.ABCSD, sparc64.ABCCD},
	ssa.BlockSPARC64POSD:  {sparc64.ABPOSD, sparc64.ABNEGD},
	ssa.BlockSPARC64NEGD:  {sparc64.ABNEGD, sparc64.ABPOSD},
	ssa.BlockSPARC64VCD:  {sparc64.ABVCD, sparc64.ABVSD},
	ssa.BlockSPARC64VSD:  {sparc64.ABVSD, sparc64.ABVCD},

	ssa.BlockSPARC64NW:  {sparc64.ABNW, obj.AJMP},
	ssa.BlockSPARC64NEW:  {sparc64.ABNEW, sparc64.ABEW},
	ssa.BlockSPARC64EW:  {sparc64.ABEW, sparc64.ABNEW},
	ssa.BlockSPARC64GW:  {sparc64.ABGW, sparc64.ABLEW},
	ssa.BlockSPARC64LEW:  {sparc64.ABLEW, sparc64.ABGW},
	ssa.BlockSPARC64GEW:  {sparc64.ABGEW, sparc64.ABLW},
	ssa.BlockSPARC64LW:  {sparc64.ABLW, sparc64.ABGEW},
	ssa.BlockSPARC64GUW:  {sparc64.ABGUW, sparc64.ABLEUW},
	ssa.BlockSPARC64LEUW:  {sparc64.ABLEUW, sparc64.ABGUW},
	ssa.BlockSPARC64CCW:  {sparc64.ABCCW, sparc64.ABCSW},
	ssa.BlockSPARC64CSW:  {sparc64.ABCSW, sparc64.ABCCW},
	ssa.BlockSPARC64POSW:  {sparc64.ABPOSW, sparc64.ABNEGW},
	ssa.BlockSPARC64NEGW:  {sparc64.ABNEGW, sparc64.ABPOSW},
	ssa.BlockSPARC64VCW:  {sparc64.ABVCW, sparc64.ABVSW},
	ssa.BlockSPARC64VSW:  {sparc64.ABVSW, sparc64.ABVCW},

	ssa.BlockSPARC64NF: {sparc64.AFBN, obj.AJMP},
	ssa.BlockSPARC64EF: {sparc64.AFBE, sparc64.AFBNE},
	ssa.BlockSPARC64NEF: {sparc64.AFBNE, sparc64.AFBE},
	ssa.BlockSPARC64LF: {sparc64.AFBL, sparc64.AFBGE},
	ssa.BlockSPARC64LEF: {sparc64.AFBLE, sparc64.AFBUG},
	ssa.BlockSPARC64GF: {sparc64.AFBG, sparc64.AFBLE},
	ssa.BlockSPARC64GEF: {sparc64.AFBGE, sparc64.AFBUL},
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
		// defer returns in RT1:
		// 0 if we should continue executing
		// 1 if we should jump to deferreturn call
		p := gc.Prog(sparc64.ACMP)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 0
		p.Reg = sparc64.REG_RT1
		p = gc.Prog(sparc64.ABNEW)
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
		p := gc.Prog(obj.AJMP)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(b.Aux.(*gc.Sym))

	case ssa.BlockSPARC64ND,
		ssa.BlockSPARC64NED,
		ssa.BlockSPARC64ED,
		ssa.BlockSPARC64GD,
		ssa.BlockSPARC64LED,
		ssa.BlockSPARC64GED,
		ssa.BlockSPARC64LD,
		ssa.BlockSPARC64GUD,
		ssa.BlockSPARC64LEUD,
		ssa.BlockSPARC64CCD,
		ssa.BlockSPARC64CSD,
		ssa.BlockSPARC64POSD,
		ssa.BlockSPARC64NEGD,
		ssa.BlockSPARC64VCD,
		ssa.BlockSPARC64VSD,
		ssa.BlockSPARC64NW,
		ssa.BlockSPARC64NEW,
		ssa.BlockSPARC64EW,
		ssa.BlockSPARC64GW,
		ssa.BlockSPARC64LEW,
		ssa.BlockSPARC64GEW,
		ssa.BlockSPARC64LW,
		ssa.BlockSPARC64GUW,
		ssa.BlockSPARC64LEUW,
		ssa.BlockSPARC64CCW,
		ssa.BlockSPARC64CSW,
		ssa.BlockSPARC64POSW,
		ssa.BlockSPARC64NEGW,
		ssa.BlockSPARC64VCW,
		ssa.BlockSPARC64VSW,
		ssa.BlockSPARC64NF,
		ssa.BlockSPARC64EF,
		ssa.BlockSPARC64NEF,
		ssa.BlockSPARC64LF,
		ssa.BlockSPARC64LEF,
		ssa.BlockSPARC64GF,
		ssa.BlockSPARC64GEF:

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

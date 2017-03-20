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
	sparc64.REG_RT2,  // for runtime, linblink and duff device
	// sparc64.REG_TMP,  // reserved for runtime and linblink
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
	// sparc64.REG_TMP2, // reserved for runtime and linblink
	sparc64.REG_L1,
	sparc64.REG_L2,
	sparc64.REG_L3,
	sparc64.REG_L4,
	sparc64.REG_L5,
	sparc64.REG_L6,
	// sparc64.REG_L7,  // reserved for runtime, to debug register windows
	// sparc64.REG_I0,  // unused to debug register windows
	// sparc64.REG_I1,  // unused to debug register windows
	// sparc64.REG_I2,  // unused to debug register windows
	// sparc64.REG_I3,  // unused to debug register windows
	// sparc64.REG_I4,  // unused to debug register windows
	// sparc64.REG_I5,  // unused to debug register windows
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

	0, // SB, pseudo
	1, // SP, pseudo
	2, // FP, pseudo
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

	case ssa.OpSPARC64ADD,
		ssa.OpSPARC64SUB,
		ssa.OpSPARC64AND,
		ssa.OpSPARC64OR,
		ssa.OpSPARC64XOR,
		ssa.OpSPARC64MULD,
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
		ssa.OpSPARC64XORconst:

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
		ssa.OpSPARC64FSQRTD:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

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
		ssa.OpSPARC64MOVDload:
		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

	case ssa.OpSPARC64MOVDstore, ssa.OpSPARC64MOVWstore, ssa.OpSPARC64MOVHstore, ssa.OpSPARC64MOVBstore:

		p := gc.Prog(v.Op.Asm())
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[1])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)

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
	case ssa.OpSPARC64CALLdefer:
		p := gc.Prog(obj.ACALL)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = gc.Linksym(gc.Deferproc.Sym)
		if gc.Maxarg < v.AuxInt {
			gc.Maxarg = v.AuxInt
		}

	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
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

	default:
		b.Unimplementedf("branch not implemented: %s. Control: %s", b.LongString(), b.Control.LongString())
	}
}

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/x86"
	"math"
)

// Generates code for v using 387 instructions.  Reports whether
// the instruction was handled by this routine.
func ssaGenValue387(s *gc.SSAGenState, v *ssa.Value) bool {
	// The SSA compiler pretends that it has an SSE backend.
	// If we don't have one of those, we need to translate
	// all the SSE ops to equivalent 387 ops. That's what this
	// function does.

	switch v.Op {
	case ssa.Op386MOVSSconst, ssa.Op386MOVSDconst:
		p := gc.Prog(loadPush(v.Type))
		p.From.Type = obj.TYPE_FCONST
		p.From.Val = math.Float64frombits(uint64(v.AuxInt))
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)
		return true
	case ssa.Op386MOVSSconst2, ssa.Op386MOVSDconst2:
		p := gc.Prog(loadPush(v.Type))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)
		return true

	case ssa.Op386MOVSSload, ssa.Op386MOVSDload, ssa.Op386MOVSSloadidx1, ssa.Op386MOVSDloadidx1, ssa.Op386MOVSSloadidx4, ssa.Op386MOVSDloadidx8:
		p := gc.Prog(loadPush(v.Type))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.From, v)
		switch v.Op {
		case ssa.Op386MOVSSloadidx1, ssa.Op386MOVSDloadidx1:
			p.From.Scale = 1
			p.From.Index = gc.SSARegNum(v.Args[1])
		case ssa.Op386MOVSSloadidx4:
			p.From.Scale = 4
			p.From.Index = gc.SSARegNum(v.Args[1])
		case ssa.Op386MOVSDloadidx8:
			p.From.Scale = 8
			p.From.Index = gc.SSARegNum(v.Args[1])
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)
		return true

	case ssa.Op386MOVSSstore, ssa.Op386MOVSDstore:
		// Push to-be-stored value on top of stack.
		push(s, v.Args[1])

		// Pop and store value.
		var op obj.As
		switch v.Op {
		case ssa.Op386MOVSSstore:
			op = x86.AFMOVFP
		case ssa.Op386MOVSDstore:
			op = x86.AFMOVDP
		}
		p := gc.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
		return true

	case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSDstoreidx1, ssa.Op386MOVSSstoreidx4, ssa.Op386MOVSDstoreidx8:
		push(s, v.Args[2])
		var op obj.As
		switch v.Op {
		case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSSstoreidx4:
			op = x86.AFMOVFP
		case ssa.Op386MOVSDstoreidx1, ssa.Op386MOVSDstoreidx8:
			op = x86.AFMOVDP
		}
		p := gc.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = gc.SSARegNum(v.Args[0])
		gc.AddAux(&p.To, v)
		switch v.Op {
		case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSDstoreidx1:
			p.To.Scale = 1
			p.To.Index = gc.SSARegNum(v.Args[1])
		case ssa.Op386MOVSSstoreidx4:
			p.To.Scale = 4
			p.To.Index = gc.SSARegNum(v.Args[1])
		case ssa.Op386MOVSDstoreidx8:
			p.To.Scale = 8
			p.To.Index = gc.SSARegNum(v.Args[1])
		}
		return true

	case ssa.Op386ADDSS, ssa.Op386ADDSD, ssa.Op386SUBSS, ssa.Op386SUBSD,
		ssa.Op386MULSS, ssa.Op386MULSD, ssa.Op386DIVSS, ssa.Op386DIVSD:
		if gc.SSARegNum(v) != gc.SSARegNum(v.Args[0]) {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}

		// Push arg1 on top of stack
		push(s, v.Args[1])

		// Set precision if needed.  64 bits is the default.
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386SUBSS, ssa.Op386MULSS, ssa.Op386DIVSS:
			p := gc.Prog(x86.AFSTCW)
			scratch387(s, &p.To)
			p = gc.Prog(x86.AFLDCW)
			p.From.Type = obj.TYPE_MEM
			p.From.Name = obj.NAME_EXTERN
			p.From.Sym = gc.Linksym(gc.Pkglookup("controlWord32", gc.Runtimepkg))
		}

		var op obj.As
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386ADDSD:
			op = x86.AFADDDP
		case ssa.Op386SUBSS, ssa.Op386SUBSD:
			op = x86.AFSUBDP
		case ssa.Op386MULSS, ssa.Op386MULSD:
			op = x86.AFMULDP
		case ssa.Op386DIVSS, ssa.Op386DIVSD:
			op = x86.AFDIVDP
		}
		p := gc.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[gc.SSARegNum(v)] + 1

		// Restore precision if needed.
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386SUBSS, ssa.Op386MULSS, ssa.Op386DIVSS:
			p := gc.Prog(x86.AFLDCW)
			scratch387(s, &p.From)
		}

		return true

	case ssa.Op386UCOMISS, ssa.Op386UCOMISD:
		push(s, v.Args[0])

		// Compare.
		p := gc.Prog(x86.AFUCOMP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[gc.SSARegNum(v.Args[1])] + 1

		// Save AX.
		p = gc.Prog(x86.AMOVL)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_AX
		scratch387(s, &p.To)

		// Move status word into AX.
		p = gc.Prog(x86.AFSTSW)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_AX

		// Then move the flags we need to the integer flags.
		gc.Prog(x86.ASAHF)

		// Restore AX.
		p = gc.Prog(x86.AMOVL)
		scratch387(s, &p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_AX

		return true

	case ssa.Op386SQRTSD:
		push(s, v.Args[0])
		gc.Prog(x86.AFSQRT)
		popAndSave(s, v)
		return true

	case ssa.Op386FCHS:
		push(s, v.Args[0])
		gc.Prog(x86.AFCHS)
		popAndSave(s, v)
		return true

	case ssa.Op386CVTSL2SS, ssa.Op386CVTSL2SD:
		p := gc.Prog(x86.AMOVL)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = gc.SSARegNum(v.Args[0])
		scratch387(s, &p.To)
		p = gc.Prog(x86.AFMOVL)
		scratch387(s, &p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)
		return true

	case ssa.Op386CVTTSD2SL, ssa.Op386CVTTSS2SL:
		push(s, v.Args[0])

		// Save control word.
		p := gc.Prog(x86.AFSTCW)
		scratch387(s, &p.To)
		p.To.Offset += 4

		// Load control word which truncates (rounds towards zero).
		p = gc.Prog(x86.AFLDCW)
		p.From.Type = obj.TYPE_MEM
		p.From.Name = obj.NAME_EXTERN
		p.From.Sym = gc.Linksym(gc.Pkglookup("controlWord64trunc", gc.Runtimepkg))

		// Now do the conversion.
		p = gc.Prog(x86.AFMOVLP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		scratch387(s, &p.To)
		p = gc.Prog(x86.AMOVL)
		scratch387(s, &p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = gc.SSARegNum(v)

		// Restore control word.
		p = gc.Prog(x86.AFLDCW)
		scratch387(s, &p.From)
		p.From.Offset += 4
		return true

	case ssa.Op386CVTSS2SD:
		// float32 -> float64 is a nop
		push(s, v.Args[0])
		popAndSave(s, v)
		return true

	case ssa.Op386CVTSD2SS:
		// Round to nearest float32.
		push(s, v.Args[0])
		p := gc.Prog(x86.AFMOVFP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		scratch387(s, &p.To)
		p = gc.Prog(x86.AFMOVF)
		scratch387(s, &p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)
		return true

	case ssa.OpLoadReg:
		if !v.Type.IsFloat() {
			return false
		}
		// Load+push the value we need.
		p := gc.Prog(loadPush(v.Type))
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
		p.To.Reg = x86.REG_F0
		// Move the value to its assigned register.
		popAndSave(s, v)
		return true

	case ssa.OpStoreReg:
		if !v.Type.IsFloat() {
			return false
		}
		push(s, v.Args[0])
		var op obj.As
		switch v.Type.Size() {
		case 4:
			op = x86.AFMOVFP
		case 8:
			op = x86.AFMOVDP
		}
		p := gc.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
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
		return true

	case ssa.OpCopy:
		if !v.Type.IsFloat() {
			return false
		}
		push(s, v.Args[0])
		popAndSave(s, v)
		return true

	case ssa.Op386CALLstatic, ssa.Op386CALLclosure, ssa.Op386CALLdefer, ssa.Op386CALLgo, ssa.Op386CALLinter:
		flush387(s)  // Calls must empty the the FP stack.
		return false // then issue the call as normal
	}
	return false
}

// push pushes v onto the floating-point stack.  v must be in a register.
func push(s *gc.SSAGenState, v *ssa.Value) {
	p := gc.Prog(x86.AFMOVD)
	p.From.Type = obj.TYPE_REG
	p.From.Reg = s.SSEto387[gc.SSARegNum(v)]
	p.To.Type = obj.TYPE_REG
	p.To.Reg = x86.REG_F0
}

// popAndSave pops a value off of the floating-point stack and stores
// it in the reigster assigned to v.
func popAndSave(s *gc.SSAGenState, v *ssa.Value) {
	r := gc.SSARegNum(v)
	if _, ok := s.SSEto387[r]; ok {
		// Pop value, write to correct register.
		p := gc.Prog(x86.AFMOVDP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[gc.SSARegNum(v)] + 1
	} else {
		// Don't actually pop value. This 387 register is now the
		// new home for the not-yet-assigned-a-home SSE register.
		// Increase the register mapping of all other registers by one.
		for rSSE, r387 := range s.SSEto387 {
			s.SSEto387[rSSE] = r387 + 1
		}
		s.SSEto387[r] = x86.REG_F0
	}
}

// loadPush returns the opcode for load+push of the given type.
func loadPush(t ssa.Type) obj.As {
	if t.Size() == 4 {
		return x86.AFMOVF
	}
	return x86.AFMOVD
}

// flush387 removes all entries from the 387 floating-point stack.
func flush387(s *gc.SSAGenState) {
	for k := range s.SSEto387 {
		p := gc.Prog(x86.AFMOVDP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		delete(s.SSEto387, k)
	}
}

// scratch387 initializes a to the scratch location used by some 387 rewrites.
func scratch387(s *gc.SSAGenState, a *obj.Addr) {
	a.Type = obj.TYPE_MEM
	a.Name = obj.NAME_AUTO
	a.Node = s.ScratchFpMem
	a.Sym = gc.Linksym(s.ScratchFpMem.Sym)
	a.Reg = x86.REG_SP
}

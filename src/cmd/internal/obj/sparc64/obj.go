// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"cmd/internal/sys"
	"fmt"
	"log"
	"math"
)

var isUncondJump = map[obj.As]bool{
	obj.ACALL:     true,
	obj.ADUFFZERO: true,
	obj.ADUFFCOPY: true,
	obj.AJMP:      true,
	obj.ARET:      true,
	AFBA:          true,
	AJMPL:         true,
	ARETRESTORE:   true,
}

var isCondJump = map[obj.As]bool{
	ABN:    true,
	ABNE:   true,
	ABE:    true,
	ABG:    true,
	ABLE:   true,
	ABGE:   true,
	ABL:    true,
	ABGU:   true,
	ABLEU:  true,
	ABCC:   true,
	ABCS:   true,
	ABPOS:  true,
	ABNEG:  true,
	ABVC:   true,
	ABVS:   true,
	ABNW:   true,
	ABNEW:  true,
	ABEW:   true,
	ABGW:   true,
	ABLEW:  true,
	ABGEW:  true,
	ABLW:   true,
	ABGUW:  true,
	ABLEUW: true,
	ABCCW:  true,
	ABCSW:  true,
	ABPOSW: true,
	ABNEGW: true,
	ABVCW:  true,
	ABVSW:  true,
	ABND:   true,
	ABNED:  true,
	ABED:   true,
	ABGD:   true,
	ABLED:  true,
	ABGED:  true,
	ABLD:   true,
	ABGUD:  true,
	ABLEUD: true,
	ABCCD:  true,
	ABCSD:  true,
	ABPOSD: true,
	ABNEGD: true,
	ABVCD:  true,
	ABVSD:  true,
	ABRZ:   true,
	ABRLEZ: true,
	ABRLZ:  true,
	ABRNZ:  true,
	ABRGZ:  true,
	ABRGEZ: true,
	AFBN:   true,
	AFBU:   true,
	AFBG:   true,
	AFBUG:  true,
	AFBL:   true,
	AFBUL:  true,
	AFBLG:  true,
	AFBNE:  true,
	AFBE:   true,
	AFBUE:  true,
	AFBGE:  true,
	AFBUGE: true,
	AFBLE:  true,
	AFBULE: true,
	AFBO:   true,
}

var isJump = make(map[obj.As]bool)

func init() {
	for k := range isUncondJump {
		isJump[k] = true
	}
	for k := range isCondJump {
		isJump[k] = true
	}
}

func stacksplit(ctxt *obj.Link, p *obj.Prog, framesize int32) *obj.Prog {
	//TODO(shawn): re-enable stacksplit once scanstack is supported
	if true {
		return p
	}

	// MOV	g_stackguard(g), L1
	p = obj.Appendp(ctxt, p)
	start := p

	p.As = AMOVD
	p.From.Type = obj.TYPE_MEM
	p.From.Reg = REG_G
	p.From.Offset = 2 * int64(ctxt.Arch.PtrSize) // G.stackguard0
	if ctxt.Cursym.Cfunc {
		p.From.Offset = 3 * int64(ctxt.Arch.PtrSize) // G.stackguard1
	}
	p.To.Type = obj.TYPE_REG
	p.To.Reg = REG_L1

	q := (*obj.Prog)(nil)
	if framesize <= obj.StackSmall {
		// small stack: SP-StackBias < stackguard
		//	CMP	stackguard, SP
		p = obj.Appendp(ctxt, p)

		p.As = ACMP
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_L1
		p.Reg = REG_RSP
	} else if framesize <= obj.StackBig {
		// large stack: SP-framesize < stackguard-StackSmall
		//	SUB	$framesize, RSP, L2
		//	CMP	stackguard, L2
		p = obj.Appendp(ctxt, p)

		p.As = ASUB
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(framesize)
		p.Reg = REG_RSP
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_L2

		p = obj.Appendp(ctxt, p)
		p.As = ACMP
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_L1
		p.Reg = REG_L2
	} else {
		// Such a large stack we need to protect against wraparound
		// if SP is close to zero.
		//	SP-stackguard+StackGuard < framesize + (StackGuard-StackSmall)
		// The +StackGuard on both sides is required to keep the left side positive:
		// SP is allowed to be slightly below stackguard. See stack.h.
		//	CMP	$StackPreempt, L1
		//	BED	label_of_call_to_morestack
		//	ADD	$StackGuard, RSP, L2
		//	SUB	L1, L2
		//	MOV	$(framesize+(StackGuard-StackSmall)), L3
		//	CMP	L3, L2
		p = obj.Appendp(ctxt, p)

		p.As = ACMP
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = obj.StackPreempt
		p.Reg = REG_L1

		p = obj.Appendp(ctxt, p)
		q = p
		p.As = ABED
		p.To.Type = obj.TYPE_BRANCH

		p = obj.Appendp(ctxt, p)
		p.As = AADD
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = obj.StackGuard
		p.Reg = REG_RSP
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_L2

		p = obj.Appendp(ctxt, p)
		p.As = ASUB
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_L1
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_L2

		p = obj.Appendp(ctxt, p)
		p.As = AMOVD
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(framesize) + (obj.StackGuard - obj.StackSmall)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_L3

		p = obj.Appendp(ctxt, p)
		p.As = ACMP
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_L3
		p.Reg = REG_L2
	}

	// BLE	do-morestack
	ble := obj.Appendp(ctxt, p)
	ble.As = ABLED
	ble.To.Type = obj.TYPE_BRANCH

	var last *obj.Prog
	for last = ctxt.Cursym.Text; last.Link != nil; last = last.Link {
	}

	spfix := obj.Appendp(ctxt, last)
	spfix.As = ARNOP
	spfix.Spadj = -framesize

	// MOV	LR, L3
	movlr := obj.Appendp(ctxt, spfix)
	movlr.As = AMOVD
	movlr.From.Type = obj.TYPE_REG
	movlr.From.Reg = REG_ILR
	movlr.To.Type = obj.TYPE_REG
	movlr.To.Reg = REG_L3
	if q != nil {
		q.Pcond = movlr
	}
	ble.Pcond = movlr

	debug := movlr
	if true {
		debug = obj.Appendp(ctxt, debug)
		debug.As = AMOVD
		debug.From.Type = obj.TYPE_CONST
		debug.From.Offset = int64(framesize)
		debug.To.Type = obj.TYPE_REG
		debug.To.Reg = REG_TMP
	}

	// CALL runtime.morestack(SB)
	call := obj.Appendp(ctxt, debug)
	call.As = obj.ACALL
	call.To.Type = obj.TYPE_MEM
	call.To.Name = obj.NAME_EXTERN
	morestack := "runtime.morestack"
	switch {
	case ctxt.Cursym.Cfunc:
		morestack = "runtime.morestackc"
	case ctxt.Cursym.Text.From3.Offset&obj.NEEDCTXT == 0:
		morestack = "runtime.morestack_noctxt"
	}
	call.To.Sym = obj.Linklookup(ctxt, morestack, 0)

	// JMP start
	jmp := obj.Appendp(ctxt, call)
	jmp.As = obj.AJMP
	jmp.To.Type = obj.TYPE_BRANCH
	jmp.Pcond = start
	jmp.Spadj = +framesize

	return ble
}

// AutoeditProg returns a new obj.Prog, with off(SP), off(FP), $off(SP),
// and $off(FP) replaced with new(RFP).
func autoeditprog(ctxt *obj.Link, p *obj.Prog) *obj.Prog {
	r := new(obj.Prog)
	*r = *p
	r.From = *autoeditaddr(ctxt, &r.From)
	r.From3 = autoeditaddr(ctxt, r.From3)
	r.To = *autoeditaddr(ctxt, &r.To)
	return r
}

// Autoeditaddr returns a new obj.Addr, with off(SP), off(FP), $off(SP),
// and $off(FP) replaced with new(RFP).
func autoeditaddr(ctxt *obj.Link, a *obj.Addr) *obj.Addr {
	if a == nil {
		return nil
	}
	if a.Type != obj.TYPE_MEM && a.Type != obj.TYPE_ADDR {
		return a
	}
	r := new(obj.Addr)
	*r = *a
	if r.Name == obj.NAME_PARAM {
		r.Reg = REG_RFP
		if ctxt.Cursym.Text.From3Offset()&obj.NOFRAME != 0 {
			// NOFRAME functions live in caller's frame.
			r.Reg = REG_RSP
		}
		r.Offset += MinStackFrameSize + StackBias
		r.Name = obj.NAME_NONE
		return r
	}
	if r.Name == obj.NAME_AUTO {
		r.Reg = REG_RFP
		r.Offset += StackBias
		r.Name = obj.NAME_NONE
	}
	return r
}

// yfix rewrites references to Y registers (issued by compiler)
// to F and D registers.
func yfix(p *obj.Prog) {
	if REG_Y0 <= p.From.Reg && p.From.Reg <= REG_Y15 {
		if isInstDouble[p.As] || isSrcDouble[p.As] {
			p.From.Reg = REG_D0 + (p.From.Reg-REG_Y0)*2
		} else if isInstFloat[p.As] || isSrcFloat[p.As] {
			p.From.Reg = REG_F0 + (p.From.Reg-REG_Y0)*2
		}
	}
	if REG_Y0 <= p.Reg && p.Reg <= REG_Y15 {
		if isInstDouble[p.As] {
			p.Reg = REG_D0 + (p.Reg-REG_Y0)*2
		} else {
			p.Reg = REG_F0 + (p.Reg-REG_Y0)*2
		}
	}
	if p.From3 != nil && REG_Y0 <= p.From3.Reg && p.From3.Reg <= REG_Y15 {
		if isInstDouble[p.As] {
			p.From3.Reg = REG_D0 + (p.From3.Reg-REG_Y0)*2
		} else {
			p.From3.Reg = REG_F0 + (p.From3.Reg-REG_Y0)*2
		}
	}
	if REG_Y0 <= p.To.Reg && p.To.Reg <= REG_Y15 {
		if isInstDouble[p.As] || isDstDouble[p.As] {
			p.To.Reg = REG_D0 + (p.To.Reg-REG_Y0)*2
		} else if isInstFloat[p.As] || isDstFloat[p.As] {
			p.To.Reg = REG_F0 + (p.To.Reg-REG_Y0)*2
		}
	}
}

// biasfix rewrites referencing to BSP and BFP to RSP and RFP and
// adding the stack bias.
func biasfix(p *obj.Prog) {
	// Only match 2-operand instructions.
	if p.From3 != nil || p.Reg != 0 {
		return
	}
	switch p.As {
	case AMOVD:
		switch aclass(p.Ctxt, &p.From) {
		case ClassReg, ClassZero:
			switch {
			// MOVD	R, BSP	-> ADD	-$STACK_BIAS, R, RSP
			case aclass(p.Ctxt, &p.To) == ClassReg|ClassBias:
				p.As = AADD
				p.Reg = p.From.Reg
				if p.From.Type == obj.TYPE_CONST {
					p.Reg = REG_ZR
				}
				p.From.Reg = 0
				p.From.Offset = -StackBias
				p.From.Type = obj.TYPE_CONST
				p.From.Class = aclass(p.Ctxt, &p.From)
				p.To.Reg -= 256 // must match a.out.go:/REG_BSP
				p.To.Class = aclass(p.Ctxt, &p.To)
			}

		case ClassReg | ClassBias:
			// MOVD	BSP, R	-> ADD	$STACK_BIAS, RSP, R
			if aclass(p.Ctxt, &p.To) == ClassReg {
				p.Reg = p.From.Reg - 256 // must match a.out.go:/REG_BSP
				p.As = AADD
				p.From.Reg = 0
				p.From.Offset = StackBias
				p.From.Type = obj.TYPE_CONST
				p.From.Class = aclass(p.Ctxt, &p.From)
			}

		// MOVD	$off(BSP), R	-> MOVD	$(off+STACK_BIAS)(RSP), R
		case ClassRegConst13 | ClassBias, ClassRegConst | ClassBias:
			p.From.Reg -= 256 // must match a.out.go:/REG_BSP
			p.From.Offset += StackBias
			p.From.Class = aclass(p.Ctxt, &p.From)
		}

	case AADD, ASUB:
		// ADD	$const, BSP	-> ADD	$const, RSP
		if isAddrCompatible(p.Ctxt, &p.From, ClassConst) && aclass(p.Ctxt, &p.To) == ClassReg|ClassBias {
			p.To.Reg -= 256 // must match a.out.go:/REG_BSP
			p.To.Class = aclass(p.Ctxt, &p.To)
		}
	}
	switch p.As {
	case AMOVD, AMOVW, AMOVUW, AMOVH, AMOVUH, AMOVB, AMOVUB,
		AFMOVD, AFMOVS:
		switch aclass(p.Ctxt, &p.From) {
		case ClassZero, ClassReg, ClassFReg, ClassDReg:
			switch {
			// MOVD	R, off(BSP)	-> MOVD	R, (off+STACK_BIAS)(RSP)
			case aclass(p.Ctxt, &p.To)&ClassBias != 0 && isAddrCompatible(p.Ctxt, &p.To, ClassIndir):
				p.To.Offset += StackBias
				p.To.Reg -= 256 // must match a.out.go:/REG_BSP
				p.To.Class = aclass(p.Ctxt, &p.To)
			}

		// MOVD	off(BSP), R	-> MOVD	(off+STACK_BIAS)(RSP), R
		case ClassIndir0 | ClassBias, ClassIndir13 | ClassBias, ClassIndir | ClassBias:
			p.From.Reg -= 256 // must match a.out.go:/REG_BSP
			p.From.Offset += StackBias
			p.From.Class = aclass(p.Ctxt, &p.From)
		}
	}
}

func progedit(ctxt *obj.Link, p *obj.Prog) {
	// Rewrite constant moves to memory to go through an intermediary
	// register
	switch p.As {
	case AMOVD:
		if (p.From.Type == obj.TYPE_CONST || p.From.Type == obj.TYPE_ADDR) && (p.To.Type == obj.TYPE_MEM) {
			q := obj.Appendp(ctxt, p)
			q.As = p.As
			q.To = p.To
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_TMP
			q.From.Offset = 0

			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_TMP
			p.To.Offset = 0
		}

	case AFMOVS:
		if (p.From.Type == obj.TYPE_FCONST || p.From.Type == obj.TYPE_ADDR) && (p.To.Type == obj.TYPE_MEM) {
			q := obj.Appendp(ctxt, p)
			q.As = p.As
			q.To = p.To
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_FTMP
			q.From.Offset = 0

			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_FTMP
			p.To.Offset = 0
		}

	case AFMOVD:
		if (p.From.Type == obj.TYPE_FCONST || p.From.Type == obj.TYPE_ADDR) && (p.To.Type == obj.TYPE_MEM) {
			q := obj.Appendp(ctxt, p)
			q.As = p.As
			q.To = p.To
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_DTMP
			q.From.Offset = 0

			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_DTMP
			p.To.Offset = 0
		}
	}

	// Rewrite 64-bit integer constants and float constants
	// to values stored in memory.
	switch p.As {
	case AMOVD:
		if aclass(p.Ctxt, &p.From) == ClassConst {
			literal := fmt.Sprintf("$i64.%016x", uint64(p.From.Offset))
			s := obj.Linklookup(ctxt, literal, 0)
			s.Size = 8
			p.From.Type = obj.TYPE_MEM
			p.From.Sym = s
			p.From.Name = obj.NAME_EXTERN
			p.From.Offset = 0
		}

	case AFMOVS:
		if p.From.Type == obj.TYPE_FCONST {
			f32 := float32(p.From.Val.(float64))
			i32 := math.Float32bits(f32)
			literal := fmt.Sprintf("$f32.%08x", uint32(i32))
			s := obj.Linklookup(ctxt, literal, 0)
			s.Size = 4
			p.From.Type = obj.TYPE_MEM
			p.From.Sym = s
			p.From.Name = obj.NAME_EXTERN
			p.From.Offset = 0
		}

	case AFMOVD:
		if p.From.Type == obj.TYPE_FCONST {
			i64 := math.Float64bits(p.From.Val.(float64))
			literal := fmt.Sprintf("$f64.%016x", uint64(i64))
			s := obj.Linklookup(ctxt, literal, 0)
			s.Size = 8
			p.From.Type = obj.TYPE_MEM
			p.From.Sym = s
			p.From.Name = obj.NAME_EXTERN
			p.From.Offset = 0
		}
	}

	// TODO(aram): remove this when compiler can use F and
	// D registers directly.
	yfix(p)

	biasfix(p)
}

func isNOFRAME(p *obj.Prog) bool {
	return p.From3Offset()&obj.NOFRAME != 0
}

// TODO(aram):
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {
	cursym.Text.Pc = 0
	cursym.Args = cursym.Text.To.Val.(int32)
	cursym.Locals = int32(cursym.Text.To.Offset)

	// Find leaf subroutines,
	// Strip NOPs.
	var q *obj.Prog
	var q1 *obj.Prog
	for p := cursym.Text; p != nil; p = p.Link {
		switch {
		case p.As == obj.ATEXT:
			p.Mark |= LEAF

		case p.As == obj.ARET:
			if cursym.Text.Mark&LEAF != 0 && p.To.Sym != nil { // RETJMP
				cursym.Text.From3.Offset |= obj.NOFRAME
			}
			break

		case p.As == obj.ANOP:
			q1 = p.Link
			q.Link = q1 /* q is non-nop */
			q1.Mark |= p.Mark
			continue

		case isUncondJump[p.As]:
			cursym.Text.Mark &^= LEAF
			fallthrough

		case isCondJump[p.As]:
			q1 = p.Pcond

			if q1 != nil {
				for q1.As == obj.ANOP {
					q1 = q1.Link
					p.Pcond = q1
				}
			}

			break
		}

		q = p
	}

	for p := cursym.Text; p != nil; p = p.Link {
		switch p.As {
		case obj.ATEXT:
			if cursym.Text.Mark&LEAF != 0 {
				cursym.Leaf = true
			}
		}
	}

	for p := cursym.Text; p != nil; p = p.Link {
		switch p.As {
		case obj.ATEXT:
			frameSize := cursym.Locals
			if frameSize < 0 {
				ctxt.Diag("%v: negative frame size %d", p, frameSize)
			}
			if frameSize%16 != 0 {
				ctxt.Diag("%v: unaligned frame size %d - must be 0 mod 16", p, frameSize)
			}
			noFrame := isNOFRAME(p)
			if frameSize != 0 && noFrame {
				ctxt.Diag("%v: non-zero framesize for NOFRAME function", p)
			}

			// Without these NOPs, DTrace changes the execution of the binary,
			// This should never happen, but these NOPs seems to fix it.
			// Keep these NOPs in here until we understand the DTrace behavior.
			p = obj.Appendp(ctxt, p)
			p.As = ARNOP
			p = obj.Appendp(ctxt, p)
			p.As = ARNOP
			if noFrame {
				break
			}

			// MOVD	$-(frameSize+MinStackFrameSize), RT1
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = -int64(frameSize + MinStackFrameSize)
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_RT1

			// SAVE	RT1, RSP
			p = obj.Appendp(ctxt, p)
			p.As = ASAVE
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_RT1
			p.Reg = REG_RSP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_RSP
			p.Spadj = frameSize + MinStackFrameSize

			// MOVD LR, (120+bias)(RSP)
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_ILR
			p.To.Type = obj.TYPE_MEM
			p.To.Reg = REG_RSP
			p.To.Offset = int64(120 + StackBias)

			// MOVD RFP, (112+bias)(RSP)
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_RFP
			p.To.Type = obj.TYPE_MEM
			p.To.Reg = REG_RSP
			p.To.Offset = int64(112 + StackBias)

			if cursym.Args == obj.ArgsSizeUnknown {
				break
			}
			args := int(cursym.Args) / 8
			// At this point we have a stack frame that
			// can be unwinded and register window spilled
			// by native target code. However, the native
			// target code (e.g. mdb(1)) doesn't know where
			// to look for function arguments. Since we
			// have argument information, we copy the
			// arguments in the place where native tools
			// expect them.

			// For all functions copy at most 6 arguments into the
			// %i registers.
			if args > 6 {
				args = 6
			}
			for i := 0; i < args; i++ {
				// MOVD	argN+8N(BFP), %iN
				p = obj.Appendp(ctxt, p)
				p.As = AMOVD
				p.From.Type = obj.TYPE_MEM
				p.From.Reg = REG_RFP
				p.From.Offset = int64(i*8 + MinStackFrameSize + StackBias)
				p.To.Type = obj.TYPE_REG
				p.To.Reg = int16(REG_I0 + i)
			}

			if !(cursym.Text.From3.Offset&obj.NOSPLIT != 0) {
				p = stacksplit(ctxt, p, frameSize) // emit split check
			}

			if cursym.Text.From3.Offset&obj.WRAPPER != 0 {
				// if(g->panic != nil && g->panic->argp == FP) g->panic->argp = bottom-of-frame
				//
				//	MOVD	g_panic(g), L1
				//	CMP	ZR, L1
				//	BED	end
				//	MOVD	panic_argp(L1), L2
				//	ADD	$(STACK_BIAS+MinStackFrameSize), RFP, L3
				//	CMP	L2, L3
				//	BNED	end
				//	ADD	$(STACK_BIAS+MinStackFrameSize), RSP, L4
				//	MOVD	L4, panic_argp(L1)
				// end:
				//	RNOP
				//
				// The RNOP is needed to give the jumps somewhere to land.
				q = obj.Appendp(ctxt, p)
				q.As = AMOVD
				q.From.Type = obj.TYPE_MEM
				q.From.Reg = REG_G
				q.From.Offset = 4 * int64(ctxt.Arch.PtrSize) // G.panic
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_L1

				q = obj.Appendp(ctxt, q)
				q.As = ACMP
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_ZR
				q.Reg = REG_L1

				q = obj.Appendp(ctxt, q)
				q.As = ABED
				q.To.Type = obj.TYPE_BRANCH
				q1 := q

				q = obj.Appendp(ctxt, q)
				q.As = AMOVD
				q.From.Type = obj.TYPE_MEM
				q.From.Reg = REG_L1
				q.From.Offset = 0 // Panic.argp
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_L2

				q = obj.Appendp(ctxt, q)
				q.As = AADD
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = StackBias + MinStackFrameSize
				q.Reg = REG_RFP
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_L3

				q = obj.Appendp(ctxt, q)
				q.As = ACMP
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_L2
				q.Reg = REG_L3

				q = obj.Appendp(ctxt, q)
				q.As = ABNED
				q.To.Type = obj.TYPE_BRANCH
				q2 := q

				q = obj.Appendp(ctxt, q)
				q.As = AADD
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = StackBias + MinStackFrameSize
				q.Reg = REG_RSP
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_L4

				q = obj.Appendp(ctxt, q)
				q.As = AMOVD
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_L4
				q.To.Type = obj.TYPE_MEM
				q.To.Reg = REG_L1
				q.To.Offset = 0 // Panic.argp

				q = obj.Appendp(ctxt, q)
				q.As = ARNOP
				q1.Pcond = q
				q2.Pcond = q
			}

		case obj.ARET:
			if isNOFRAME(cursym.Text) {
				if p.To.Sym != nil { // RETJMP
					p.As = obj.AJMP
				}
				break
			}

			// MOVD (112+bias)(RSP), RFP
			q = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_MEM
			p.From.Reg = REG_RSP
			p.From.Offset = int64(112 + StackBias)
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_RFP

			q.As = ARETRESTORE
			q.From.Type = obj.TYPE_NONE
			q.From.Offset = 0
			q.Reg = 0
			q.To.Type = obj.TYPE_NONE
			q.To.Reg = 0
			// The SP restore operation needs a Spadj of
			// -(cursym.Locals + MinStackFrameSize),
			// and the JMP operation needs a Spadj of
			// +(cursym.Locals + MinStackFrameSize).
			//
			// Since this operation does both, they cancel out
			// so we don't do any Spadj adjustment.
			//
			// The best solution would be to split RETRESTORE
			// into the constituent instructions, but that requires
			// more sophisticated delay-slot processing,
			// since the RESTORE has to be in the delay
			// slot of the branch.

		case AADD, ASUB:
			if p.To.Type == obj.TYPE_REG && p.To.Reg == REG_BSP && p.From.Type == obj.TYPE_CONST {
				if p.As == AADD {
					p.Spadj = int32(-p.From.Offset)
				} else {
					p.Spadj = int32(+p.From.Offset)
				}
			}
		}
	}

	// Schedule delay-slots. Only RNOPs for now.
	for p := cursym.Text; p != nil; p = p.Link {
		if !isJump[p.As] || p.As == ARETRESTORE {
			continue
		}
		if p.Link != nil && p.Link.As == ARNOP {
			continue
		}
		p = obj.Appendp(ctxt, p)
		p.As = ARNOP
	}

	// For future use by oplook and friends.
	for p := cursym.Text; p != nil; p = p.Link {
		p.From.Class = aclass(ctxt, &p.From)
		if p.From3 != nil {
			p.From3.Class = aclass(ctxt, p.From3)
		}
		p.To.Class = aclass(ctxt, &p.To)
	}
}

func relinv(a obj.As) obj.As {
	switch a {
	case obj.AJMP:
		return ABN
	case ABN:
		return obj.AJMP
	case ABE:
		return ABNE
	case ABNE:
		return ABE
	case ABG:
		return ABLE
	case ABLE:
		return ABG
	case ABGE:
		return ABL
	case ABL:
		return ABGE
	case ABGU:
		return ABLEU
	case ABLEU:
		return ABGU
	case ABCC:
		return ABCS
	case ABCS:
		return ABCC
	case ABPOS:
		return ABNEG
	case ABNEG:
		return ABPOS
	case ABVC:
		return ABVS
	case ABVS:
		return ABVC
	case ABNW:
		return obj.AJMP
	case ABEW:
		return ABNEW
	case ABNEW:
		return ABEW
	case ABGW:
		return ABLEW
	case ABLEW:
		return ABGW
	case ABGEW:
		return ABLW
	case ABLW:
		return ABGEW
	case ABGUW:
		return ABLEUW
	case ABLEUW:
		return ABGUW
	case ABCCW:
		return ABCSW
	case ABCSW:
		return ABCCW
	case ABPOSW:
		return ABNEGW
	case ABNEGW:
		return ABPOSW
	case ABVCW:
		return ABVSW
	case ABVSW:
		return ABVCW
	case ABND:
		return obj.AJMP
	case ABED:
		return ABNED
	case ABNED:
		return ABED
	case ABGD:
		return ABLED
	case ABLED:
		return ABGD
	case ABGED:
		return ABLD
	case ABLD:
		return ABGED
	case ABGUD:
		return ABLEUD
	case ABLEUD:
		return ABGUD
	case ABCCD:
		return ABCSD
	case ABCSD:
		return ABCCD
	case ABPOSD:
		return ABNEGD
	case ABNEGD:
		return ABPOSD
	case ABVCD:
		return ABVSD
	case ABVSD:
		return ABVCD
	case AFBN:
		return AFBA
	case AFBA:
		return AFBN
	case AFBE:
		return AFBNE
	case AFBNE:
		return AFBE
	case AFBG:
		return AFBLE
	case AFBLE:
		return AFBG
	case AFBGE:
		return AFBL
	case AFBL:
		return AFBGE
	}

	log.Fatalf("unknown relation: %s", obj.Aconv(a))
	return 0
}

var unaryDst = map[obj.As]bool{
	obj.ACALL: true,
	obj.AJMP:  true,
	AWORD:     true,
	ADWORD:    true,
	ABNW:      true,
	ABNEW:     true,
	ABEW:      true,
	ABGW:      true,
	ABLEW:     true,
	ABGEW:     true,
	ABLW:      true,
	ABGUW:     true,
	ABLEUW:    true,
	ABCCW:     true,
	ABCSW:     true,
	ABPOSW:    true,
	ABNEGW:    true,
	ABVCW:     true,
	ABVSW:     true,
	ABND:      true,
	ABNED:     true,
	ABED:      true,
	ABGD:      true,
	ABLED:     true,
	ABGED:     true,
	ABLD:      true,
	ABGUD:     true,
	ABLEUD:    true,
	ABCCD:     true,
	ABCSD:     true,
	ABPOSD:    true,
	ABNEGD:    true,
	ABVCD:     true,
	ABVSD:     true,
}

var Linksparc64 = obj.LinkArch{
	Arch:       sys.ArchSPARC64,
	Preprocess: preprocess,
	Assemble:   span,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
}

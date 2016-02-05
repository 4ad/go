// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

var isUncondJump = map[int16]bool{
	obj.ACALL:     true,
	obj.ADUFFZERO: true,
	obj.ADUFFCOPY: true,
	obj.AJMP:      true,
	obj.ARET:      true,
	AFBA:          true,
}

var isCondJump = map[int16]bool{
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

var isJump = make(map[int16]bool)

func init() {
	for k := range isUncondJump {
		isJump[k] = true
	}
	for k := range isCondJump {
		isJump[k] = true
	}
}

// Autoedit rewrites off(SP), off(FP), $off(SP), and $off(FP) to new(RFP).
func autoedit(a *obj.Addr) {
	if a == nil {
		return
	}
	if a.Type != obj.TYPE_MEM && a.Type != obj.TYPE_ADDR {
		return
	}
	if a.Name == obj.NAME_PARAM {
		a.Reg = REG_RFP
		a.Offset += MinStackFrameSize + StackBias
		a.Name = obj.TYPE_NONE
		return
	}
	if a.Name == obj.NAME_AUTO {
		a.Reg = REG_RFP
		a.Offset += StackBias
		a.Name = obj.TYPE_NONE
	}
}

func progedit(ctxt *obj.Link, p *obj.Prog) {
	// Rewite virtual SP/FP references to RFP. This has to happen
	// before any call to aclass, so it must be first here.
	autoedit(&p.From)
	autoedit(p.From3)
	autoedit(&p.To)

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
		if aclass(&p.From) == ClassConst {
			literal := fmt.Sprintf("$i64.%016x", p.From.Offset)
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
				cursym.Leaf = 1
			}
		}
	}

	for p := cursym.Text; p != nil; p = p.Link {
		switch p.As {
		case obj.ATEXT:
			// TODO(aram):
			// 	mdb(1) seems to work with this leaf prolog, but DTrace
			// 	doesn't. We need Dtrace, so we disable it for now.
			//
			// if cursym.Leaf == 1 {
			// 	frameSize := cursym.Locals
			// 	if frameSize == -8 {
			// 		frameSize = 0
			// }
			// if frameSize % 8 != 0 {
			// 	ctxt.Diag("%v: unaligned frame size %d - must be 0 mod 8", p, frameSize)
			// }

			// 	// MOVD RFP, (112+bias)(RSP)
			// 	p = obj.Appendp(ctxt, p)
			// 	p.As = AMOVD
			// 	p.From.Type = obj.TYPE_REG
			// 	p.From.Reg = REG_RFP
			// 	p.To.Type = obj.TYPE_MEM
			// 	p.To.Reg = REG_RSP
			// 	p.To.Offset = int64(112 + StackBias)

			// 	// ADD -(frame+128), RSP
			// 	p = obj.Appendp(ctxt, p)
			// 	p.As = AADD
			// 	p.From.Type = obj.TYPE_CONST
			// 	p.From.Offset = -int64(frameSize + WindowSaveAreaSize)
			// 	p.To.Type = obj.TYPE_REG
			// 	p.To.Reg = REG_RSP

			// 	// SUB -(frame+128), RSP, RFP
			// 	p = obj.Appendp(ctxt, p)
			// 	p.As = ASUB
			// 	p.From.Type = obj.TYPE_CONST
			// 	p.From.Offset = -int64(frameSize + WindowSaveAreaSize)
			// 	p.Reg = REG_RSP
			// 	p.To.Type = obj.TYPE_REG
			// 	p.To.Reg = REG_RFP

			// 	break
			// }

			frameSize := cursym.Locals
			if frameSize == -8 {
				frameSize = 0
			}
			if frameSize%8 != 0 {
				ctxt.Diag("%v: unaligned frame size %d - must be 0 mod 8", p, frameSize)
			}

			// MOVD RFP, (112+bias)(RSP)
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_RFP
			p.To.Type = obj.TYPE_MEM
			p.To.Reg = REG_RSP
			p.To.Offset = int64(112 + StackBias)

			// MOVD R31, (120+bias)(RSP)
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_R31
			p.To.Type = obj.TYPE_MEM
			p.To.Reg = REG_RSP
			p.To.Offset = int64(120 + StackBias)

			// ADD -(frame+128|176), RSP
			p = obj.Appendp(ctxt, p)
			p.As = AADD
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = -int64(frameSize + MinStackFrameSize)
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_RSP

			// SUB -(frame+128|176), RSP, RFP
			p = obj.Appendp(ctxt, p)
			p.As = ASUB
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = -int64(frameSize + MinStackFrameSize)
			p.Reg = REG_RSP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_RFP

			// MOVD LR, R31
			p = obj.Appendp(ctxt, p)
			p.As = AMOVD
			p.From.Type = obj.TYPE_REG
			p.From.Reg = REG_LR
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_R31

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

			// TODO(aram): arrange for the compiler to set the %i
			// registers and remove this code.

			// For all functions copy at most 6 arguments into the
			// %i registers.
			if args > 6 {
				args = 6
			}
			for i := 0; i < args; i++ {
				// MOVD	argN+8N(FP), %iN
				p = obj.Appendp(ctxt, p)
				p.As = AMOVD
				p.From.Type = obj.TYPE_MEM
				p.From.Reg = REG_RFP
				p.From.Offset = int64(i*8 + MinStackFrameSize + StackBias)
				p.To.Type = obj.TYPE_REG
				p.To.Reg = int16(REG_R24 + i) // %i0+i
			}

			// If the function is not a leaf, save incoming arguments
			// into the slot reserved for register window spill.
			if cursym.Leaf == 1 {
				break
			}
			for i := 0; i < 6; i++ {
				// MOVD %iN, (bias+64+8N)(RSP)
				p = obj.Appendp(ctxt, p)
				p.As = AMOVD
				p.From.Type = obj.TYPE_REG
				p.From.Reg = int16(REG_R24 + i) // %i0+i
				p.To.Type = obj.TYPE_MEM
				p.To.Reg = REG_RSP
				p.To.Offset = int64(i*8 + 64 + StackBias)
			}

		case obj.ARET:
			// TODO(aram):
			// 	mdb(1) seems to work with this leaf epilog, but DTrace
			// 	doesn't. We need Dtrace, so we disable it for now.
			//
			// if cursym.Leaf == 1 {
			// 	// MOVD RFP, TMP
			// 	q1 = p
			// 	p = obj.Appendp(ctxt, p)
			// 	p.As = obj.ARET
			// 	q1.As = AMOVD
			// 	q1.From.Type = obj.TYPE_REG
			// 	q1.From.Reg = REG_RFP
			// 	q1.To.Type = obj.TYPE_REG
			// 	q1.To.Reg = REG_TMP

			// 	// MOVD (112+StackBias)(RFP), RFP
			// 	q1 = obj.Appendp(ctxt, q1)
			// 	q1.As = AMOVD
			// 	q1.From.Type = obj.TYPE_MEM
			// 	q1.From.Reg = REG_RFP
			// 	q1.From.Offset = 112 + StackBias
			// 	q1.To.Type = obj.TYPE_REG
			// 	q1.To.Reg = REG_RFP

			// 	// MOVD TMP, RSP
			// 	q1 = obj.Appendp(ctxt, q1)
			// 	q1.As = AMOVD
			// 	q1.From.Type = obj.TYPE_REG
			// 	q1.From.Reg = REG_TMP
			// 	q1.To.Type = obj.TYPE_REG
			// 	q1.To.Reg = REG_RSP

			// 	break
			// }

			// MOVD RFP, TMP
			q1 = p
			p = obj.Appendp(ctxt, p)
			p.As = obj.ARET
			q1.As = AMOVD
			q1.From.Type = obj.TYPE_REG
			q1.From.Reg = REG_RFP
			q1.To.Type = obj.TYPE_REG
			q1.To.Reg = REG_TMP

			// MOVD R31, LR
			q1 = obj.Appendp(ctxt, q1)
			q1.As = AMOVD
			q1.From.Type = obj.TYPE_REG
			q1.From.Reg = REG_R31
			q1.To.Type = obj.TYPE_REG
			q1.To.Reg = REG_LR

			// MOVD (120+StackBias)(RFP), R31
			q1 = obj.Appendp(ctxt, q1)
			q1.As = AMOVD
			q1.From.Type = obj.TYPE_MEM
			q1.From.Reg = REG_RFP
			q1.From.Offset = 120 + StackBias
			q1.To.Type = obj.TYPE_REG
			q1.To.Reg = REG_R31

			// MOVD (112+StackBias)(RFP), RFP
			q1 = obj.Appendp(ctxt, q1)
			q1.As = AMOVD
			q1.From.Type = obj.TYPE_MEM
			q1.From.Reg = REG_RFP
			q1.From.Offset = 112 + StackBias
			q1.To.Type = obj.TYPE_REG
			q1.To.Reg = REG_RFP

			// MOVD TMP, RSP
			q1 = obj.Appendp(ctxt, q1)
			q1.As = AMOVD
			q1.From.Type = obj.TYPE_REG
			q1.From.Reg = REG_TMP
			q1.To.Type = obj.TYPE_REG
			q1.To.Reg = REG_RSP
		}
	}

	// Schedule delay-slots. Only RNOPs for now.
	for p := cursym.Text; p != nil; p = p.Link {
		if !isJump[p.As] {
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
		p.From.Class = aclass(&p.From)
		if p.From3 != nil {
			p.From3.Class = aclass(p.From3)
		}
		p.To.Class = aclass(&p.To)
	}
}

func relinv(a int) int {
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
		return AFBLG
	case AFBLG:
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

var unaryDst = map[int]bool{
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
	ByteOrder:  binary.BigEndian,
	Name:       "sparc64",
	Thechar:    'u',
	Preprocess: preprocess,
	Assemble:   span,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
	Minlc:      4,
	Ptrsize:    8,
	Regsize:    8,
}

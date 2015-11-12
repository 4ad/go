// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
	"fmt"
)

func defframe(ptxt *obj.Prog) {
	var n *gc.Node

	// fill in argument size, stack size
	ptxt.To.Type = obj.TYPE_TEXTSIZE

	ptxt.To.Val = int32(gc.Rnd(gc.Curfn.Type.Argwid, int64(gc.Widthptr)))
	frame := uint32(gc.Rnd(gc.Stksize+gc.Maxarg, int64(gc.Widthreg)))

	// sparc64 requires that the frame size (not counting saved LR)
	// be empty or be 8 mod 16. If not, pad it.
	if frame != 0 && frame%16 != 8 {
		frame += 8
	}

	ptxt.To.Offset = int64(frame)

	// insert code to zero ambiguously live variables
	// so that the garbage collector only sees initialized values
	// when it looks for pointers.
	p := ptxt

	hi := int64(0)
	lo := hi

	// iterate through declarations - they are sorted in decreasing xoffset order.
	for l := gc.Curfn.Func.Dcl; l != nil; l = l.Next {
		n = l.N
		if !n.Name.Needzero {
			continue
		}
		if n.Class != gc.PAUTO {
			gc.Fatalf("needzero class %d", n.Class)
		}
		if n.Type.Width%int64(gc.Widthptr) != 0 || n.Xoffset%int64(gc.Widthptr) != 0 || n.Type.Width == 0 {
			gc.Fatalf("var %v has size %d offset %d", gc.Nconv(n, obj.FmtLong), int(n.Type.Width), int(n.Xoffset))
		}

		if lo != hi && n.Xoffset+n.Type.Width >= lo-int64(2*gc.Widthreg) {
			// merge with range we already have
			lo = n.Xoffset

			continue
		}

		// zero old range
		p = zerorange(p, int64(frame), lo, hi)

		// set new range
		hi = n.Xoffset + n.Type.Width

		lo = n.Xoffset
	}

	// zero final range
	zerorange(p, int64(frame), lo, hi)
}

var darwin = obj.Getgoos() == "darwin"

func zerorange(p *obj.Prog, frame int64, lo int64, hi int64) *obj.Prog {
	cnt := hi - lo
	if cnt == 0 {
		return p
	}
	if cnt < int64(4*gc.Widthptr) {
		for i := int64(0); i < cnt; i += int64(gc.Widthptr) {
			p = appendpp(p, sparc64.AMOVD, obj.TYPE_REG, sparc64.REG_ZR, 0, obj.TYPE_MEM, sparc64.REG_RSP, 8+frame+lo+i)
		}
	} else if cnt <= int64(128*gc.Widthptr) && !darwin { // darwin ld64 cannot handle BR26 reloc with non-zero addend
		p = appendpp(p, sparc64.AMOVD, obj.TYPE_REG, sparc64.REG_RSP, 0, obj.TYPE_REG, sparc64.REG_RT1, 0)
		p = appendpp(p, sparc64.AADD, obj.TYPE_CONST, 0, 8+frame+lo-8, obj.TYPE_REG, sparc64.REG_RT1, 0)
		p.Reg = sparc64.REG_RT1
		p = appendpp(p, obj.ADUFFZERO, obj.TYPE_NONE, 0, 0, obj.TYPE_MEM, 0, 0)
		f := gc.Sysfunc("duffzero")
		gc.Naddr(&p.To, f)
		gc.Afunclit(&p.To, f)
		p.To.Offset = 4 * (128 - cnt/int64(gc.Widthptr))
	} else {
		p = appendpp(p, sparc64.AMOVD, obj.TYPE_CONST, 0, 8+frame+lo-8, obj.TYPE_REG, sparc64.REG_TMP, 0)
		p = appendpp(p, sparc64.AMOVD, obj.TYPE_REG, sparc64.REG_RSP, 0, obj.TYPE_REG, sparc64.REG_RT1, 0)
		p = appendpp(p, sparc64.AADD, obj.TYPE_REG, sparc64.REG_TMP, 0, obj.TYPE_REG, sparc64.REG_RT1, 0)
		p.Reg = sparc64.REG_RT1
		p = appendpp(p, sparc64.AMOVD, obj.TYPE_CONST, 0, cnt, obj.TYPE_REG, sparc64.REG_TMP, 0)
		p = appendpp(p, sparc64.AADD, obj.TYPE_REG, sparc64.REG_TMP, 0, obj.TYPE_REG, sparc64.REG_RT2, 0)
		p.Reg = sparc64.REG_RT1
		p = appendpp(p, sparc64.AMOVD, obj.TYPE_REG, sparc64.REG_ZR, 0, obj.TYPE_MEM, sparc64.REG_RT1, int64(gc.Widthptr))
		p.Scond = sparc64.C_XPRE
		p1 := p
		p = appendpp(p, sparc64.ACMP, obj.TYPE_REG, sparc64.REG_RT1, 0, obj.TYPE_NONE, 0, 0)
		p.Reg = sparc64.REG_RT2
		p = appendpp(p, sparc64.ABNE, obj.TYPE_NONE, 0, 0, obj.TYPE_BRANCH, 0, 0)
		gc.Patch(p, p1)
	}

	return p
}

func appendpp(p *obj.Prog, as int, ftype int, freg int, foffset int64, ttype int, treg int, toffset int64) *obj.Prog {
	q := gc.Ctxt.NewProg()
	gc.Clearp(q)
	q.As = int16(as)
	q.Lineno = p.Lineno
	q.From.Type = int16(ftype)
	q.From.Reg = int16(freg)
	q.From.Offset = foffset
	q.To.Type = int16(ttype)
	q.To.Reg = int16(treg)
	q.To.Offset = toffset
	q.Link = p.Link
	p.Link = q
	return q
}

func ginsnop() {
	var con gc.Node
	gc.Nodconst(&con, gc.Types[gc.TINT], 0)
	gins(sparc64.ARNOP, &con, nil)
}

var panicdiv *gc.Node

/*
 * generate division.
 * generates one of:
 *	res = nl / nr
 *	res = nl % nr
 * according to op.
 */
func dodiv(op int, nl *gc.Node, nr *gc.Node, res *gc.Node) {
	// Have to be careful about handling
	// most negative int divided by -1 correctly.
	// The hardware will generate undefined result.
	// Also need to explicitly trap on division on zero,
	// the hardware will silently generate undefined result.
	// DIVW will leave unpredicable result in higher 32-bit,
	// so always use DIVD/DIVDU.
	t := nl.Type

	t0 := t
	check := false
	if gc.Issigned[t.Etype] {
		check = true
		if gc.Isconst(nl, gc.CTINT) && nl.Int() != -(1<<uint64(t.Width*8-1)) {
			check = false
		} else if gc.Isconst(nr, gc.CTINT) && nr.Int() != -1 {
			check = false
		}
	}

	if t.Width < 8 {
		if gc.Issigned[t.Etype] {
			t = gc.Types[gc.TINT64]
		} else {
			t = gc.Types[gc.TUINT64]
		}
		check = false
	}

	a := optoas(gc.ODIV, t)

	var tl gc.Node
	gc.Regalloc(&tl, t0, nil)
	var tr gc.Node
	gc.Regalloc(&tr, t0, nil)
	if nl.Ullman >= nr.Ullman {
		gc.Cgen(nl, &tl)
		gc.Cgen(nr, &tr)
	} else {
		gc.Cgen(nr, &tr)
		gc.Cgen(nl, &tl)
	}

	if t != t0 {
		// Convert
		tl2 := tl

		tr2 := tr
		tl.Type = t
		tr.Type = t
		gmove(&tl2, &tl)
		gmove(&tr2, &tr)
	}

	// Handle divide-by-zero panic.
	p1 := gins(optoas(gc.OCMP, t), &tr, nil)
	p1.Reg = sparc64.REG_ZR
	p1 = gc.Gbranch(optoas(gc.ONE, t), nil, +1)
	if panicdiv == nil {
		panicdiv = gc.Sysfunc("panicdivide")
	}
	gc.Ginscall(panicdiv, -1)
	gc.Patch(p1, gc.Pc)

	var p2 *obj.Prog
	if check {
		var nm1 gc.Node
		gc.Nodconst(&nm1, t, -1)
		gcmp(optoas(gc.OCMP, t), &tr, &nm1)
		p1 := gc.Gbranch(optoas(gc.ONE, t), nil, +1)
		if op == gc.ODIV {
			// a / (-1) is -a.
			gins(optoas(gc.OMINUS, t), &tl, &tl)

			gmove(&tl, res)
		} else {
			// a % (-1) is 0.
			var nz gc.Node
			gc.Nodconst(&nz, t, 0)

			gmove(&nz, res)
		}

		p2 = gc.Gbranch(obj.AJMP, nil, 0)
		gc.Patch(p1, gc.Pc)
	}

	p1 = gins(a, &tr, &tl)
	if op == gc.ODIV {
		gc.Regfree(&tr)
		gmove(&tl, res)
	} else {
		// A%B = A-(A/B*B)
		var tm gc.Node
		gc.Regalloc(&tm, t, nil)

		// patch div to use the 3 register form
		// TODO(minux): add gins3?
		p1.Reg = p1.To.Reg

		p1.To.Reg = tm.Reg
		gins(optoas(gc.OMUL, t), &tr, &tm)
		gc.Regfree(&tr)
		gins(optoas(gc.OSUB, t), &tm, &tl)
		gc.Regfree(&tm)
		gmove(&tl, res)
	}

	gc.Regfree(&tl)
	if check {
		gc.Patch(p2, gc.Pc)
	}
}

/*
 * generate high multiply:
 *   res = (nl*nr) >> width
 */
func cgen_hmul(nl *gc.Node, nr *gc.Node, res *gc.Node) {
	// largest ullman on left.
	if nl.Ullman < nr.Ullman {
		nl, nr = nr, nl
	}

	t := (*gc.Type)(nl.Type)
	w := int(int(t.Width * 8))
	var n1 gc.Node
	gc.Cgenr(nl, &n1, res)
	var n2 gc.Node
	gc.Cgenr(nr, &n2, nil)
	switch gc.Simtype[t.Etype] {
	case gc.TINT8,
		gc.TINT16,
		gc.TINT32:
		gins(optoas(gc.OMUL, t), &n2, &n1)
		p := (*obj.Prog)(gins(sparc64.ASRAW, nil, &n1))
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(w)

	case gc.TUINT8,
		gc.TUINT16,
		gc.TUINT32:
		gins(optoas(gc.OMUL, t), &n2, &n1)
		p := (*obj.Prog)(gins(sparc64.ASRLW, nil, &n1))
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(w)

	// TODO(aram):
	// case gc.TINT64,
	// 	gc.TUINT64:
	// 	if gc.Issigned[t.Etype] {
	// 		gins(sparc64.ASMULH, &n2, &n1)
	// 	} else {
	// 		gins(sparc64.AUMULH, &n2, &n1)
	// 	}

	default:
		gc.Fatalf("cgen_hmul %v", t)
	}

	gc.Cgen(&n1, res)
	gc.Regfree(&n1)
	gc.Regfree(&n2)
}

/*
 * generate shift according to op, one of:
 *	res = nl << nr
 *	res = nl >> nr
 */
func cgen_shift(op int, bounded bool, nl *gc.Node, nr *gc.Node, res *gc.Node) {
	a := int(optoas(op, nl.Type))

	if nr.Op == gc.OLITERAL {
		var n1 gc.Node
		gc.Regalloc(&n1, nl.Type, res)
		gc.Cgen(nl, &n1)
		sc := uint64(nr.Int())
		if sc >= uint64(nl.Type.Width*8) {
			// large shift gets 2 shifts by width-1
			var n3 gc.Node
			gc.Nodconst(&n3, gc.Types[gc.TUINT32], nl.Type.Width*8-1)

			gins(a, &n3, &n1)
			gins(a, &n3, &n1)
		} else {
			gins(a, nr, &n1)
		}
		gmove(&n1, res)
		gc.Regfree(&n1)
		return
	}

	if nl.Ullman >= gc.UINF {
		var n4 gc.Node
		gc.Tempname(&n4, nl.Type)
		gc.Cgen(nl, &n4)
		nl = &n4
	}

	if nr.Ullman >= gc.UINF {
		var n5 gc.Node
		gc.Tempname(&n5, nr.Type)
		gc.Cgen(nr, &n5)
		nr = &n5
	}

	// Allow either uint32 or uint64 as shift type,
	// to avoid unnecessary conversion from uint32 to uint64
	// just to do the comparison.
	tcount := gc.Types[gc.Simtype[nr.Type.Etype]]

	if tcount.Etype < gc.TUINT32 {
		tcount = gc.Types[gc.TUINT32]
	}

	var n1 gc.Node
	gc.Regalloc(&n1, nr.Type, nil) // to hold the shift type in CX
	var n3 gc.Node
	gc.Regalloc(&n3, tcount, &n1) // to clear high bits of CX

	var n2 gc.Node
	gc.Regalloc(&n2, nl.Type, res)

	if nl.Ullman >= nr.Ullman {
		gc.Cgen(nl, &n2)
		gc.Cgen(nr, &n1)
		gmove(&n1, &n3)
	} else {
		gc.Cgen(nr, &n1)
		gmove(&n1, &n3)
		gc.Cgen(nl, &n2)
	}

	gc.Regfree(&n3)

	// test and fix up large shifts
	if !bounded {
		gc.Nodconst(&n3, tcount, nl.Type.Width*8)
		gcmp(optoas(gc.OCMP, tcount), &n1, &n3)
		p1 := (*obj.Prog)(gc.Gbranch(optoas(gc.OLT, tcount), nil, +1))
		if op == gc.ORSH && gc.Issigned[nl.Type.Etype] {
			gc.Nodconst(&n3, gc.Types[gc.TUINT32], nl.Type.Width*8-1)
			gins(a, &n3, &n2)
		} else {
			gc.Nodconst(&n3, nl.Type, 0)
			gmove(&n3, &n2)
		}

		gc.Patch(p1, gc.Pc)
	}

	gins(a, &n1, &n2)

	gmove(&n2, res)

	gc.Regfree(&n1)
	gc.Regfree(&n2)
}

func clearfat(nl *gc.Node) {
	/* clear a fat object */
	if gc.Debug['g'] != 0 {
		fmt.Printf("clearfat %v (%v, size: %d)\n", nl, nl.Type, nl.Type.Width)
	}

	w := uint64(uint64(nl.Type.Width))

	// Avoid taking the address for simple enough types.
	if gc.Componentgen(nil, nl) {
		return
	}

	c := uint64(w % 8) // bytes
	q := uint64(w / 8) // dwords

	var r0 gc.Node
	gc.Nodreg(&r0, gc.Types[gc.TUINT64], sparc64.REG_ZR)
	var dst gc.Node

	// REGRT1 is reserved on sparc64, see sparc64/gsubr.go.
	gc.Nodreg(&dst, gc.Types[gc.Tptr], sparc64.REG_RT1)
	gc.Agen(nl, &dst)

	var boff uint64
	if q > 128 {
		p := gins(sparc64.ASUB, nil, &dst)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 8

		var end gc.Node
		gc.Regalloc(&end, gc.Types[gc.Tptr], nil)
		p = gins(sparc64.AMOVD, &dst, &end)
		p.From.Type = obj.TYPE_ADDR
		p.From.Offset = int64(q * 8)

		p = gins(sparc64.AMOVD, &r0, &dst)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = 8
		p.Scond = sparc64.C_XPRE
		pl := (*obj.Prog)(p)

		p = gcmp(sparc64.ACMP, &dst, &end)
		gc.Patch(gc.Gbranch(sparc64.ABNE, nil, 0), pl)

		gc.Regfree(&end)

		// The loop leaves R16 on the last zeroed dword
		boff = 8
	} else if q >= 4 && !darwin { // darwin ld64 cannot handle BR26 reloc with non-zero addend
		p := gins(sparc64.ASUB, nil, &dst)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = 8
		f := (*gc.Node)(gc.Sysfunc("duffzero"))
		p = gins(obj.ADUFFZERO, nil, f)
		gc.Afunclit(&p.To, f)

		// 4 and 128 = magic constants: see ../../runtime/asm_sparc64x.s
		p.To.Offset = int64(4 * (128 - q))

		// duffzero leaves R16 on the last zeroed dword
		boff = 8
	} else {
		var p *obj.Prog
		for t := uint64(0); t < q; t++ {
			p = gins(sparc64.AMOVD, &r0, &dst)
			p.To.Type = obj.TYPE_MEM
			p.To.Offset = int64(8 * t)
		}

		boff = 8 * q
	}

	var p *obj.Prog
	for t := uint64(0); t < c; t++ {
		p = gins(sparc64.AMOVB, &r0, &dst)
		p.To.Type = obj.TYPE_MEM
		p.To.Offset = int64(t + boff)
	}
}

// Called after regopt and peep have run.
// Expand CHECKNIL pseudo-op into actual nil pointer check.
func expandchecks(firstp *obj.Prog) {
	var p1 *obj.Prog

	for p := (*obj.Prog)(firstp); p != nil; p = p.Link {
		if gc.Debug_checknil != 0 && gc.Ctxt.Debugvlog != 0 {
			fmt.Printf("expandchecks: %v\n", p)
		}
		if p.As != obj.ACHECKNIL {
			continue
		}
		if gc.Debug_checknil != 0 && p.Lineno > 1 { // p->lineno==1 in generated wrappers
			gc.Warnl(int(p.Lineno), "generated nil check")
		}
		if p.From.Type != obj.TYPE_REG {
			gc.Fatalf("invalid nil check %v\n", p)
		}

		// check is
		//	BRNZ arg, 2(PC)
		//	MOVD ZR, 0(arg)
		p1 = gc.Ctxt.NewProg()
		gc.Clearp(p1)
		p1.Link = p.Link
		p.Link = p1
		p1.Lineno = p.Lineno
		p1.Pc = 9999

		p.As = sparc64.ABRNZ
		p.To.Type = obj.TYPE_BRANCH
		p.To.Val = p1.Link

		// crash by write to memory address 0.
		p1.As = sparc64.AMOVD
		p1.From.Type = obj.TYPE_REG
		p1.From.Reg = sparc64.REG_ZR
		p1.To.Type = obj.TYPE_MEM
		p1.To.Reg = p.From.Reg
		p1.To.Offset = 0
	}
}

// res = runtime.getg()
func getg(res *gc.Node) {
	var n1 gc.Node
	gc.Nodreg(&n1, res.Type, sparc64.REG_G)
	gmove(&n1, res)
}

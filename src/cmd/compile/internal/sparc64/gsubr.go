// Derived from Inferno utils/6c/txt.c
// http://code.google.com/p/inferno-os/source/browse/utils/6c/txt.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sparc64

import (
	"cmd/compile/internal/big"
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
	"fmt"
)

var resvd = []int{
	sparc64.REG_ZR,
	sparc64.REG_RT1,
	sparc64.REG_CTXT,
	sparc64.REG_G,
	sparc64.REG_RT2,
	sparc64.REG_TMP,
	sparc64.REG_G6,
	sparc64.REG_TLS,
	sparc64.REG_RSP,
	sparc64.REG_OLR,
	sparc64.REG_TMP2,
	sparc64.REG_L7,
	sparc64.REG_I0, // TODO(aram): revisit this.
	sparc64.REG_I1, // TODO(aram): revisit this.
	sparc64.REG_I2, // TODO(aram): revisit this.
	sparc64.REG_I3, // TODO(aram): revisit this.
	sparc64.REG_I4, // TODO(aram): revisit this.
	sparc64.REG_I5, // TODO(aram): revisit this.
	sparc64.REG_RFP,
	sparc64.REG_ILR,
	sparc64.REG_YTMP,
	sparc64.REG_YTWO,
}

/*
 * generate
 *	as $c, n
 */
func ginscon(as obj.As, c int64, n2 *gc.Node) {
	var n1 gc.Node

	gc.Nodconst(&n1, gc.Types[gc.TINT64], c)

	if as != sparc64.AMOVD && (c < -sparc64.BIG || c > sparc64.BIG) || as == sparc64.AMULD || n2 != nil && n2.Op != gc.OREGISTER {
		// cannot have more than 13-bit of immediate in ADD, etc.
		// instead, MOV into register first.
		var ntmp gc.Node
		gc.Regalloc(&ntmp, gc.Types[gc.TINT64], nil)

		gins(sparc64.AMOVD, &n1, &ntmp)
		gins(as, &ntmp, n2)
		gc.Regfree(&ntmp)
		return
	}

	rawgins(as, &n1, n2)
}

/*
 * generate
 *	as n, $c (CMP)
 */
func ginscon2(as obj.As, n2 *gc.Node, c int64) {
	var n1 gc.Node

	gc.Nodconst(&n1, gc.Types[gc.TINT64], c)

	switch as {
	default:
		gc.Fatalf("ginscon2")

	case sparc64.ACMP:
		if -sparc64.BIG <= c && c <= sparc64.BIG {
			gcmp(as, n2, &n1)
			return
		}
	}

	// MOV n1 into register first
	var ntmp gc.Node
	gc.Regalloc(&ntmp, gc.Types[gc.TINT64], nil)

	rawgins(sparc64.AMOVD, &n1, &ntmp)
	gcmp(as, n2, &ntmp)
	gc.Regfree(&ntmp)
}

func ginscmp(op gc.Op, t *gc.Type, n1, n2 *gc.Node, likely int) *obj.Prog {
	if gc.Isint[t.Etype] && n1.Op == gc.OLITERAL && n2.Op != gc.OLITERAL {
		// Reverse comparison to place constant last.
		op = gc.Brrev(op)
		n1, n2 = n2, n1
	}

	var r1, r2, g1, g2 gc.Node
	gc.Regalloc(&r1, t, n1)
	gc.Regalloc(&g1, n1.Type, &r1)
	gc.Cgen(n1, &g1)
	gmove(&g1, &r1)
	if gc.Isint[t.Etype] && gc.Isconst(n2, gc.CTINT) {
		ginscon2(optoas(gc.OCMP, t), &r1, n2.Int64())
	} else {
		gc.Regalloc(&r2, t, n2)
		gc.Regalloc(&g2, n1.Type, &r2)
		gc.Cgen(n2, &g2)
		gmove(&g2, &r2)
		gcmp(optoas(gc.OCMP, t), &r1, &r2)
		gc.Regfree(&g2)
		gc.Regfree(&r2)
	}
	gc.Regfree(&g1)
	gc.Regfree(&r1)
	return gc.Gbranch(optoas(op, t), nil, likely)
}

// set up nodes representing 2^63
var (
	bigi         gc.Node
	bigf         gc.Node
	bignodes_did bool
)

func bignodes() {
	if bignodes_did {
		return
	}
	bignodes_did = true

	var i big.Int
	i.SetInt64(1)
	i.Lsh(&i, 63)

	gc.Nodconst(&bigi, gc.Types[gc.TUINT64], 0)
	bigi.SetBigInt(&i)

	bigi.Convconst(&bigf, gc.Types[gc.TFLOAT64])
}

/*
 * generate move:
 *	t = f
 * hard part is conversions.
 */
func gmove(f *gc.Node, t *gc.Node) {
	if gc.Debug['M'] != 0 {
		fmt.Printf("gmove %v -> %v\n", gc.Nconv(f, gc.FmtLong), gc.Nconv(t, gc.FmtLong))
	}

	ft := int(gc.Simsimtype(f.Type))
	tt := int(gc.Simsimtype(t.Type))
	cvt := t.Type

	if gc.Iscomplex[ft] || gc.Iscomplex[tt] {
		gc.Complexmove(f, t)
		return
	}

	// cannot have two memory operands
	var r1 gc.Node
	var r2 gc.Node
	var a obj.As
	if gc.Ismem(f) && gc.Ismem(t) {
		goto hard
	}

	// convert constant to desired type
	if f.Op == gc.OLITERAL {
		var con gc.Node
		switch tt {
		default:
			f.Convconst(&con, t.Type)

		case gc.TINT32,
			gc.TINT16,
			gc.TINT8:
			var con gc.Node
			f.Convconst(&con, gc.Types[gc.TINT64])
			var r1 gc.Node
			gc.Regalloc(&r1, con.Type, t)
			gins(sparc64.AMOVD, &con, &r1)
			gmove(&r1, t)
			gc.Regfree(&r1)
			return

		case gc.TUINT32,
			gc.TUINT16,
			gc.TUINT8:
			var con gc.Node
			f.Convconst(&con, gc.Types[gc.TUINT64])
			var r1 gc.Node
			gc.Regalloc(&r1, con.Type, t)
			gins(sparc64.AMOVD, &con, &r1)
			gmove(&r1, t)
			gc.Regfree(&r1)
			return
		}

		f = &con
		ft = tt // so big switch will choose a simple mov

		// constants can't move directly to memory.
		if gc.Ismem(t) {
			goto hard
		}
	}

	// value -> value copy, first operand in memory.
	// any floating point operand requires register
	// src, so goto hard to copy to register first.
	if gc.Ismem(f) && ft != tt && (gc.Isfloat[ft] || gc.Isfloat[tt]) {
		cvt = gc.Types[ft]
		goto hard
	}

	// value -> value copy, only one memory operand.
	// figure out the instruction to use.
	// break out of switch for one-instruction gins.
	// goto fsrccpy for "float operation, and source is an integer register".
	// goto fdstcpy for "float operation, and destination is an integer register".
	// goto rdst for "destination must be register".
	// goto hard for "convert to cvt type first".
	// otherwise handle and return.

	switch uint32(ft)<<16 | uint32(tt) {
	default:
		gc.Fatalf("gmove %v -> %v", gc.Tconv(f.Type, gc.FmtLong), gc.Tconv(t.Type, gc.FmtLong))

		/*
		 * integer copy and truncate
		 */
	case gc.TINT8<<16 | gc.TINT8, // same size
		gc.TUINT8<<16 | gc.TINT8,
		gc.TINT16<<16 | gc.TINT8,
		// truncate
		gc.TUINT16<<16 | gc.TINT8,
		gc.TINT32<<16 | gc.TINT8,
		gc.TUINT32<<16 | gc.TINT8,
		gc.TINT64<<16 | gc.TINT8,
		gc.TUINT64<<16 | gc.TINT8:
		a = sparc64.AMOVB

	case gc.TINT8<<16 | gc.TUINT8, // same size
		gc.TUINT8<<16 | gc.TUINT8,
		gc.TINT16<<16 | gc.TUINT8,
		// truncate
		gc.TUINT16<<16 | gc.TUINT8,
		gc.TINT32<<16 | gc.TUINT8,
		gc.TUINT32<<16 | gc.TUINT8,
		gc.TINT64<<16 | gc.TUINT8,
		gc.TUINT64<<16 | gc.TUINT8:
		a = sparc64.AMOVUB

	case gc.TINT16<<16 | gc.TINT16, // same size
		gc.TUINT16<<16 | gc.TINT16,
		gc.TINT32<<16 | gc.TINT16,
		// truncate
		gc.TUINT32<<16 | gc.TINT16,
		gc.TINT64<<16 | gc.TINT16,
		gc.TUINT64<<16 | gc.TINT16:
		a = sparc64.AMOVH

	case gc.TINT16<<16 | gc.TUINT16, // same size
		gc.TUINT16<<16 | gc.TUINT16,
		gc.TINT32<<16 | gc.TUINT16,
		// truncate
		gc.TUINT32<<16 | gc.TUINT16,
		gc.TINT64<<16 | gc.TUINT16,
		gc.TUINT64<<16 | gc.TUINT16:
		a = sparc64.AMOVUH

	case gc.TINT32<<16 | gc.TINT32, // same size
		gc.TUINT32<<16 | gc.TINT32,
		gc.TINT64<<16 | gc.TINT32,
		// truncate
		gc.TUINT64<<16 | gc.TINT32:
		a = sparc64.AMOVW

	case gc.TINT32<<16 | gc.TUINT32, // same size
		gc.TUINT32<<16 | gc.TUINT32,
		gc.TINT64<<16 | gc.TUINT32,
		gc.TUINT64<<16 | gc.TUINT32:
		a = sparc64.AMOVUW

	case gc.TINT64<<16 | gc.TINT64, // same size
		gc.TINT64<<16 | gc.TUINT64,
		gc.TUINT64<<16 | gc.TINT64,
		gc.TUINT64<<16 | gc.TUINT64:
		a = sparc64.AMOVD

		/*
		 * integer up-conversions
		 */
	case gc.TINT8<<16 | gc.TINT16, // sign extend int8
		gc.TINT8<<16 | gc.TUINT16,
		gc.TINT8<<16 | gc.TINT32,
		gc.TINT8<<16 | gc.TUINT32,
		gc.TINT8<<16 | gc.TINT64,
		gc.TINT8<<16 | gc.TUINT64:
		a = sparc64.AMOVB

		goto rdst

	case gc.TUINT8<<16 | gc.TINT16, // zero extend uint8
		gc.TUINT8<<16 | gc.TUINT16,
		gc.TUINT8<<16 | gc.TINT32,
		gc.TUINT8<<16 | gc.TUINT32,
		gc.TUINT8<<16 | gc.TINT64,
		gc.TUINT8<<16 | gc.TUINT64:
		a = sparc64.AMOVUB

		goto rdst

	case gc.TINT16<<16 | gc.TINT32, // sign extend int16
		gc.TINT16<<16 | gc.TUINT32,
		gc.TINT16<<16 | gc.TINT64,
		gc.TINT16<<16 | gc.TUINT64:
		a = sparc64.AMOVH

		goto rdst

	case gc.TUINT16<<16 | gc.TINT32, // zero extend uint16
		gc.TUINT16<<16 | gc.TUINT32,
		gc.TUINT16<<16 | gc.TINT64,
		gc.TUINT16<<16 | gc.TUINT64:
		a = sparc64.AMOVUH

		goto rdst

	case gc.TINT32<<16 | gc.TINT64, // sign extend int32
		gc.TINT32<<16 | gc.TUINT64:
		a = sparc64.AMOVW

		goto rdst

	case gc.TUINT32<<16 | gc.TINT64, // zero extend uint32
		gc.TUINT32<<16 | gc.TUINT64:
		a = sparc64.AMOVUW

		goto rdst

	//return;
	// algorithm is:
	//	if small enough, use native float64 -> int64 conversion.
	//	otherwise, subtract 2^63, convert, and add it back.
	/*
	* float to integer
	 */
	case gc.TFLOAT32<<16 | gc.TINT32,
		gc.TFLOAT32<<16 | gc.TINT64,
		gc.TFLOAT32<<16 | gc.TINT16,
		gc.TFLOAT32<<16 | gc.TINT8,
		gc.TFLOAT32<<16 | gc.TUINT16,
		gc.TFLOAT32<<16 | gc.TUINT8,
		gc.TFLOAT32<<16 | gc.TUINT32,
		gc.TFLOAT32<<16 | gc.TUINT64:
		cvt = gc.Types[gc.TFLOAT64]

		goto hard

	case gc.TFLOAT64<<16 | gc.TINT32,
		gc.TFLOAT64<<16 | gc.TINT64,
		gc.TFLOAT64<<16 | gc.TINT16,
		gc.TFLOAT64<<16 | gc.TINT8,
		gc.TFLOAT64<<16 | gc.TUINT16,
		gc.TFLOAT64<<16 | gc.TUINT8,
		gc.TFLOAT64<<16 | gc.TUINT32,
		gc.TFLOAT64<<16 | gc.TUINT64:
		bignodes()

		var r1 gc.Node
		gc.Regalloc(&r1, gc.Types[ft], f)
		gmove(f, &r1)
		if tt == gc.TUINT64 {
			gc.Regalloc(&r2, gc.Types[gc.TFLOAT64], nil)
			gmove(&bigf, &r2)
			gins(sparc64.AFCMPD, &r2, &r1)
			p1 := gc.Gbranch(optoas(gc.OGT, gc.Types[gc.TFLOAT64]), nil, +1)
			gins(sparc64.AFSUBD, &r2, &r1)
			gc.Patch(p1, gc.Pc)
			gc.Regfree(&r2)
		}

		gc.Regalloc(&r2, gc.Types[gc.TFLOAT64], nil)
		var r3 gc.Node
		gc.Regalloc(&r3, gc.Types[gc.TINT64], t)
		gins(sparc64.AFDTOX, &r1, &r2)
		p1 := gins(sparc64.AFMOVD, &r2, nil)
		p1.To.Type = obj.TYPE_MEM
		p1.To.Reg = sparc64.REG_RSP
		p1.To.Offset = -8 + sparc64.StackBias
		p1 = gins(sparc64.AMOVD, nil, &r3)
		p1.From.Type = obj.TYPE_MEM
		p1.From.Reg = sparc64.REG_RSP
		p1.From.Offset = -8 + sparc64.StackBias
		gc.Regfree(&r2)
		gc.Regfree(&r1)
		if tt == gc.TUINT64 {
			p1 := gc.Gbranch(optoas(gc.OGT, gc.Types[gc.TFLOAT64]), nil, +1)
			gc.Nodreg(&r1, gc.Types[gc.TINT64], sparc64.REG_RT1)
			gins(sparc64.AMOVD, &bigi, &r1)
			gins(sparc64.AADD, &r1, &r3)
			gc.Patch(p1, gc.Pc)
		}

		gmove(&r3, t)
		gc.Regfree(&r3)
		return

		//warn("gmove: convert int to float not implemented: %N -> %N\n", f, t);
	//return;
	// algorithm is:
	//	if small enough, use native int64 -> uint64 conversion.
	//	otherwise, halve (rounding to odd?), convert, and double.
	/*
	 * integer to float
	 */
	case gc.TINT32<<16 | gc.TFLOAT32,
		gc.TINT32<<16 | gc.TFLOAT64,
		gc.TINT16<<16 | gc.TFLOAT32,
		gc.TINT16<<16 | gc.TFLOAT64,
		gc.TINT8<<16 | gc.TFLOAT32,
		gc.TINT8<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT64]

		goto hard

	case gc.TUINT16<<16 | gc.TFLOAT32,
		gc.TUINT16<<16 | gc.TFLOAT64,
		gc.TUINT8<<16 | gc.TFLOAT32,
		gc.TUINT8<<16 | gc.TFLOAT64,
		gc.TUINT32<<16 | gc.TFLOAT32,
		gc.TUINT32<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TUINT64]

		goto hard

	case gc.TINT64<<16 | gc.TFLOAT32,
		gc.TINT64<<16 | gc.TFLOAT64,
		gc.TUINT64<<16 | gc.TFLOAT32,
		gc.TUINT64<<16 | gc.TFLOAT64:
		bignodes()

		// The algorithm is:
		//	if small enough, use native int64 -> float64 conversion,
		//	otherwise halve (x -> (x>>1)|(x&1)), convert, and double.
		var r1 gc.Node
		gc.Regalloc(&r1, gc.Types[gc.TINT64], nil)
		gmove(f, &r1)
		if ft == gc.TUINT64 {
			gc.Nodreg(&r2, gc.Types[gc.TUINT64], sparc64.REG_RT1)
			gmove(&bigi, &r2)
			gins(sparc64.ACMP, &r1, &r2)
			p1 := gc.Gbranch(sparc64.ABLEUD, nil, +1)
			var r3 gc.Node
			gc.Regalloc(&r3, gc.Types[gc.TUINT64], nil)
			p2 := gins(sparc64.AAND, nil, &r3) // andi.
			p2.Reg = r1.Reg
			p2.From.Type = obj.TYPE_CONST
			p2.From.Offset = 1
			p3 := gins(sparc64.ASRLD, nil, &r1)
			p3.From.Type = obj.TYPE_CONST
			p3.From.Offset = 1
			gins(sparc64.AOR, &r3, &r1)
			gc.Regfree(&r3)
			gc.Patch(p1, gc.Pc)
		}

		gc.Regalloc(&r2, gc.Types[gc.TFLOAT64], t)
		p1 := gins(sparc64.AMOVD, &r1, nil)
		p1.To.Type = obj.TYPE_MEM
		p1.To.Reg = sparc64.REG_RSP
		p1.To.Offset = -8 + sparc64.StackBias
		p1 = gins(sparc64.AFMOVD, nil, &r2)
		p1.From.Type = obj.TYPE_MEM
		p1.From.Reg = sparc64.REG_RSP
		p1.From.Offset = -8 + sparc64.StackBias
		gins(sparc64.AFXTOD, &r2, &r2)
		gc.Regfree(&r1)
		if ft == gc.TUINT64 {
			p1 := gc.Gbranch(sparc64.ABLEUD, nil, +1)
			gc.Nodreg(&r1, gc.Types[gc.TFLOAT64], sparc64.REG_YTWO)
			gins(sparc64.AFMULD, &r1, &r2)
			gc.Patch(p1, gc.Pc)
		}

		gmove(&r2, t)
		gc.Regfree(&r2)
		return

	/*
	 * float to float
	 */
	case gc.TFLOAT32<<16 | gc.TFLOAT32:
		a = sparc64.AFMOVS

	case gc.TFLOAT64<<16 | gc.TFLOAT64:
		a = sparc64.AFMOVD

	case gc.TFLOAT32<<16 | gc.TFLOAT64:
		a = sparc64.AFSTOD
		goto rdst

	case gc.TFLOAT64<<16 | gc.TFLOAT32:
		a = sparc64.AFDTOS
		goto rdst
	}

	gins(a, f, t)
	return

	// requires register destination
rdst:
	gc.Regalloc(&r1, t.Type, t)

	gins(a, f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return

	// requires register intermediate
hard:
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return
}

// gins is called by the front end.
// It synthesizes some multiple-instruction sequences
// so the front end can stay simpler.
func gins(as obj.As, f, t *gc.Node) *obj.Prog {
	if as >= obj.A_ARCHSPECIFIC {
		if x, ok := f.IntLiteral(); ok {
			ginscon(as, x, t)
			return nil // caller must not use
		}
	}
	if as == sparc64.ACMP {
		if x, ok := t.IntLiteral(); ok {
			ginscon2(as, f, x)
			return nil // caller must not use
		}
	}
	return rawgins(as, f, t)
}

/*
 * generate one instruction:
 *	as f, t
 */
func rawgins(as obj.As, f *gc.Node, t *gc.Node) *obj.Prog {
	// TODO(austin): Add self-move test like in 6g (but be careful
	// of truncation moves)

	p := gc.Prog(as)
	gc.Naddr(&p.From, f)
	gc.Naddr(&p.To, t)

	switch as {
	case sparc64.ACMP, sparc64.AFCMPS, sparc64.AFCMPD:
		if t != nil {
			if f.Op != gc.OREGISTER {
				gc.Fatalf("bad operands to gcmp")
			}
			p.From = p.To
			p.To = obj.Addr{}
			raddr(f, p)
		}
	}

	// Bad things the front end has done to us. Crash to find call stack.
	switch as {
	case sparc64.AAND, sparc64.AMULD:
		if p.From.Type == obj.TYPE_CONST {
			gc.Debug['h'] = 1
			gc.Fatalf("bad inst: %v", p)
		}
	case sparc64.ACMP:
		if p.From.Type == obj.TYPE_MEM || p.To.Type == obj.TYPE_MEM {
			gc.Debug['h'] = 1
			gc.Fatalf("bad inst: %v", p)
		}
	}

	if gc.Debug['g'] != 0 {
		fmt.Printf("%v\n", p)
	}

	w := int32(0)
	switch as {
	case sparc64.AMOVB,
		sparc64.AMOVUB:
		w = 1

	case sparc64.AMOVH,
		sparc64.AMOVUH:
		w = 2

	case sparc64.AMOVW,
		sparc64.AMOVUW:
		w = 4

	case sparc64.AMOVD:
		if p.From.Type == obj.TYPE_CONST || p.From.Type == obj.TYPE_ADDR {
			break
		}
		w = 8
	}

	if w != 0 && ((f != nil && p.From.Width < int64(w)) || (t != nil && p.To.Type != obj.TYPE_REG && p.To.Width > int64(w))) {
		gc.Dump("f", f)
		gc.Dump("t", t)
		gc.Fatalf("bad width: %v (%d, %d)\n", p, p.From.Width, p.To.Width)
	}

	return p
}

/*
 * insert n into reg slot of p
 */
func raddr(n *gc.Node, p *obj.Prog) {
	var a obj.Addr

	gc.Naddr(&a, n)
	if a.Type != obj.TYPE_REG {
		if n != nil {
			gc.Fatalf("bad in raddr: %v", n.Op)
		} else {
			gc.Fatalf("bad in raddr: <null>")
		}
		p.Reg = 0
	} else {
		p.Reg = a.Reg
	}
}

func gcmp(as obj.As, lhs *gc.Node, rhs *gc.Node) *obj.Prog {
	if lhs.Op != gc.OREGISTER {
		gc.Fatalf("bad operands to gcmp: %v %v", lhs.Op, rhs.Op)
	}

	p := rawgins(as, rhs, nil)
	raddr(lhs, p)
	return p
}

/*
 * return Axxx for Oxxx on type t.
 */
func optoas(op gc.Op, t *gc.Type) obj.As {
	if t == nil {
		gc.Fatalf("optoas: t is nil")
	}

	// avoid constant conversions in switches below
	const (
		OMINUS_ = uint32(gc.OMINUS) << 16
		OLSH_   = uint32(gc.OLSH) << 16
		ORSH_   = uint32(gc.ORSH) << 16
		OADD_   = uint32(gc.OADD) << 16
		OSUB_   = uint32(gc.OSUB) << 16
		OMUL_   = uint32(gc.OMUL) << 16
		ODIV_   = uint32(gc.ODIV) << 16
		OOR_    = uint32(gc.OOR) << 16
		OAND_   = uint32(gc.OAND) << 16
		OXOR_   = uint32(gc.OXOR) << 16
		OEQ_    = uint32(gc.OEQ) << 16
		ONE_    = uint32(gc.ONE) << 16
		OLT_    = uint32(gc.OLT) << 16
		OLE_    = uint32(gc.OLE) << 16
		OGE_    = uint32(gc.OGE) << 16
		OGT_    = uint32(gc.OGT) << 16
		OCMP_   = uint32(gc.OCMP) << 16
		OAS_    = uint32(gc.OAS) << 16
		OHMUL_  = uint32(gc.OHMUL) << 16
		OSQRT_  = uint32(gc.OSQRT) << 16
	)

	a := obj.AXXX
	switch uint32(op)<<16 | uint32(gc.Simtype[t.Etype]) {
	default:
		gc.Fatalf("optoas: no entry for op=%v type=%v", op, t)

	case OEQ_ | gc.TBOOL,
		OEQ_ | gc.TINT8,
		OEQ_ | gc.TUINT8,
		OEQ_ | gc.TINT16,
		OEQ_ | gc.TUINT16,
		OEQ_ | gc.TINT32,
		OEQ_ | gc.TUINT32,
		OEQ_ | gc.TPTR32:
		a = sparc64.ABEW

	case OEQ_ | gc.TINT64,
		OEQ_ | gc.TUINT64,
		OEQ_ | gc.TPTR64:
		a = sparc64.ABED

	case OEQ_ | gc.TFLOAT32,
		OEQ_ | gc.TFLOAT64:
		a = sparc64.AFBE

	case ONE_ | gc.TBOOL,
		ONE_ | gc.TINT8,
		ONE_ | gc.TUINT8,
		ONE_ | gc.TINT16,
		ONE_ | gc.TUINT16,
		ONE_ | gc.TINT32,
		ONE_ | gc.TUINT32,
		ONE_ | gc.TPTR32:
		a = sparc64.ABNEW

	case ONE_ | gc.TINT64,
		ONE_ | gc.TUINT64,
		ONE_ | gc.TPTR64:
		a = sparc64.ABNED

	case ONE_ | gc.TFLOAT32,
		ONE_ | gc.TFLOAT64:
		a = sparc64.AFBNE

	case OLT_ | gc.TINT8,
		OLT_ | gc.TINT16,
		OLT_ | gc.TINT32:
		a = sparc64.ABLW

	case OLT_ | gc.TINT64:
		a = sparc64.ABLD

	case OLT_ | gc.TUINT8,
		OLT_ | gc.TUINT16,
		OLT_ | gc.TUINT32:
		a = sparc64.ABCSW

	case OLT_ | gc.TUINT64:
		a = sparc64.ABCSD

	case OLT_ | gc.TFLOAT32,
		OLT_ | gc.TFLOAT64:
		a = sparc64.AFBL

	case OLE_ | gc.TINT8,
		OLE_ | gc.TINT16,
		OLE_ | gc.TINT32:
		a = sparc64.ABLEW

	case OLE_ | gc.TINT64:
		a = sparc64.ABLED

	case OLE_ | gc.TUINT8,
		OLE_ | gc.TUINT16,
		OLE_ | gc.TUINT32:
		a = sparc64.ABLEUW

	case OLE_ | gc.TUINT64:
		a = sparc64.ABLEUD

	case OLE_ | gc.TFLOAT32,
		OLE_ | gc.TFLOAT64:
		a = sparc64.AFBLE

	case OGT_ | gc.TINT8,
		OGT_ | gc.TINT16,
		OGT_ | gc.TINT32:
		a = sparc64.ABGW

	case OGT_ | gc.TINT64:
		a = sparc64.ABGD

	case OGT_ | gc.TFLOAT32,
		OGT_ | gc.TFLOAT64:
		a = sparc64.AFBG

	case OGT_ | gc.TUINT8,
		OGT_ | gc.TUINT16,
		OGT_ | gc.TUINT32:
		a = sparc64.ABGUW

	case OGT_ | gc.TUINT64:
		a = sparc64.ABGUD

	case OGE_ | gc.TINT8,
		OGE_ | gc.TINT16,
		OGE_ | gc.TINT32:
		a = sparc64.ABGEW

	case OGE_ | gc.TINT64:
		a = sparc64.ABGED

	case OGE_ | gc.TFLOAT32,
		OGE_ | gc.TFLOAT64:
		a = sparc64.AFBGE

	case OGE_ | gc.TUINT8,
		OGE_ | gc.TUINT16,
		OGE_ | gc.TUINT32:
		a = sparc64.ABCCW

	case OGE_ | gc.TUINT64:
		a = sparc64.ABCCD

	case OCMP_ | gc.TBOOL,
		OCMP_ | gc.TINT8,
		OCMP_ | gc.TINT16,
		OCMP_ | gc.TINT32,
		OCMP_ | gc.TPTR32,
		OCMP_ | gc.TINT64,
		OCMP_ | gc.TUINT8,
		OCMP_ | gc.TUINT16,
		OCMP_ | gc.TUINT32,
		OCMP_ | gc.TUINT64,
		OCMP_ | gc.TPTR64:
		a = sparc64.ACMP

	case OCMP_ | gc.TFLOAT32:
		a = sparc64.AFCMPS

	case OCMP_ | gc.TFLOAT64:
		a = sparc64.AFCMPD

	case OAS_ | gc.TBOOL,
		OAS_ | gc.TINT8:
		a = sparc64.AMOVB

	case OAS_ | gc.TUINT8:
		a = sparc64.AMOVUB

	case OAS_ | gc.TINT16:
		a = sparc64.AMOVH

	case OAS_ | gc.TUINT16:
		a = sparc64.AMOVUH

	case OAS_ | gc.TINT32:
		a = sparc64.AMOVW

	case OAS_ | gc.TUINT32,
		OAS_ | gc.TPTR32:
		a = sparc64.AMOVUW

	case OAS_ | gc.TINT64,
		OAS_ | gc.TUINT64,
		OAS_ | gc.TPTR64:
		a = sparc64.AMOVD

	case OAS_ | gc.TFLOAT32:
		a = sparc64.AFMOVS

	case OAS_ | gc.TFLOAT64:
		a = sparc64.AFMOVD

	case OADD_ | gc.TINT8,
		OADD_ | gc.TUINT8,
		OADD_ | gc.TINT16,
		OADD_ | gc.TUINT16,
		OADD_ | gc.TINT32,
		OADD_ | gc.TUINT32,
		OADD_ | gc.TPTR32,
		OADD_ | gc.TINT64,
		OADD_ | gc.TUINT64,
		OADD_ | gc.TPTR64:
		a = sparc64.AADD

	case OADD_ | gc.TFLOAT32:
		a = sparc64.AFADDS

	case OADD_ | gc.TFLOAT64:
		a = sparc64.AFADDD

	case OSUB_ | gc.TINT8,
		OSUB_ | gc.TUINT8,
		OSUB_ | gc.TINT16,
		OSUB_ | gc.TUINT16,
		OSUB_ | gc.TINT32,
		OSUB_ | gc.TUINT32,
		OSUB_ | gc.TPTR32,
		OSUB_ | gc.TINT64,
		OSUB_ | gc.TUINT64,
		OSUB_ | gc.TPTR64:
		a = sparc64.ASUB

	case OSUB_ | gc.TFLOAT32:
		a = sparc64.AFSUBS

	case OSUB_ | gc.TFLOAT64:
		a = sparc64.AFSUBD

	case OMINUS_ | gc.TINT8,
		OMINUS_ | gc.TUINT8,
		OMINUS_ | gc.TINT16,
		OMINUS_ | gc.TUINT16,
		OMINUS_ | gc.TINT32,
		OMINUS_ | gc.TUINT32,
		OMINUS_ | gc.TPTR32,
		OMINUS_ | gc.TINT64,
		OMINUS_ | gc.TUINT64,
		OMINUS_ | gc.TPTR64:
		a = sparc64.ANEG

	case OMINUS_ | gc.TFLOAT32:
		a = sparc64.AFNEGS

	case OMINUS_ | gc.TFLOAT64:
		a = sparc64.AFNEGD

	case OAND_ | gc.TINT8,
		OAND_ | gc.TUINT8,
		OAND_ | gc.TINT16,
		OAND_ | gc.TUINT16,
		OAND_ | gc.TINT32,
		OAND_ | gc.TUINT32,
		OAND_ | gc.TPTR32,
		OAND_ | gc.TINT64,
		OAND_ | gc.TUINT64,
		OAND_ | gc.TPTR64:
		a = sparc64.AAND

	case OOR_ | gc.TINT8,
		OOR_ | gc.TUINT8,
		OOR_ | gc.TINT16,
		OOR_ | gc.TUINT16,
		OOR_ | gc.TINT32,
		OOR_ | gc.TUINT32,
		OOR_ | gc.TPTR32,
		OOR_ | gc.TINT64,
		OOR_ | gc.TUINT64,
		OOR_ | gc.TPTR64:
		a = sparc64.AOR

	case OXOR_ | gc.TINT8,
		OXOR_ | gc.TUINT8,
		OXOR_ | gc.TINT16,
		OXOR_ | gc.TUINT16,
		OXOR_ | gc.TINT32,
		OXOR_ | gc.TUINT32,
		OXOR_ | gc.TPTR32,
		OXOR_ | gc.TINT64,
		OXOR_ | gc.TUINT64,
		OXOR_ | gc.TPTR64:
		a = sparc64.AXOR

		// TODO(minux): handle rotates
	//case CASE(OLROT, TINT8):
	//case CASE(OLROT, TUINT8):
	//case CASE(OLROT, TINT16):
	//case CASE(OLROT, TUINT16):
	//case CASE(OLROT, TINT32):
	//case CASE(OLROT, TUINT32):
	//case CASE(OLROT, TPTR32):
	//case CASE(OLROT, TINT64):
	//case CASE(OLROT, TUINT64):
	//case CASE(OLROT, TPTR64):
	//	a = 0//???; RLDC?
	//	break;

	case OLSH_ | gc.TINT8,
		OLSH_ | gc.TUINT8,
		OLSH_ | gc.TINT16,
		OLSH_ | gc.TUINT16,
		OLSH_ | gc.TINT32,
		OLSH_ | gc.TUINT32,
		OLSH_ | gc.TPTR32:
		a = sparc64.ASLLW

	case OLSH_ | gc.TINT64,
		OLSH_ | gc.TUINT64,
		OLSH_ | gc.TPTR64:
		a = sparc64.ASLLD

	case ORSH_ | gc.TUINT8,
		ORSH_ | gc.TUINT16,
		ORSH_ | gc.TUINT32,
		ORSH_ | gc.TPTR32:
		a = sparc64.ASRLW

	case ORSH_ | gc.TUINT64,
		ORSH_ | gc.TPTR64:
		a = sparc64.ASRLD

	case ORSH_ | gc.TINT8,
		ORSH_ | gc.TINT16,
		ORSH_ | gc.TINT32:
		a = sparc64.ASRAW

	case ORSH_ | gc.TINT64:
		a = sparc64.ASRAD

	// TODO(shawn): handle rotates, likely via addcc/addxc and PSR icc
	// overflow bit
	//case CASE(ORROTC, TINT8):
	//case CASE(ORROTC, TUINT8):
	//case CASE(ORROTC, TINT16):
	//case CASE(ORROTC, TUINT16):
	//case CASE(ORROTC, TINT32):
	//case CASE(ORROTC, TUINT32):
	//case CASE(ORROTC, TINT64):
	//case CASE(ORROTC, TUINT64):
	//	a = 0//??? RLDC??
	//	break;

	// TODO(aram): handle high-multiply via mulx?
	//case OHMUL_ | gc.TINT64:
	//	a = sparc64.ASMULH
	//
	//case OHMUL_ | gc.TUINT64,
	//	OHMUL_ | gc.TPTR64:
	//	a = sparc64.AUMULH

	case OMUL_ | gc.TINT8,
		OMUL_ | gc.TINT16,
		OMUL_ | gc.TINT32,
		OMUL_ | gc.TINT64,
		OMUL_ | gc.TUINT8,
		OMUL_ | gc.TUINT16,
		OMUL_ | gc.TUINT32,
		OMUL_ | gc.TPTR32,
		OMUL_ | gc.TUINT64,
		OMUL_ | gc.TPTR64:
		a = sparc64.AMULD

	case OMUL_ | gc.TFLOAT32:
		a = sparc64.AFMULS

	case OMUL_ | gc.TFLOAT64:
		a = sparc64.AFMULD

	case ODIV_ | gc.TINT8,
		ODIV_ | gc.TINT16,
		ODIV_ | gc.TINT32,
		ODIV_ | gc.TINT64:
		a = sparc64.ASDIVD

	case ODIV_ | gc.TUINT8,
		ODIV_ | gc.TUINT16,
		ODIV_ | gc.TUINT32,
		ODIV_ | gc.TPTR32,
		ODIV_ | gc.TUINT64,
		ODIV_ | gc.TPTR64:
		a = sparc64.AUDIVD

	case ODIV_ | gc.TFLOAT32:
		a = sparc64.AFDIVS

	case ODIV_ | gc.TFLOAT64:
		a = sparc64.AFDIVD

	case OSQRT_ | gc.TFLOAT64:
		a = sparc64.AFSQRTD
	}

	return a
}

const (
	ODynam   = 1 << 0
	OAddable = 1 << 1
)

func xgen(n *gc.Node, a *gc.Node, o int) bool {
	// TODO(minux)

	return -1 != 0 /*TypeKind(100016)*/
}

func sudoclean() {
	return
}

/*
 * generate code to compute address of n,
 * a reference to a (perhaps nested) field inside
 * an array or struct.
 * return 0 on failure, 1 on success.
 * on success, leaves usable address in a.
 *
 * caller is responsible for calling sudoclean
 * after successful sudoaddable,
 * to release the register used for a.
 */
func sudoaddable(as obj.As, n *gc.Node, a *obj.Addr) bool {
	// TODO(minux)

	*a = obj.Addr{}
	return false
}

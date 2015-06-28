// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
	"fmt"
	"log"
)

// TODO(aram):
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {}

func follow(ctxt *obj.Link, s *obj.LSym) {
	ctxt.Cursym = s

	firstp := ctxt.NewProg()
	lastp := firstp
	xfol(ctxt, s.Text, &lastp)
	lastp.Link = nil
	s.Text = firstp.Link
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
	}

	log.Fatalf("unknown relation: %s", Anames[a])
	return 0
}

func xfol(ctxt *obj.Link, p *obj.Prog, last **obj.Prog) {
	var q *obj.Prog
	var r *obj.Prog
	var a int
	var i int

loop:
	if p == nil {
		return
	}
	a = int(p.As)
	if a == obj.AJMP {
		q = p.Pcond
		if q != nil {
			p.Mark |= FOLL
			p = q
			if !(p.Mark&FOLL != 0) {
				goto loop
			}
		}
	}

	if p.Mark&FOLL != 0 {
		i = 0
		q = p
		for ; i < 4; i, q = i+1, q.Link {
			if q == *last || q == nil {
				break
			}
			a = int(q.As)
			if a == obj.ANOP {
				i--
				continue
			}

			if a == obj.AJMP || a == obj.ARET {
				goto copy
			}
			if q.Pcond == nil || (q.Pcond.Mark&FOLL != 0) {
				continue
			}
			if a != ABE && a != ABNE {
				continue
			}

		copy:
			for {
				r = ctxt.NewProg()
				*r = *p
				if !(r.Mark&FOLL != 0) {
					fmt.Printf("cant happen 1\n")
				}
				r.Mark |= FOLL
				if p != q {
					p = p.Link
					(*last).Link = r
					*last = r
					continue
				}

				(*last).Link = r
				*last = r
				if a == obj.AJMP || a == obj.ARET {
					return
				}
				if a == ABNE {
					r.As = ABE
				} else {
					r.As = ABNE
				}
				r.Pcond = p.Link
				r.Link = p.Pcond
				if !(r.Link.Mark&FOLL != 0) {
					xfol(ctxt, r.Link, last)
				}
				if !(r.Pcond.Mark&FOLL != 0) {
					fmt.Printf("cant happen 2\n")
				}
				return
			}
		}

		a = obj.AJMP
		q = ctxt.NewProg()
		q.As = int16(a)
		q.Lineno = p.Lineno
		q.To.Type = obj.TYPE_BRANCH
		q.To.Offset = p.Pc
		q.Pcond = p
		p = q
	}

	p.Mark |= FOLL
	(*last).Link = p
	*last = p
	if a == obj.AJMP || a == obj.ARET {
		return
	}
	if p.Pcond != nil {
		if a != ABL && p.Link != nil {
			q = obj.Brchain(ctxt, p.Link)
			if a != obj.ATEXT {
				if q != nil && (q.Mark&FOLL != 0) {
					p.As = int16(relinv(a))
					p.Link = p.Pcond
					p.Pcond = q
				}
			}

			xfol(ctxt, p.Link, last)
			q = obj.Brchain(ctxt, p.Pcond)
			if q == nil {
				q = p.Pcond
			}
			if q.Mark&FOLL != 0 {
				p.Pcond = q
				return
			}

			p = q
			goto loop
		}
	}

	p = p.Link
	goto loop
}

var unaryDst = map[int]bool{
	AWORD:   true,
	ADWORD:  true,
	ARDPC:   true,
	ARDTICK: true,
	ARDCCR:  true,
}

var Linksparc64 = obj.LinkArch{
	ByteOrder:  binary.BigEndian,
	Name:       "sparc64",
	Thechar:    'u',
	Preprocess: preprocess,
	Follow:     follow,
	UnaryDst:   unaryDst,
	Minlc:      4,
	Ptrsize:    8,
	Regsize:    8,
}

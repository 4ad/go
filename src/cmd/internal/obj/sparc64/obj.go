// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
	"log"
)

// TODO(aram):
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {}

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

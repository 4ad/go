// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
)

// TODO(aram):
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {}

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
	UnaryDst:   unaryDst,
	Minlc:      4,
	Ptrsize:    8,
	Regsize:    8,
}

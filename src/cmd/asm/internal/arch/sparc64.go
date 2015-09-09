// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the SPARC64
// instruction set, to minimize its interaction with the core of the
// assembler.

package arch

import (
	"cmd/internal/obj/sparc64"
)

var sparc64Jump = map[string]bool{
	"ABN":    true,
	"ABNE":   true,
	"ABE":    true,
	"ABG":    true,
	"ABLE":   true,
	"ABGE":   true,
	"ABL":    true,
	"ABGU":   true,
	"ABLEU":  true,
	"ABCC":   true,
	"ABCS":   true,
	"ABPOS":  true,
	"ABNEG":  true,
	"ABVC":   true,
	"ABVS":   true,
	"ABRZ":   true,
	"ABRLEZ": true,
	"ABRLZ":  true,
	"ABRNZ":  true,
	"ABRGZ":  true,
	"ABRGEZ": true,
	"AFBA":   true,
	"AFBN":   true,
	"AFBU":   true,
	"AFBG":   true,
	"AFBUG":  true,
	"AFBL":   true,
	"AFBUL":  true,
	"AFBLG":  true,
	"AFBNE":  true,
	"AFBE":   true,
	"AFBUE":  true,
	"AFBGE":  true,
	"AFBUGE": true,
	"AFBLE":  true,
	"AFBULE": true,
	"AFBO":   true,
	"AJMPL":  true,
}

// IsSPARC64MEMBAR reports whether the op (as defined by an sparc64.A*
// constant) is one of the special MEMBAR instructions that require
// special handling.
func IsSPARC64MEMBAR(op int) bool {
	switch op {
	case sparc64.AMEMBAR:
		return true
	}
	return false
}

func jumpSparc64(word string) bool {
	return sparc64Jump[word]
}

func sparc64RegisterNumber(name string, n int16) (int16, bool) {
	switch name {
	case "D":
		if 0 <= n && n <= 30 && n%2 == 0 {
			return sparc64.REG_D0 + n, true
		}
		if 32 <= n && n <= 62 && n%2 == 0 {
			return sparc64.REG_D0 + n - 31, true
		}
	case "F":
		if 0 <= n && n <= 31 {
			return sparc64.REG_F0 + n, true
		}
	case "R":
		if 1 <= n && n <= 30 { // not 1, not 31
			return sparc64.REG_R0 + n, true
		}
	}
	return 0, false
}

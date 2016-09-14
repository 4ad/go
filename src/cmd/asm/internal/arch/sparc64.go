// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the SPARC64
// instruction set, to minimize its interaction with the core of the
// assembler.

package arch

import (
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
	"strings"
)

var sparc64Jump = map[string]bool{
	"BN":    true,
	"BNE":   true,
	"BE":    true,
	"BG":    true,
	"BLE":   true,
	"BGE":   true,
	"BL":    true,
	"BGU":   true,
	"BLEU":  true,
	"BCC":   true,
	"BCS":   true,
	"BPOS":  true,
	"BNEG":  true,
	"BVC":   true,
	"BVS":   true,
	"BNW":   true,
	"BNEW":  true,
	"BEW":   true,
	"BGW":   true,
	"BLEW":  true,
	"BGEW":  true,
	"BLW":   true,
	"BGUW":  true,
	"BLEUW": true,
	"BCCW":  true,
	"BCSW":  true,
	"BPOSW": true,
	"BNEGW": true,
	"BVCW":  true,
	"BVSW":  true,
	"BND":   true,
	"BNED":  true,
	"BED":   true,
	"BGD":   true,
	"BLED":  true,
	"BGED":  true,
	"BLD":   true,
	"BGUD":  true,
	"BLEUD": true,
	"BCCD":  true,
	"BCSD":  true,
	"BPOSD": true,
	"BNEGD": true,
	"BVCD":  true,
	"BVSD":  true,
	"BRZ":   true,
	"BRLEZ": true,
	"BRLZ":  true,
	"BRNZ":  true,
	"BRGZ":  true,
	"BRGEZ": true,
	"CALL":  true,
	"FBA":   true,
	"FBN":   true,
	"FBU":   true,
	"FBG":   true,
	"FBUG":  true,
	"FBL":   true,
	"FBUL":  true,
	"FBLG":  true,
	"FBNE":  true,
	"FBE":   true,
	"FBUE":  true,
	"FBGE":  true,
	"FBUGE": true,
	"FBLE":  true,
	"FBULE": true,
	"FBO":   true,
	"JMP":   true,
	"JMPL":  true,
}

// IsSPARC64CMP reports whether the op (as defined by an arm.A* constant) is
// one of the comparison instructions that require special handling.
func IsSPARC64CMP(op obj.As) bool {
	switch op {
	case sparc64.ACMP, sparc64.AFCMPD, sparc64.AFCMPS:
		return true
	}
	return false
}

func jumpSparc64(word string) bool {
	return sparc64Jump[word]
}

// SPARC64Suffix handles the special suffix for the SPARC64.
// It returns a boolean to indicate success; failure means
// cond was unrecognized.
func SPARC64Suffix(prog *obj.Prog, cond string) bool {
	if cond == "" {
		return true
	}
	bits, ok := ParseSPARC64Suffix(cond)
	if !ok {
		return false
	}
	prog.Scond = bits
	return true
}

// ParseSPARC64Suffix parses the suffix attached to an SPARC64 instruction.
func ParseSPARC64Suffix(cond string) (uint8, bool) {
	if cond == "" {
		return 0, true
	}
	if strings.HasPrefix(cond, ".") {
		cond = cond[1:]
	}
	names := strings.Split(cond, ".")
	if len(names) != 1 {
		return 0, false
	}
	if names[0] == "PN" {
		return 1, true
	}
	return 0, false
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
	case "G":
		if 0 <= n && n <= 5 { // not 6, 7
			return sparc64.REG_G0 + n, true
		}
	case "O":
		if 0 <= n && n <= 5 { // not 6, 7
			return sparc64.REG_O0 + n, true
		}
	case "L":
		if 0 <= n && n <= 7 {
			return sparc64.REG_L0 + n, true
		}
	case "I":
		if 0 <= n && n <= 5 { // not 6, 7
			return sparc64.REG_I0 + n, true
		}
	}
	return 0, false
}

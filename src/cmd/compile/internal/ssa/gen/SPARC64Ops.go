// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "strings"

// Notes:
//  - Less-than-64-bit integer types live in the low portion of registers.
//    For now, the upper portion is junk; sign/zero-extension might be optimized in the future, but not yet.
//  - Boolean types are zero or 1; stored in a byte, but loaded with AMOVUB so the upper bytes of a register are zero.
//  - *const instructions may use a constant larger than the instuction can encode.
//    In this case the assembler expands to multiple instructions and uses TMP
//    register.

var regNamesSPARC64 = []string{
	"RT1",
	"CTXT",
	"G",
	"RT2",
	"O0",
	"O1",
	"O2",
	"O3",
	"O4",
	"O5",
	"RSP",
	"L1",
	"L2",
	"L3",
	"L4",
	"L5",
	"L6",
	"RFP",
	"Y0",
	"Y1",
	"Y2",
	"Y3",
	"Y4",
	"Y5",
	"Y6",
	"Y7",
	"Y8",
	"Y9",
	"Y10",
	"Y11",
	"Y12",
	"Y13",

	// pseudo-registers
	"SB",
}

func init() {
	// Make map from reg names to reg integers.
	if len(regNamesSPARC64) > 64 {
		panic("too many registers")
	}
	num := map[string]int{}
	for i, name := range regNamesSPARC64 {
		num[name] = i
	}
	buildReg := func(s string) regMask {
		m := regMask(0)
		for _, r := range strings.Split(s, " ") {
			if n, ok := num[r]; ok {
				m |= regMask(1) << uint(n)
				continue
			}
			panic("register " + r + " not found")
		}
		return m
	}

	var (
		gp = buildReg("O0 O1 O2 O3 O4 O5 L1 L2 L3 L4 L5 L6")
		fp = buildReg("Y0 Y1 Y2 Y3 Y4 Y5 Y6 Y7 Y8 Y9 Y10 Y11 Y12 Y13")

		gp11 = regInfo{inputs: []regMask{gp}, outputs: []regMask{gp}}
		gp21 = regInfo{inputs: []regMask{gp, gp}, outputs: []regMask{gp}}
		fp21 = regInfo{inputs: []regMask{fp, fp}, outputs: []regMask{fp}}
	)
	ops := []opData{
		{name: "ADD", argLength: 2, reg: gp21, asm: "ADD", commutative: true}, // arg0 + arg1
		{name: "ADDconst", argLength: 1, reg: gp11, asm: "ADD", aux: "Int64"}, // arg0 + auxInt
		{name: "SUB", argLength: 2, reg: gp21, asm: "SUB"}, // arg0 - arg1
		{name: "SUBconst", argLength: 1, reg: gp11, asm: "SUB", aux: "Int64"}, // arg0 - auxInt
		{name: "MULD", argLength: 2, reg: gp21, asm: "MULD", commutative: true},     // arg0 * arg1

		{name: "FADDS", argLength: 2, reg: fp21, asm: "FADDS", commutative: true}, // arg0 + arg1
		{name: "FADDD", argLength: 2, reg: fp21, asm: "FADDD", commutative: true}, // arg0 + arg1
		{name: "FSUBS", argLength: 2, reg: fp21, asm: "FSUBS"},                    // arg0 - arg1
		{name: "FSUBD", argLength: 2, reg: fp21, asm: "FSUBD"},                    // arg0 - arg1
		{name: "FMULS", argLength: 2, reg: fp21, asm: "FMULS", commutative: true}, // arg0 * arg1
		{name: "FMULD", argLength: 2, reg: fp21, asm: "FMULD", commutative: true}, // arg0 * arg1
	}

	blocks := []blockData{
		{name: "N"},
		{name: "NE"},
		{name: "E"},
		{name: "G"},
		{name: "LE"},
		{name: "GE"},
		{name: "L"},
		{name: "GU"},
		{name: "LEU"},
		{name: "CC"},
		{name: "CS"},
		{name: "POS"},
		{name: "NEG"},
		{name: "VC"},
		{name: "VS"},
		// TODO(aram): float?
	}

	archs = append(archs, arch{
		name:            "SPARC64",
		pkg:             "cmd/internal/obj/sparc64",
		genfile:         "../../sparc64/ssa.go",
		ops:             ops,
		blocks:          blocks,
		regnames:        regNamesSPARC64,
		gpregmask:       gp,
		fpregmask:       fp,
		framepointerreg: int8(num["RFP"]),
	})
}

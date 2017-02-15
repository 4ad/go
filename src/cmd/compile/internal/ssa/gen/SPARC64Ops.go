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
	"g",
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
	"SP",
	"FP",
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

		gp01 = regInfo{inputs: nil, outputs: []regMask{gp}}
		gp11 = regInfo{inputs: []regMask{gp}, outputs: []regMask{gp}}
		gp21 = regInfo{inputs: []regMask{gp, gp}, outputs: []regMask{gp}}
		fp01        = regInfo{inputs: nil, outputs: []regMask{fp}}
		fp11 = regInfo{inputs: []regMask{fp}, outputs: []regMask{fp}}
		fp21 = regInfo{inputs: []regMask{fp, fp}, outputs: []regMask{fp}}
	)
	ops := []opData{
		{name: "ADD", argLength: 2, reg: gp21, asm: "ADD", commutative: true}, // arg0 + arg1
		{name: "SUB", argLength: 2, reg: gp21, asm: "SUB"}, // arg0 - arg1

		{name: "AND", argLength: 2, reg: gp21, asm: "AND", commutative: true}, // arg0 & arg1
		{name: "OR", argLength: 2, reg: gp21, asm: "OR", commutative: true},  // arg0 | arg1
		{name: "XOR", argLength: 2, reg: gp21, asm: "XOR", commutative: true}, // arg0 ^ arg1

		{name: "ADDconst", argLength: 1, reg: gp11, asm: "ADD", aux: "Int64"}, // arg0 + auxInt
		{name: "SUBconst", argLength: 1, reg: gp11, asm: "SUB", aux: "Int64"}, // arg0 - auxInt
		{name: "ANDconst", argLength: 1, reg: gp11, asm: "AND", aux: "Int64"}, // arg0 & auxInt
		{name: "ORconst", argLength: 1, reg: gp11, asm: "OR", aux: "Int64"},  // arg0 | auxInt
		{name: "XORconst", argLength: 1, reg: gp11, asm: "XOR", aux: "Int64"}, // arg0 ^ auxInt

		{name: "MULD", argLength: 2, reg: gp21, typ: "Int64", asm: "MULD", commutative: true},     // arg0 * arg1
		{name: "SDIVD", argLength: 2, reg: gp21, typ: "Int64", asm: "SDIVD"},                       // arg0 / arg1, signed
		{name: "UDIVD", argLength: 2, reg: gp21, typ: "UInt64", asm: "UDIVD"},                       // arg0 / arg1, unsigned

		{name: "FADDS", argLength: 2, reg: fp21, asm: "FADDS", commutative: true}, // arg0 + arg1
		{name: "FADDD", argLength: 2, reg: fp21, asm: "FADDD", commutative: true}, // arg0 + arg1
		{name: "FSUBS", argLength: 2, reg: fp21, asm: "FSUBS"},                    // arg0 - arg1
		{name: "FSUBD", argLength: 2, reg: fp21, asm: "FSUBD"},                    // arg0 - arg1
		{name: "FMULS", argLength: 2, reg: fp21, asm: "FMULS", commutative: true}, // arg0 * arg1
		{name: "FMULD", argLength: 2, reg: fp21, asm: "FMULD", commutative: true}, // arg0 * arg1
		{name: "FDIVS", argLength: 2, reg: fp21, asm: "FDIVS"},                    // arg0 / arg1
		{name: "FDIVD", argLength: 2, reg: fp21, asm: "FDIVD"},                    // arg0 / arg1

		// unary ops
		{name: "NEG", argLength: 1, reg: gp11, asm: "NEG"},       // -arg0
		{name: "FNEGS", argLength: 1, reg: fp11, asm: "FNEGS"},   // -arg0, float32
		{name: "FNEGD", argLength: 1, reg: fp11, asm: "FNEGD"},   // -arg0, float64
		{name: "FSQRTD", argLength: 1, reg: fp11, asm: "FSQRTD"}, // sqrt(arg0), float64
		// moves
		{name: "MOVDconst", argLength: 0, reg: gp01, aux: "Int64", asm: "MOVD", rematerializeable: true},
		{name: "MOVWconst", argLength: 0, reg: gp01, aux: "Int32", asm: "MOVW", rematerializeable: true},     // 32 low bits of auxint
		{name: "FMOVDconst", argLength: 0, reg: fp01, aux: "Float64", asm: "FMOVD", rematerializeable: true},
		{name: "FMOVSconst", argLength: 0, reg: fp01, aux: "Float32", asm: "FMOVS", rematerializeable: true},

		// conversions
		{name: "MOVBreg", argLength: 1, reg: gp11, asm: "MOVB"},   // move from arg0, sign-extended from byte
		{name: "MOVUBreg", argLength: 1, reg: gp11, asm: "MOVUB"}, // move from arg0, unsign-extended from byte
		{name: "MOVHreg", argLength: 1, reg: gp11, asm: "MOVH"},   // move from arg0, sign-extended from half
		{name: "MOVUHreg", argLength: 1, reg: gp11, asm: "MOVUH"}, // move from arg0, unsign-extended from half
		{name: "MOVWreg", argLength: 1, reg: gp11, asm: "MOVW"},   // move from arg0, sign-extended from word
		{name: "MOVUWreg", argLength: 1, reg: gp11, asm: "MOVUW"}, // move from arg0, unsign-extended from word
		{name: "MOVDreg", argLength: 1, reg: gp11, asm: "MOVD"},   // move from arg0
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

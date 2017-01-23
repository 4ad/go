// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
)

var ssaRegToReg = []int16{
	sparc64.REG_ZR,
	sparc64.REG_RT1,
	sparc64.REG_CTXT,
	sparc64.REG_G,
	sparc64.REG_RT2,
	sparc64.REG_TMP,
	sparc64.REG_TLS,
	sparc64.REG_G7,
	sparc64.REG_O0,
	sparc64.REG_O1,
	sparc64.REG_O2,
	sparc64.REG_O3,
	sparc64.REG_O4,
	sparc64.REG_O5,
	sparc64.REG_RSP,
	sparc64.REG_OLR,
	sparc64.REG_TMP2,
	sparc64.REG_L1,
	sparc64.REG_L2,
	sparc64.REG_L3,
	sparc64.REG_L4,
	sparc64.REG_L5,
	sparc64.REG_L6,
	sparc64.REG_L7,
	sparc64.REG_I0,
	sparc64.REG_I1,
	sparc64.REG_I2,
	sparc64.REG_I3,
	sparc64.REG_I4,
	sparc64.REG_I5,
	sparc64.REG_RFP,
	sparc64.REG_ILR,

	sparc64.REG_Y0,
	sparc64.REG_Y1,
	sparc64.REG_Y2,
	sparc64.REG_Y3,
	sparc64.REG_Y4,
	sparc64.REG_Y5,
	sparc64.REG_Y6,
	sparc64.REG_Y7,
	sparc64.REG_Y8,
	sparc64.REG_Y9,
	sparc64.REG_Y10,
	sparc64.REG_Y11,
	sparc64.REG_Y12,
	sparc64.REG_Y13,
	sparc64.REG_YTWO,
	sparc64.REG_YTMP,
}

// Smallest possible faulting page at address zero,
// see ../../../../runtime/mheap.go:/minPhysPageSize
const minZeroPage = 4096

// loadByType returns the load instruction of the given type.
func loadByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return sparc64.AFMOVS
		case 8:
			return sparc64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			if t.IsSigned() {
				return sparc64.AMOVB
			} else {
				return sparc64.AMOVUB
			}
		case 2:
			if t.IsSigned() {
				return sparc64.AMOVH
			} else {
				return sparc64.AMOVUH
			}
		case 4:
			if t.IsSigned() {
				return sparc64.AMOVW
			} else {
				return sparc64.AMOVUW
			}
		case 8:
			return sparc64.AMOVD
		}
	}
	panic("bad load type")
}

// storeByType returns the store instruction of the given type.
func storeByType(t ssa.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return sparc64.AFMOVS
		case 8:
			return sparc64.AFMOVD
		}
	} else {
		switch t.Size() {
		case 1:
			return sparc64.AMOVB
		case 2:
			return sparc64.AMOVH
		case 4:
			return sparc64.AMOVW
		case 8:
			return sparc64.AMOVD
		}
	}
	panic("bad store type")
}

func ssaGenValue(s *gc.SSAGenState, v *ssa.Value) {
	s.SetLineno(v.Line)
	switch v.Op {
	default:
		v.Unimplementedf("genValue not implemented: %s", v.LongString())
	}
}

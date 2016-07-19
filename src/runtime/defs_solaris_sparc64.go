// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

/*
Input to cgo.

GOARCH=sparc64 go tool cgo -godefs defs_solaris_sparc64.go >defs1_solaris_sparc64.go
*/

package runtime

/*
#include <sys/types.h>
#include <sys/regset.h>
*/
import "C"

const (
	REG_CCR = C.REG_CCR
	REG_PC  = C.REG_PC
	REG_nPC = C.REG_nPC
	REG_G1  = C.REG_G1
	REG_G2  = C.REG_G2
	REG_G3  = C.REG_G3
	REG_G4  = C.REG_G4
	REG_G5  = C.REG_G5
	REG_G6  = C.REG_G6
	REG_G7  = C.REG_G7
	REG_O0  = C.REG_O0
	REG_O1  = C.REG_O1
	REG_O2  = C.REG_O2
	REG_O3  = C.REG_O3
	REG_O4  = C.REG_O4
	REG_O5  = C.REG_O5
	REG_O6  = C.REG_O6
	REG_O7  = C.REG_O7
)

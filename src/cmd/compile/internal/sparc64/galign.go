// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/internal/obj/sparc64"
)

var MAXWIDTH int64 = 1 << 50

func betypeinit() {
}

func Main() {
	gc.Thearch.LinkArch = &sparc64.Linksparc64

	gc.Thearch.REGSP = sparc64.REG_BSP
	gc.Thearch.REGCTXT = sparc64.REG_CTXT
	gc.Thearch.REGCALLX = sparc64.REG_RT1
	gc.Thearch.REGCALLX2 = sparc64.REG_RT2
	gc.Thearch.REGRETURN = sparc64.REG_RT1
	gc.Thearch.REGMIN = sparc64.REG_MIN
	gc.Thearch.REGMAX = sparc64.REG_MAX
	gc.Thearch.REGZERO = sparc64.REG_ZR
	gc.Thearch.FREGMIN = sparc64.REG_Y0
	gc.Thearch.FREGMAX = sparc64.REG_Y15
	gc.Thearch.MAXWIDTH = MAXWIDTH
	gc.Thearch.ReservedRegs = resvd

	gc.Thearch.Betypeinit = betypeinit
	gc.Thearch.Cgen_hmul = cgen_hmul
	gc.Thearch.Cgen_shift = cgen_shift
	gc.Thearch.Clearfat = clearfat
	gc.Thearch.Defframe = defframe
	gc.Thearch.Dodiv = dodiv
	gc.Thearch.Excise = excise
	gc.Thearch.Expandchecks = expandchecks
	gc.Thearch.Getg = getg
	gc.Thearch.Gins = gins
	gc.Thearch.Ginscmp = ginscmp
	gc.Thearch.Ginscon = ginscon
	gc.Thearch.Ginsnop = ginsnop
	gc.Thearch.Gmove = gmove
	gc.Thearch.Peep = peep
	gc.Thearch.Proginfo = proginfo
	gc.Thearch.Regtyp = regtyp
	gc.Thearch.Sameaddr = sameaddr
	gc.Thearch.Smallindir = smallindir
	gc.Thearch.Stackaddr = stackaddr
	gc.Thearch.Blockcopy = blockcopy
	gc.Thearch.Sudoaddable = sudoaddable
	gc.Thearch.Sudoclean = sudoclean
	gc.Thearch.Excludedregs = excludedregs
	gc.Thearch.RtoB = RtoB
	gc.Thearch.FtoB = RtoB
	gc.Thearch.BtoR = BtoR
	gc.Thearch.BtoF = BtoF
	gc.Thearch.Optoas = optoas
	gc.Thearch.Doregbits = doregbits
	gc.Thearch.Regnames = regnames

	gc.Thearch.SSARegToReg = ssaRegToReg
	gc.Thearch.SSAMarkMoves = func(s *gc.SSAGenState, b *ssa.Block) {}
	gc.Thearch.SSAGenValue = ssaGenValue
	gc.Thearch.SSAGenBlock = ssaGenBlock

	gc.Main()
	gc.Exit(0)
}

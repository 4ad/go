// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
)

const (
	LeftRdwr  uint32 = gc.LeftRead | gc.LeftWrite
	RightRdwr uint32 = gc.RightRead | gc.RightWrite
)

// This table gives the basic information about instruction
// generated by the compiler and processed in the optimizer.
// See opt.h for bit definitions.
//
// Instructions not generated need not be listed.
// As an exception to that rule, we typically write down all the
// size variants of an operation even if we just use a subset.
//
// The table is formatted for 8-space tabs.
var progtable = [sparc64.ALAST]obj.ProgInfo{
	obj.ATYPE:     {Flags: gc.Pseudo | gc.Skip},
	obj.ATEXT:     {Flags: gc.Pseudo},
	obj.AFUNCDATA: {Flags: gc.Pseudo},
	obj.APCDATA:   {Flags: gc.Pseudo},
	obj.AUNDEF:    {Flags: gc.Break},
	obj.AUSEFIELD: {Flags: gc.OK},
	obj.ACHECKNIL: {Flags: gc.LeftRead},
	obj.AVARDEF:   {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARKILL:  {Flags: gc.Pseudo | gc.RightWrite},

	// NOP is an internal no-op that also stands
	// for USED and SET annotations, not the Power opcode.
	obj.ANOP:      {Flags: gc.LeftRead | gc.RightWrite},
	sparc64.AHINT: {Flags: gc.OK},

	// Integer
	sparc64.AADD:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ASUB:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ANEG:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AAND:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AORR:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AXOR:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AMUL:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ASMULL: {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AUMULL: {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ASMULH: {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AUMULH: {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ASDIV:  {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AUDIV:  {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ALSL:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ALSR:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AASR:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.ACMP:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead},

	// Floating point.
	sparc64.AFADDD:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFADDS:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFSUBD:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFSUBS:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFNEGD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite},
	sparc64.AFNEGS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite},
	sparc64.AFSQRTD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite},
	sparc64.AFMULD:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFMULS:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFDIVD:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFDIVS:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	sparc64.AFCMPD:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead},
	sparc64.AFCMPS:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead},

	// float -> integer
	sparc64.AFCVTZSD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZSS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZSDW: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZSSW: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZUD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZUS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZUDW: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTZUSW: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},

	// float -> float
	sparc64.AFCVTSD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AFCVTDS: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},

	// integer -> float
	sparc64.ASCVTFD:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.ASCVTFS:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.ASCVTFWD: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.ASCVTFWS: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AUCVTFD:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AUCVTFS:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AUCVTFWD: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	sparc64.AUCVTFWS: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},

	// Moves
	sparc64.AMOVB:  {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVBU: {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVH:  {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVHU: {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVW:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVWU: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AMOVD:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Move},
	sparc64.AFMOVS: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	sparc64.AFMOVD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move},

	// Jumps
	sparc64.AB:    {Flags: gc.Jump | gc.Break},
	sparc64.ABL:   {Flags: gc.Call},
	sparc64.ABEQ:  {Flags: gc.Cjmp},
	sparc64.ABNE:  {Flags: gc.Cjmp},
	sparc64.ABGE:  {Flags: gc.Cjmp},
	sparc64.ABLT:  {Flags: gc.Cjmp},
	sparc64.ABGT:  {Flags: gc.Cjmp},
	sparc64.ABLE:  {Flags: gc.Cjmp},
	sparc64.ABLO:  {Flags: gc.Cjmp},
	sparc64.ABLS:  {Flags: gc.Cjmp},
	sparc64.ABHI:  {Flags: gc.Cjmp},
	sparc64.ABHS:  {Flags: gc.Cjmp},
	sparc64.ACBZ:  {Flags: gc.Cjmp},
	sparc64.ACBNZ: {Flags: gc.Cjmp},
	obj.ARET:      {Flags: gc.Break},
	obj.ADUFFZERO: {Flags: gc.Call},
	obj.ADUFFCOPY: {Flags: gc.Call},
}

func proginfo(p *obj.Prog) {
	info := &p.Info
	*info = progtable[p.As]
	if info.Flags == 0 {
		gc.Fatalf("proginfo: unknown instruction %v", p)
	}

	if (info.Flags&gc.RegRead != 0) && p.Reg == 0 {
		info.Flags &^= gc.RegRead
		info.Flags |= gc.RightRead /*CanRegRead |*/
	}

	if (p.From.Type == obj.TYPE_MEM || p.From.Type == obj.TYPE_ADDR) && p.From.Reg != 0 {
		info.Regindex |= RtoB(int(p.From.Reg))
		if p.Scond != 0 {
			info.Regset |= RtoB(int(p.From.Reg))
		}
	}

	if (p.To.Type == obj.TYPE_MEM || p.To.Type == obj.TYPE_ADDR) && p.To.Reg != 0 {
		info.Regindex |= RtoB(int(p.To.Reg))
		if p.Scond != 0 {
			info.Regset |= RtoB(int(p.To.Reg))
		}
	}

	if p.From.Type == obj.TYPE_ADDR && p.From.Sym != nil && (info.Flags&gc.LeftRead != 0) {
		info.Flags &^= gc.LeftRead
		info.Flags |= gc.LeftAddr
	}

	if p.As == obj.ADUFFZERO {
		info.Reguse |= RtoB(sparc64.REGRT1)
		info.Regset |= RtoB(sparc64.REGRT1)
	}

	if p.As == obj.ADUFFCOPY {
		// TODO(austin) Revisit when duffcopy is implemented
		info.Reguse |= RtoB(sparc64.REGRT1) | RtoB(sparc64.REGRT2) | RtoB(sparc64.REG_R5)

		info.Regset |= RtoB(sparc64.REGRT1) | RtoB(sparc64.REGRT2)
	}
}

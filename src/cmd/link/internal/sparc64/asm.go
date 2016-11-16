// Inferno utils/5l/asm.c
// http://code.google.com/p/inferno-os/source/browse/utils/5l/asm.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sparc64

import (
	"cmd/internal/obj"
	"cmd/link/internal/ld"
	"fmt"
	"log"
)

func gentext() {
	if !ld.DynlinkingGo() {
		return
	}
	log.Fatalf("gentext() not implemented")
}

func addgotsym(s *ld.LSym) {
	if s.Got >= 0 {
		return
	}

	ld.Adddynsym(ld.Ctxt, s)
	got := ld.Linklookup(ld.Ctxt, ".got", 0)
	s.Got = int32(got.Size)
	ld.Adduint64(ld.Ctxt, got, 0)

	if ld.Iself {
		rela := ld.Linklookup(ld.Ctxt, ".rela.got", 0)
		ld.Addaddrplus(ld.Ctxt, rela, got, int64(s.Got))
		ld.Adduint64(ld.Ctxt, rela, ld.ELF64_R_INFO(uint32(s.Dynid), ld.R_SPARC_GLOB_DAT))
		ld.Adduint64(ld.Ctxt, rela, 0)
	} else {
		ld.Diag("addgotsym: unsupported binary format")
	}
}

func adddynrela(rela *ld.LSym, s *ld.LSym, r *ld.Reloc) {
	log.Fatalf("adddynrela not implemented")
}

func adddynrel(s *ld.LSym, r *ld.Reloc) {
	targ := r.Sym
	ld.Ctxt.Cursym = s

	switch r.Type {
	default:
		if r.Type >= 256 {
			ld.Diag("unexpected relocation type %d (%d)", r.Type, r.Type-256)
			return
		}
	case 256 + ld.R_SPARC_PC10:
		if targ.Type == obj.SDYNIMPORT {
			ld.Diag("unexpected R_SPARC_PC10 relocation for dynamic symbol %s", targ.Name)
		}
		if targ.Type == 0 || targ.Type == obj.SXREF {
			ld.Diag("unknown symbol %s in pcrel", targ.Name)
		}
		r.Type = obj.R_PCREL
		r.Add += int64(r.Siz)
		println("R_SPARC_PC10 relocation for symbol ", targ.Name)
		return

	case 256 + ld.R_SPARC_PC22:
		if targ.Type == obj.SDYNIMPORT {
			ld.Diag("unexpected R_SPARC_PC22 relocation for dynamic symbol %s", targ.Name)
		}
		if targ.Type == 0 || targ.Type == obj.SXREF {
			ld.Diag("unknown symbol %s in pcrel", targ.Name)
		}
		r.Type = obj.R_PCREL
		r.Add += int64(r.Siz)
		println("R_SPARC_PC22 relocation for symbol ", targ.Name)
		return

	case 256 + ld.R_SPARC_WPLT30:
		r.Add += int64(r.Siz)
		if targ.Type == obj.SDYNIMPORT {
			addpltsym(targ)
			r.Sym = ld.Linklookup(ld.Ctxt, ".plt", 0)
			r.Add += int64(targ.Plt)
		}
		r.Type = obj.R_CALLSPARC64
		println("R_SPARC_WPLT30 relocation for symbol ", targ.Name)
		fmt.Printf("r %#v\n", r)
		return

	// TODO(shawn):
	// The R_SPARC_GOTDATA_OP* relocations are an optimized form of
	// relocation that only supports a range of +/- 2 Gbytes.  We should
	// eventually support these, but for now, simplify them to standard GOT
	// relocations for simplicity in implementation.
	case 256 + ld.R_SPARC_GOT10, 256 + ld.R_SPARC_GOTDATA_OP_LOX10:
		addgotsym(targ)
		r.Sym = ld.Linklookup(ld.Ctxt, ".got", 0)
		r.Add += int64(targ.Got)
		r.Type = obj.R_GOTOFF
		println("R_SPARC_GOT10 relocation for symbol ", targ.Name)
		return

	case 256 + ld.R_SPARC_GOT22, 256 + ld.R_SPARC_GOTDATA_OP_HIX22:
		addgotsym(targ)
		r.Sym = ld.Linklookup(ld.Ctxt, ".got", 0)
		r.Add += int64(targ.Got)
		r.Type = obj.R_GOTOFF
		println("R_SPARC_GOT22 relocation for symbol ", targ.Name)
		return

	case 256 + ld.R_SPARC_GOTDATA_OP:
		r.Type = ld.R_SPARC_NONE
		return
	}

	// Handle references to ELF symbols from our own object files.
	if targ.Type != obj.SDYNIMPORT {
		return
	}

	switch r.Type {
	case obj.R_ADDRSPARC64HI, obj.R_ADDRSPARC64LO:
		if s.Type == obj.STEXT && ld.Iself {
			addpltsym(targ)
			r.Sym = ld.Linkrlookup(ld.Ctxt, ".plt", 0)
			r.Add += int64(targ.Plt)
			return
		}
	}

	ld.Ctxt.Cursym = s
	ld.Diag("unsupported relocation for dynamic symbol %s (type=%d stype=%d)", targ.Name, r.Type, targ.Type)

}

func elfreloc1(r *ld.Reloc, sectoff int64) int {
	ld.Thearch.Vput(uint64(sectoff))

	elfsym := r.Xsym.ElfsymForReloc()
	switch r.Type {
	default:
		return -1

	case obj.R_ADDR:
		switch r.Siz {
		case 4:
			ld.Thearch.Vput(ld.R_SPARC_32 | uint64(elfsym)<<32)
		case 8:
			ld.Thearch.Vput(ld.R_SPARC_64 | uint64(elfsym)<<32)
		default:
			return -1
		}

	case obj.R_ADDRSPARC64LO:
		ld.Thearch.Vput(ld.R_SPARC_LM22 | uint64(elfsym)<<32)
		ld.Thearch.Vput(uint64(r.Xadd))
		ld.Thearch.Vput(uint64(sectoff + 4))
		ld.Thearch.Vput(ld.R_SPARC_LO10 | uint64(elfsym)<<32)

	case obj.R_ADDRSPARC64HI:
		ld.Thearch.Vput(ld.R_SPARC_HH22 | uint64(elfsym)<<32)
		ld.Thearch.Vput(uint64(r.Xadd))
		ld.Thearch.Vput(uint64(sectoff + 4))
		ld.Thearch.Vput(ld.R_SPARC_HM10 | uint64(elfsym)<<32)

	case obj.R_SPARC64_TLS_LE:
		ld.Thearch.Vput(ld.R_SPARC_TLS_LE_HIX22 | uint64(elfsym)<<32)
		ld.Thearch.Vput(uint64(r.Xadd))
		ld.Thearch.Vput(uint64(sectoff + 4))
		ld.Thearch.Vput(ld.R_SPARC_TLS_LE_LOX10 | uint64(elfsym)<<32)

	case obj.R_CALLSPARC64:
		if r.Siz != 4 {
			return -1
		}
		ld.Thearch.Vput(ld.R_SPARC_WDISP30 | uint64(elfsym)<<32)

	case obj.R_PCREL:
		if r.Siz != 4 {
			return -1
		}
		ld.Thearch.Vput(ld.R_SPARC_RELATIVE | uint64(elfsym)<<32)
	}
	ld.Thearch.Vput(uint64(r.Xadd))

	return 0
}

func elfsetupplt() {
	plt := ld.Linklookup(ld.Ctxt, ".plt", 0)
	if plt.Size == 0 {
		// .plt entries are aligned at 32-byte boundaries, but the
		// entire section at 256-byte boundaries.
		plt.Align = 256

		// Runtime linker will provide the initial plt; each entry is
		// 32 bytes; reserve the first four entries for its use.
		plt.Size = 4 * 32

		// Create relocation table for .plt
		rela := ld.Linklookup(ld.Ctxt, ".rela.plt", 0)
		rela.Align = int32(ld.SysArch.RegSize)

		// Create global offset table; first entry reserved for
		// address of .dynamic section.
		got := ld.Linklookup(ld.Ctxt, ".got", 0)
		dyn := ld.Linklookup(ld.Ctxt, ".dynamic", 0)
		ld.Addaddrplus(ld.Ctxt, got, dyn, 0)

		// TODO(srwalker): pad end of plt with 10 entries for elfedit, etc.
		// .strtab too; aslr-tagging, etc.
	}
}

func machoreloc1(r *ld.Reloc, sectoff int64) int {
	log.Fatalf("machoreloc1 not implemented")
	return 0
}

func archreloc(r *ld.Reloc, s *ld.LSym, val *int64) int {
	if ld.Linkmode == ld.LinkExternal {
		switch r.Type {
		default:
			ld.Diag("unsupported LinkExternal archreloc %s", r.Type)
			return -1

		case obj.R_ADDRSPARC64LO, obj.R_ADDRSPARC64HI:
			r.Done = 0

			// set up addend for eventual relocation via outer symbol.
			rs := r.Sym
			r.Xadd = r.Add
			for rs.Outer != nil {
				r.Xadd += ld.Symaddr(rs) - ld.Symaddr(rs.Outer)
				rs = rs.Outer
			}

			if rs.Type != obj.SHOSTOBJ && rs.Type != obj.SDYNIMPORT && rs.Sect == nil {
				ld.Diag("missing section for %s", rs.Name)
			}
			r.Xsym = rs

			return 0

		case obj.R_CALLSPARC64, obj.R_SPARC64_TLS_LE:
			r.Done = 0
			r.Xsym = r.Sym
			r.Xadd = r.Add
			return 0
		}
	}

	switch r.Type {
	case obj.R_CONST:
		*val = r.Add
		return 0

	case obj.R_ADDRSPARC64LO:
		t := ld.Symaddr(r.Sym) + r.Add

		o0 := uint32(*val >> 32)
		o1 := uint32(*val)

		o0 |= uint32(t) >> 10
		o1 |= uint32(t) & 0x3ff

		*val = int64(o0)<<32 | int64(o1)
		return 0

	case obj.R_ADDRSPARC64HI:
		t := ld.Symaddr(r.Sym) + r.Add

		o0 := uint32(*val >> 32)
		o1 := uint32(*val)

		o0 |= uint32(uint64(t)>>32) >> 10
		o1 |= uint32(uint64(t)>>32) & 0x3ff

		*val = int64(o0)<<32 | int64(o1)
		return 0

	case obj.R_CALLSPARC64:
		t := (ld.Symaddr(r.Sym) + r.Add) - (s.Value + int64(r.Off))
		if t > 1<<31-4 || t < -1<<31 {
			ld.Diag("program too large, call relocation distance = %d", t)
		}
		*val |= (t >> 2) & 0x3fffffff
		return 0

	case obj.R_GOTOFF:
		// TODO(shawn): GOT10 needs (val) & 0x3ff
		// GOT22 needs (val) >> 10
		*val = ld.Symaddr(r.Sym) + r.Add - ld.Symaddr(ld.Linklookup(ld.Ctxt, ".got", 0))
		return 0

	case obj.R_SPARC64_TLS_LE:
		// The thread pointer points to the TCB, and then the
		// address of the first TLS block follows, giving an
		// offset of -16 for our static TLS variables.
		v := r.Sym.Value - 16
		if v < -4096 || 4095 < v {
			ld.Diag("TLS offset out of range %d", v)
		}
		*val = (*val &^ 0x1fff) | (v & 0x1fff)
		return 0
	}

	ld.Diag("unsupported LinkInternal archreloc %s", r.Type)
	return -1
}

func archrelocvariant(r *ld.Reloc, s *ld.LSym, t int64) int64 {
	log.Fatalf("unexpected relocation variant")
	return -1
}

func addpltsym(s *ld.LSym) {
	if s.Plt >= 0 {
		return
	}

	ld.Adddynsym(ld.Ctxt, s)

	if ld.Iself {
		elfsetupplt()
		plt := ld.Linkrlookup(ld.Ctxt, ".plt", 0)
		rela := ld.Linkrlookup(ld.Ctxt, ".rela.plt", 0)

		// Each of the first 32,767 procedure linkage table entries occupies
		// 8 words (32 bytes), and must be aligned on a 32-byte boundary.
		//
		// NOTE: This only supports "near" .plt; entries beyond
		// 32,767 are considered "far" and have a different format.

		// The first eight bytes of each entry (excluding the initially
		// reserved ones) should transfer control to the first or second
		// reserved plt entry.  For our use, the second reserved entry (.PLT1)
		// should always be the target.
		//
		// 03 00 00 80 sethi (.-.PLT0), %g1 sethi     %hi(0x20000), %g1
		sethi := uint32(0x03000000)
		sethi |= uint32(plt.Size)
		ld.Adduint32(ld.Ctxt, plt, sethi)

		// 30 6f ff e7 ba,a,pt   %xcc, .PLT1
		ba := uint32(0x30680000)
		ba |= (((-uint32(plt.Size)) + 32) >> 2) & ((1 << (19)) - 1)
		ld.Adduint32(ld.Ctxt, plt, ba)

		// Fill remaining 24 bytes with nop; these will be provided by the
		// runtime linker.
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)
		ld.Adduint32(ld.Ctxt, plt, 0x01000000)

		// rela
		// offset
		ld.Addaddrplus(ld.Ctxt, rela, plt, plt.Size-32)
		// info
		ld.Adduint64(ld.Ctxt, rela, ld.ELF64_R_INFO(uint32(s.Dynid), ld.R_SPARC_JMP_SLOT))
		// addend
		ld.Adduint64(ld.Ctxt, rela, 0)

		s.Plt = int32(plt.Size - 32)
	} else {
		ld.Diag("addpltsym: unsupported binary format")
	}

	return
}

func asmb() {
	if ld.Debug['v'] != 0 {
		fmt.Fprintf(ld.Bso, "%5.2f asmb\n", obj.Cputime())
	}
	ld.Bso.Flush()

	if ld.Iself {
		ld.Asmbelfsetup()
	}

	sect := ld.Segtext.Sect
	ld.Cseek(int64(sect.Vaddr - ld.Segtext.Vaddr + ld.Segtext.Fileoff))
	ld.CodeblkPad(int64(sect.Vaddr), int64(sect.Length), []byte{0x0, 0x0d, 0xea, 0xd1})
	for sect = sect.Next; sect != nil; sect = sect.Next {
		ld.Cseek(int64(sect.Vaddr - ld.Segtext.Vaddr + ld.Segtext.Fileoff))
		ld.Datblk(int64(sect.Vaddr), int64(sect.Length))
	}

	if ld.Segrodata.Filelen > 0 {
		if ld.Debug['v'] != 0 {
			fmt.Fprintf(ld.Bso, "%5.2f rodatblk\n", obj.Cputime())
		}
		ld.Bso.Flush()

		ld.Cseek(int64(ld.Segrodata.Fileoff))
		ld.Datblk(int64(ld.Segrodata.Vaddr), int64(ld.Segrodata.Filelen))
	}

	if ld.Debug['v'] != 0 {
		fmt.Fprintf(ld.Bso, "%5.2f datblk\n", obj.Cputime())
	}
	ld.Bso.Flush()

	ld.Cseek(int64(ld.Segdata.Fileoff))
	ld.Datblk(int64(ld.Segdata.Vaddr), int64(ld.Segdata.Filelen))

	/* output symbol table */
	ld.Symsize = 0

	ld.Lcsize = 0
	symo := uint32(0)
	if ld.Debug['s'] == 0 {
		// TODO: rationalize
		if ld.Debug['v'] != 0 {
			fmt.Fprintf(ld.Bso, "%5.2f sym\n", obj.Cputime())
		}
		ld.Bso.Flush()
		switch ld.HEADTYPE {
		default:
			if ld.Iself {
				symo = uint32(ld.Segdata.Fileoff + ld.Segdata.Filelen)
				symo = uint32(ld.Rnd(int64(symo), int64(ld.INITRND)))
			}
		}

		ld.Cseek(int64(symo))
		switch ld.HEADTYPE {
		default:
			if ld.Iself {
				if ld.Debug['v'] != 0 {
					fmt.Fprintf(ld.Bso, "%5.2f elfsym\n", obj.Cputime())
				}
				ld.Asmelfsym()
				ld.Cflush()
				ld.Cwrite(ld.Elfstrdat)

				if ld.Debug['v'] != 0 {
					fmt.Fprintf(ld.Bso, "%5.2f dwarf\n", obj.Cputime())
				}

				if ld.Linkmode == ld.LinkExternal {
					ld.Elfemitreloc()
				}
			}
		}
	}

	ld.Ctxt.Cursym = nil
	if ld.Debug['v'] != 0 {
		fmt.Fprintf(ld.Bso, "%5.2f header\n", obj.Cputime())
	}
	ld.Bso.Flush()
	ld.Cseek(0)
	switch ld.HEADTYPE {
	default:

	case obj.Hlinux,
		obj.Hfreebsd,
		obj.Hnetbsd,
		obj.Hopenbsd,
		obj.Hsolaris,
		obj.Hnacl:
		ld.Asmbelf(int64(symo))
	}

	ld.Cflush()
	if ld.Debug['c'] != 0 {
		fmt.Printf("textsize=%d\n", ld.Segtext.Filelen)
		fmt.Printf("datsize=%d\n", ld.Segdata.Filelen)
		fmt.Printf("bsssize=%d\n", ld.Segdata.Length-ld.Segdata.Filelen)
		fmt.Printf("symsize=%d\n", ld.Symsize)
		fmt.Printf("lcsize=%d\n", ld.Lcsize)
		fmt.Printf("total=%d\n", ld.Segtext.Filelen+ld.Segdata.Length+uint64(ld.Symsize)+uint64(ld.Lcsize))
	}
}

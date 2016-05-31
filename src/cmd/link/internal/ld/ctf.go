// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"cmd/internal/obj"
	"encoding/binary"
)

var (
	ctfo      int64
	ctfsize   int64
	ctfsym    *LSym
	ctfsympos int64
	ctffile   CtfFile
)

// Ctfemitdebugsections is the main entry point for generating ctf.
func Ctfemitdebugsections() {
	if Debug['t'] != 0 || goos != "solaris" { // disable ctf
		return
	}

	ctffile.Magic = 0xcff1
	ctffile.Version = 2
	ctffile.addString("")

	var label CtfLblent
	label.Label = ctffile.addString(obj.Getgoversion())
	binary.Write(&ctffile.Labels, Ctxt.Arch.ByteOrder, label)

	ctffile.addType("unsafe.Pointer", CTF_TYPE_INFO(CTF_K_POINTER, true, 0), 0)
	ctffile.addType("uintptr", CTF_TYPE_INFO(CTF_K_INTEGER, true, 0), 4)
	ctffile.putUint32(CTF_INT_DATA(0, 0, 64))

	off := ctffile.Labels.Len()
	ctffile.Objtoff = uint32(off)
	off += ctffile.Objects.Len()
	ctffile.Funcoff = uint32(off)
	off += ctffile.Functions.Len()
	ctffile.Typeoff = uint32(off)
	off += ctffile.Types.Len()
	ctffile.Stroff = uint32(off)
	ctffile.Strlen = uint32(len(ctffile.Strings))

	ctfo = Cpos()
	binary.Write(&coutbuf, Ctxt.Arch.ByteOrder, ctffile.CtfHeader)
	Cwrite(ctffile.Labels.Bytes())
	Cwrite(ctffile.Objects.Bytes())
	Cwrite(ctffile.Functions.Bytes())
	Cwrite(ctffile.Types.Bytes())
	Cwrite(ctffile.Strings)
	ctfsize = Cpos() - ctfo
}

const (
	ElfStrCtf = iota
	NElfStrCtf
)

var elfstrctf [NElfStrCtf]int64

func ctfaddshstrings(shstrtab *LSym) {
	if Debug['t'] != 0 || goos != "solaris" { // disable ctf
		return
	}
	elfstrctf[ElfStrCtf] = Addstring(shstrtab, ".SUNW_ctf")
	if Linkmode == LinkExternal {
		// TODO(aram): why only LinkExternal?
		ctfsym = Linklookup(Ctxt, ".SUNW_ctf", 0)
	}
}

func ctfaddelfsectionsyms() {
	if Debug['t'] != 0 || goos != "solaris" { // disable ctf
		return
	}
	if ctfsym != nil {
		ctfsympos = Cpos()
		putelfsectionsym(ctfsym, 0)
	}
}

func ctfaddelfheaders() {
	if Debug['t'] != 0 || goos != "solaris" { // disable ctf
		return
	}
	sh := newElfShdr(elfstrctf[ElfStrCtf])
	sh.type_ = SHT_PROGBITS
	sh.off = uint64(ctfo)
	sh.size = uint64(ctfsize)
	sh.addralign = 4
	if ctfsympos > 0 {
		putelfsymshndx(ctfsympos, sh.shnum)
	}
}

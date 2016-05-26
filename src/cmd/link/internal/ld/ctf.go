// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

var ctfo int64

var ctfsize int64

var ctfsym *LSym

var ctfsympos int64

// Ctfemitdebugsections is the main entry point for generating ctf.
func Ctfemitdebugsections() {
	if Debug['t'] != 0 || goos != "solaris" { // disable ctf
		return
	}
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
	elfstrdbg[ElfStrCtf] = Addstring(shstrtab, ".SUNW_ctf")
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

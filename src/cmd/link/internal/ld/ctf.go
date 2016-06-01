// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"cmd/internal/obj"
	"encoding/binary"
	"fmt"
	"strings"
)

var (
	ctfo      int64
	ctfsize   int64
	ctfsym    *LSym
	ctfsympos int64
	ctffile   CtfFile
)

var typPending = map[string]bool{}

// Define ctftype, for composite ones recurse into constituents.
func ctftype(gotype *LSym) (typno uint16, pending bool) {
	if gotype == nil {
		return 0, false
	}

	if !strings.HasPrefix(gotype.Name, "type.") {
		Diag("ctf: type name doesn't start with \"type.\": %s", gotype.Name)
		return 0, false
	}

	name := gotype.Name[5:] // could also decode from Type.string
	typno, ok := ctffile.Types.byName[name]
	if ok {
		return typno, false
	}
	if typPending[name] {
		return 0, true
	}

	if false && Debug['v'] > 2 {
		fmt.Printf("new type: %v\n", gotype)
	}

	kind := decodetype_kind(gotype)
	bytesize := decodetype_size(gotype)

	switch kind {
	case obj.KindBool:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_INTEGER, true, 0), 4)
		ctffile.putUint32(CTF_INT_DATA(CTF_INT_BOOL, 0, 1))

	case obj.KindInt,
		obj.KindInt8,
		obj.KindInt16,
		obj.KindInt32,
		obj.KindInt64:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_INTEGER, true, 0), 4)
		ctffile.putUint32(CTF_INT_DATA(CTF_INT_SIGNED, 0, uint32(bytesize*8)))

	case obj.KindUint,
		obj.KindUint8,
		obj.KindUint16,
		obj.KindUint32,
		obj.KindUint64,
		obj.KindUintptr:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_INTEGER, true, 0), 4)
		ctffile.putUint32(CTF_INT_DATA(0, 0, uint32(bytesize*8)))

	case obj.KindFloat32:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_FLOAT, true, 0), 4)
		ctffile.putUint32(CTF_FP_DATA(CTF_FP_SINGLE, 0, 32))

	case obj.KindFloat64:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_FLOAT, true, 0), 4)
		ctffile.putUint32(CTF_FP_DATA(CTF_FP_DOUBLE, 0, 64))

	case obj.KindComplex64:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_FLOAT, true, 0), 4)
		ctffile.putUint32(CTF_FP_DATA(CTF_FP_CPLX, 0, 64))

	case obj.KindComplex128:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_FLOAT, true, 0), 4)
		ctffile.putUint32(CTF_FP_DATA(CTF_FP_DCPLX, 0, 128))

	case obj.KindArray:
		// TODO(aram):

	case obj.KindChan:
		// TODO(aram):

	case obj.KindFunc:
		// TODO(aram):

	case obj.KindInterface:
		// TODO(aram):

	case obj.KindMap:
		// TODO(aram):

	case obj.KindPtr:
		vtypno, _ := ctftype(decodetype_ptrelem(gotype))
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_POINTER, true, 0), vtypno)

	case obj.KindSlice:
		// TODO(aram):

	case obj.KindString:
		// TODO(aram):

	case obj.KindStruct:
		if typPending[name] {
			return 0, true
		}
		typPending[name] = true
		nfields := decodetype_structfieldcount(gotype)
		var f string
		var s *LSym
		var mb CtfMember
		for i := 0; i < nfields; i++ {
			f = decodetype_structfieldname(gotype, i)
			s = decodetype_structfieldtype(gotype, i)
			s = decodetype_structfieldtype(gotype, i)
			if f == "" {
				f = s.Name[5:] // skip "type."
			}
			ctftype(s)
		}
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_STRUCT, true, uint16(nfields)), uint16(bytesize))
		for i := 0; i < nfields; i++ {
			f = decodetype_structfieldname(gotype, i)
			s = decodetype_structfieldtype(gotype, i)
			if f == "" {
				f = s.Name[5:] // skip "type."
			}
			mtypno, _ := ctftype(s)
			mb = CtfMember{
				Name:   ctffile.addString(f),
				Type:   mtypno,
				Offset: uint16(decodetype_structfieldoffs(gotype, i)),
			}
			binary.Write(&ctffile.Types, Ctxt.Arch.ByteOrder, mb)
			delete(typPending, name)
		}
		ctffile.addType(name, CTF_TYPE_INFO(CTF_K_TYPEDEF, true, 0), typno)

	case obj.KindUnsafePointer:
		typno = ctffile.addType(name, CTF_TYPE_INFO(CTF_K_POINTER, true, 0), 0)

	default:
		Diag("ctf: definition of unknown kind %d: %s", kind, gotype.Name)
	}

	return typno, false
}

// For use with pass.c::genasmsym
func defctfsymb(sym *LSym, s string, t int, v int64, size int64, ver int, gotype *LSym) {
	if strings.HasPrefix(s, "go.string.") {
		return
	}
	if strings.HasPrefix(s, "runtime.gcbits.") {
		return
	}

	if strings.HasPrefix(s, "type.") && s != "type.*" && !strings.HasPrefix(s, "type..") {
		ctftype(sym)
		return
	}
	switch t {
	default:
		return

	case 'a', 'p':
		ctftype(gotype)
	}
}

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

	genasmsym(defctfsymb)
	if Debug['v'] != 0 {
		fmt.Fprintf(&Bso, "%5.2f ctf\n", obj.Cputime())
	}

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
	sh.link = uint32(elfshname(".symtab").shnum)
	if ctfsympos > 0 {
		putelfsymshndx(ctfsympos, sh.shnum)
	}
}

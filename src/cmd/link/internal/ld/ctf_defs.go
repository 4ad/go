// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"bytes"
)

// All data is in target-endian order.

type CtfPreamble struct {
	Magic   uint16 /* magic number (CTF_MAGIC) */
	Version uint8  /* data format version number (CTF_VERSION) */
	Flags   uint8  /* flags (see below) */
}

const CTF_F_COMPRESS = 1

type CtfHeader struct {
	CtfPreamble
	Parlabel uint32 /* ref to name of parent lbl uniq'd against */
	Parname  uint32 /* ref to basename of parent */
	Lbloff   uint32 /* offset of label section */
	Objtoff  uint32 /* offset of object section */
	Funcoff  uint32 /* offset of function section */
	Typeoff  uint32 /* offset of type section */
	Stroff   uint32 /* offset of string section */
	Strlen   uint32 /* length of string section in bytes */
}

const (
	CTF_K_UNKNOWN  = 0
	CTF_K_INTEGER  = 1
	CTF_K_FLOAT    = 2
	CTF_K_POINTER  = 3
	CTF_K_ARRAY    = 4
	CTF_K_FUNCTION = 5
	CTF_K_STRUCT   = 6
	CTF_K_UNION    = 7
	CTF_K_ENUM     = 8
	CTF_K_FORWARD  = 9
	CTF_K_TYPEDEF  = 10
	CTF_K_VOLATILE = 11
	CTF_K_CONST    = 12
	CTF_K_RESTRICT = 13
)

const CTF_MAX_VLEN = 0x3ff

func CTF_INFO_KIND(info uint16) uint16 {
	return (info & 0xf800) >> 11
}

func CTF_INFO_ISROOT(info uint16) bool {
	return (info&0x0400)>>10 == 1
}

func CTF_INFO_VLEN(info uint16) uint16 {
	return info & CTF_MAX_VLEN
}

func CTF_TYPE_INFO(kind uint16, isroot bool, vlen uint16) uint16 {
	var info uint16 = (kind << 11) | (vlen & CTF_MAX_VLEN)
	if isroot {
		info |= 1 << 10
	}
	return info
}

type CtfLblent struct {
	Label   uint32 /* ref to name of label */
	Typeidx uint32 /* last type associated with this label */
}

const (
	CTF_MAX_SIZE   = 0xfffe /* max size of a type in bytes */
	CTF_LSIZE_SENT = 0xffff /* sentinel for ctt_size */
	CTF_MAX_LSIZE  = 1<<64 - 1
)

type CtfStype struct {
	Name       uint32 /* reference to name in string table */
	Info       uint16 /* encoded kind, variant length */
	SizeOrType uint16 /* size of entire type in bytes, or reference to another type */
}

type CtfType struct {
	Name       uint32 /* reference to name in string table */
	Info       uint16 /* encoded kind, variant length */
	SizeOrType uint16 /* always CTF_LSIZE_SENT */
	SizeHi     uint32 /* high 32 bits of type size in bytes */
	SizeLo     uint32 /* low 32 bits of type size in bytes */
}

func CTF_INT_ENCODING(data uint32) uint32 {
	return (data & 0xff000000) >> 24
}

func CTF_INT_OFFSET(data uint32) uint32 {
	return (data & 0x00ff0000) >> 16
}

func CTF_INT_BITS(data uint32) uint32 {
	return data & 0x0000ffff
}

func CTF_INT_DATA(encoding, offset, bits uint32) uint32 {
	return (encoding << 24) | (offset << 16) | bits
}

const (
	CTF_INT_SIGNED  = 0x01
	CTF_INT_CHAR    = 0x02
	CTF_INT_BOOL    = 0x04
	CTF_INT_VARARGS = 0x08
)

func CTF_FP_ENCODING(data uint32) uint32 {
	return (data & 0xff000000) >> 24
}

func CTF_FP_OFFSET(data uint32) uint32 {
	return (data & 0x00ff0000) >> 16
}

func CTF_FP_BITS(data uint32) uint32 {
	return data & 0x0000ffff
}

func CTF_FP_DATA(encoding, offset, bits uint32) uint32 {
	return (encoding << 24) | (offset << 16) | bits
}

const (
	CTF_FP_SINGLE   = 1  /* IEEE 32-bit float encoding */
	CTF_FP_DOUBLE   = 2  /* IEEE 64-bit float encoding */
	CTF_FP_CPLX     = 3  /* Complex encoding */
	CTF_FP_DCPLX    = 4  /* Double complex encoding */
	CTF_FP_LDCPLX   = 5  /* Long double complex encoding */
	CTF_FP_LDOUBLE  = 6  /* Long double encoding */
	CTF_FP_INTRVL   = 7  /* Interval (2x32-bit) encoding */
	CTF_FP_DINTRVL  = 8  /* Double interval (2x64-bit) encoding */
	CTF_FP_LDINTRVL = 9  /* Long double interval (2x128-bit) encoding */
	CTF_FP_IMAGRY   = 10 /* Imaginary (32-bit) encoding */
	CTF_FP_DIMAGRY  = 11 /* Long imaginary (64-bit) encoding */
	CTF_FP_LDIMAGRY = 12 /* Long double imaginary (128-bit) encoding */
)

type CtfArray struct {
	Contents uint16 /* reference to type of array contents */
	Index    uint16 /* reference to type of array index */
	Nelems   uint32 /* number of elements */
}

type CtfMember struct {
	Name   uint32 /* reference to name in string table */
	Type   uint16 /* reference to type of member */
	Offset uint16 /* offset of this member in bits */
}

type CtfLmember struct {
	Name     uint32 /* reference to name in string table */
	Type     uint16 /* reference to type of member */
	_        uint16 /* padding */
	OffsetHi uint32 /* high 32 bits of member offset in bits */
	OffsetLo uint32 /* low 32 bits of member offset in bits */
}

type CtfEnum struct {
	Name  uint32 /* reference to name in string table */
	Value int32  /* value associated with this name */
}

type CtfFile struct {
	CtfHeader
	Labels    bytes.Buffer
	Objects   bytes.Buffer
	Functions bytes.Buffer
	Types     bytes.Buffer
	Strings   []byte
}

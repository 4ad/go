// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/sys"
	"unsafe"
)

const (
	c0 = uintptr((8-sys.PtrSize)/4*2860486313 + (sys.PtrSize-4)/4*33054211828000289)
	c1 = uintptr((8-sys.PtrSize)/4*3267000013 + (sys.PtrSize-4)/4*23344194077549503)
)

// type algorithms - known to compiler
const (
	alg_MEM = iota
	alg_MEM0
	alg_MEM8
	alg_MEM16
	alg_MEM32
	alg_MEM64
	alg_MEM128
	alg_NOEQ
	alg_NOEQ0
	alg_NOEQ8
	alg_NOEQ16
	alg_NOEQ32
	alg_NOEQ64
	alg_NOEQ128
	alg_STRING
	alg_INTER
	alg_NILINTER
	alg_SLICE
	alg_FLOAT32
	alg_FLOAT64
	alg_CPLX64
	alg_CPLX128
	alg_max
)

type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

func memhash0(p unsafe.Pointer, h uintptr) uintptr {
	return h
}
func memhash8(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 1)
}
func memhash16(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 2)
}
func memhash32(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 4)
}
func memhash64(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 8)
}
func memhash128(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 16)
}

// memhash_varlen is defined in assembly because it needs access
// to the closure.  It appears here to provide an argument
// signature for the assembly routine.
func memhash_varlen(p unsafe.Pointer, h uintptr) uintptr

var algarray = [alg_max]typeAlg{
	alg_MEM:      {nil, nil}, // not used
	alg_MEM0:     {nil, nil}, //{memhash0, memequal0},
	alg_MEM8:     {nil, nil}, //{memhash8, memequal8},
	alg_MEM16:    {nil, nil}, //{memhash16, memequal16},
	alg_MEM32:    {nil, nil}, //{memhash32, memequal32},
	alg_MEM64:    {nil, nil}, //{memhash64, memequal64},
	alg_MEM128:   {nil, nil}, //{memhash128, memequal128},
	alg_NOEQ:     {nil, nil},
	alg_NOEQ0:    {nil, nil},
	alg_NOEQ8:    {nil, nil},
	alg_NOEQ16:   {nil, nil},
	alg_NOEQ32:   {nil, nil},
	alg_NOEQ64:   {nil, nil},
	alg_NOEQ128:  {nil, nil},
	alg_STRING:   {nil, nil}, //{strhash, strequal},
	alg_INTER:    {nil, nil}, //{interhash, interequal},
	alg_NILINTER: {nil, nil}, //{nilinterhash, nilinterequal},
	alg_SLICE:    {nil, nil},
	alg_FLOAT32:  {nil, nil}, //{f32hash, f32equal},
	alg_FLOAT64:  {nil, nil}, //{f64hash, f64equal},
	alg_CPLX64:   {nil, nil}, //{c64hash, c64equal},
	alg_CPLX128:  {nil, nil}, //{c128hash, c128equal},
}

func f32hash(p unsafe.Pointer, h uintptr) uintptr {
	f := *(*float32)(p)
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		return c1 * (c0 ^ h ^ uintptr(fastrand1())) // any kind of NaN
	default:
		return memhash(p, h, 4)
	}
}

func f64hash(p unsafe.Pointer, h uintptr) uintptr {
	f := *(*float64)(p)
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		return c1 * (c0 ^ h ^ uintptr(fastrand1())) // any kind of NaN
	default:
		return memhash(p, h, 8)
	}
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func strhash(a unsafe.Pointer, h uintptr) uintptr {
	x := (*stringStruct)(a)
	return memhash(x.str, h, uintptr(x.len))
}

func interhash(p unsafe.Pointer, h uintptr) uintptr {
	a := (*iface)(p)
	tab := a.tab
	if tab == nil {
		return h
	}
	t := tab._type
	fn := t.alg.hash
	if fn == nil {
		panic(errorString("hash of unhashable type " + *t._string))
	}
	if isDirectIface(t) {
		return c1 * fn(unsafe.Pointer(&a.data), h^c0)
	} else {
		return c1 * fn(a.data, h^c0)
	}
}

func nilinterhash(p unsafe.Pointer, h uintptr) uintptr {
	a := (*eface)(p)
	t := a._type
	if t == nil {
		return h
	}
	fn := t.alg.hash
	if fn == nil {
		panic(errorString("hash of unhashable type " + *t._string))
	}
	if isDirectIface(t) {
		return c1 * fn(unsafe.Pointer(&a.data), h^c0)
	} else {
		return c1 * fn(a.data, h^c0)
	}
}

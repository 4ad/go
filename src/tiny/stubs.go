// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm64

package runtime

import "unsafe"

// Declarations for runtime services implemented in C or assembly.

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
const regSize = 4 << (^uintreg(0) >> 63) // unsafe.Sizeof(uintreg(0)) but an ideal const

// Should be a built-in for unsafe.Pointer?
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to a single xor instruction.
// USE CAREFULLY!
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func getg() *g

//go:noescape
func getcallersp(unsafe.Pointer) uintptr

//go:noescape
func getcallerpc(unsafe.Pointer) uintptr

//go:noescape
func jmpdefer(fv *funcval, argp uintptr)

func return0()

// defined in sys_$GOOS_$GOARCH.s
func read(fd int32, p unsafe.Pointer, n int32) int32
func close(fd int32) int32

func exit(code int32)
func nanotime() int64
func usleep(usec uint32)

func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) unsafe.Pointer
func munmap(addr unsafe.Pointer, n uintptr)

//go:noescape
func write(fd uintptr, p unsafe.Pointer, n int32) int32

//go:noescape
func open(name *byte, mode, perm int32) int32

func madvise(addr unsafe.Pointer, n uintptr, flags int32)

func brk(addr uintptr) uintptr

func Exit(x int32) {
	exit(x)
}

// defined in mem{clr,move}_$GOARCH.s
//go:noescape
func memmove(to unsafe.Pointer, frm unsafe.Pointer, length uintptr)

//go:noescape
func memclr(ptr unsafe.Pointer, length uintptr)

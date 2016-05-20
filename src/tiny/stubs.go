// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/sys"
	"unsafe"
)

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

func close(fd int32) int32

func Exit(x int32) {
	exit(x)
}

// defined in mem{clr,move}_$GOARCH.s
//go:noescape
func memmove(to unsafe.Pointer, frm unsafe.Pointer, length uintptr)

//go:noescape
func memclr(ptr unsafe.Pointer, length uintptr)

//go:noescape
func asmcgocall(fn, arg unsafe.Pointer) int32

//go:nosplit
func exitsyscall(dummy int32) {}

//go:nosplit
func entersyscallblock(dummy int32) {}

func fastrand1() uint32

func GOMAXPROCS(n int) int {
	return 1
}

const GOARCH string = sys.TheGoarch
const GOOS string = sys.TheGoos

func GOROOT() string {
	return sys.DefaultGoroot
}

func GC() {}

func SetFinalizer(obj interface{}, finalizer interface{}) {
}

//go:linkname reflect_chanlen reflect.chanlen
func reflect_chanlen(c *hchan) int {
	if c == nil {
		return 0
	}
	return int(c.qcount)
}

//go:linkname reflect_chancap reflect.chancap
func reflect_chancap(c *hchan) int {
	if c == nil {
		return 0
	}
	return int(c.dataqsiz)
}

func callwritebarrier(typ *_type, frame unsafe.Pointer, framesize, retoffset uintptr) {}

func badreflectcall() {
	panic("runtime: arg size to reflect.call more than 1GB")
}
